package appl

import (
	"agent_common/pkg/config"
	"agent_common/pkg/logger"
	"agent_common/pkg/util/collection"
	"agent_common/pkg/util/types"
	"agent_common/pkg/util/writer"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/sys/unix"
)

type AgentApplication interface {
	Run(ctx context.Context) error 
}

type implAgentApplication [DB any, DBCONF types.StrCloneToGetter]struct {
	levelLogger       logger.LevelLogger
	dataLoggers       map[int]logger.DataLogger[any]
	applicationConfig *AgentApplicationConfig
	configDb          *sql.DB

	genDbLoggerFn      types.GenDbLoggerFn[DBCONF]
	genTargetDbPoolrFn types.GenTargetDbPoolFn[DB, DBCONF]
	intervals         []types.IntervalRegister[DB]

	applCustomConfig      map[string]string
	
	lazy struct {
		dbLogger          logger.DbLogger
		targetDbPool      types.TargetDbPool[DB]
		dbConfig          types.StrCloneToGetter

		threads struct {
			interval InternvalThread[DB]
			collect  []CollectThread[DB]
			push     []PushThread
		}

		queue struct {
			collectQueue types.Queue[types.CollectFn[DB]]
			pushQueue     types.Queue[types.CollectFnRet]
		}

		closes struct {
			dbLoggerCloser          io.Closer
			collectQCloser          io.Closer
			pushQCloser             io.Closer
			targetDbPoolCloser      io.Closer
		}
	}
}

type AgentApplicationExtendInitConfig[DB any, DBCONF types.StrCloneToGetter] struct {
	GenDbLoggerFn     types.GenDbLoggerFn[DBCONF]
	GenTargetDbPoolrFn types.GenTargetDbPoolFn[DB, DBCONF]
	DataLoggers       map[int]logger.DataLogger[any]
	Intervals         []types.IntervalRegister[DB]
}

func InitAgentApplication[DB any, DBCONF types.StrCloneToGetter](configPath string, 
		extend AgentApplicationExtendInitConfig[DB, DBCONF],
		applCustomConfig map[string]string) (AgentApplication, error) {
	var initConfig AgentApplicationInitConfig
	initData, readErr := os.ReadFile(configPath)
	if readErr != nil {
		return nil, readErr
	}
	if err := toml.Unmarshal(initData, &initConfig); err != nil {
		return nil, err
	}
		
	logWriter, writerErr := writer.NewLogWriter(initConfig.LogConfig.Path)
	if writerErr != nil {
		return nil, writerErr
	}

	levelLogger, loggerErr := logger.NewSlogLogger(logWriter, logger.LogLevel(initConfig.LogConfig.Level))
	if loggerErr != nil {
		logWriter.Close()
		return nil, loggerErr
	}

	db, dbErr := sql.Open(initConfig.ConfigDb.Drvier, initConfig.ConfigDb.Dsn)
	if dbErr != nil {
		levelLogger.Error("open config db error : ", dbErr)
		logWriter.Close()
		return nil, dbErr
	}

	configGen := config.NewSQLConfigure[*AgentApplicationConfig](nil, initConfig.ObjectId, levelLogger)
	configData, configGenErr := configGen.Load()
	if configGenErr != nil {
		logWriter.Close()
		db.Close()
		return nil, configGenErr
	}

	return &implAgentApplication[DB, DBCONF]{
		levelLogger:       levelLogger,
		dataLoggers:       extend.DataLoggers,
		applicationConfig: configData,
		configDb:          db,
		genDbLoggerFn: extend.GenDbLoggerFn,
		intervals: extend.Intervals,
		genTargetDbPoolrFn : extend.GenTargetDbPoolrFn,
		applCustomConfig : applCustomConfig,
	}, nil
}

func (ia *implAgentApplication[DB,DBCONF]) asyncRun(ctx context.Context) {
	isStop := true
	queuePtr := &ia.lazy.queue

	for !isStop {
		nowSec := time.Now().Second() % 60

		if nowSec == 0 {
			ia.levelLogger.Info(fmt.Sprintf("CollectQueue [%d/%d]",
				queuePtr.collectQueue.Count(), queuePtr.collectQueue.Count()))

			ia.levelLogger.Info(fmt.Sprintf("PushQueue [%d/%d]",
				queuePtr.pushQueue.Count(), queuePtr.pushQueue.Count()))
		}

		if nowSec == 30 {
			var mem runtime.MemStats
			var rusage unix.Rusage
			runtime.ReadMemStats(&mem)
			err := unix.Getrusage(unix.RUSAGE_SELF, &rusage)

			ia.levelLogger.Info(fmt.Sprintf("ALLOC [alloc:%d], HEAP [isuse : %d, idle :%d]", mem.Alloc, mem.HeapInuse, mem.HeapIdle))

			if err == nil {
				ia.levelLogger.Debug(fmt.Sprintf(
					"CPU TIME [user:%d, sys:%d]", rusage.Utime.Usec, rusage.Stime.Usec))
				ia.levelLogger.Debug(fmt.Sprintf(
					"Memory [rss : %d]", rusage.Maxrss))
			} else {
				ia.levelLogger.Error("Get failed rusage : ", err.Error())
			}
		}


		select {
		case <-ctx.Done():
			isStop = true
		default:
		}

		time.Sleep(time.Second * 1)
	}

	ia.close()
} 

func (ia *implAgentApplication[DB, DBCONF]) genIntervalMap() (map[int][]types.IntervalRegister[DB],error) {
	ret := make(map[int][]types.IntervalRegister[DB])

	for _, confInterval := range ia.applicationConfig.intervals {
		for _, codeInterval := range ia.intervals {
			if confInterval.key == codeInterval.Name {
				if ret[confInterval.inetervalSec] == nil {
					ret[confInterval.inetervalSec] = make([]types.IntervalRegister[DB], 0)
				}

				ret[confInterval.inetervalSec] = append(ret[confInterval.inetervalSec], codeInterval)
				ia.levelLogger.Info(fmt.Sprintf("collect thread [%s] init interval [sec:%d]", codeInterval.Name, confInterval.inetervalSec))
			}
		}
	}

	return ret, nil
}

func (ia *implAgentApplication[DB, DBCONF]) lazyInit() error {
	lazyPtr := &ia.lazy
	
	dbconf := *new(DBCONF)
	dbconf.CloneFromGetter(ia.applicationConfig)

	dblogger, dbloggerCloser, dbLoggerErr := ia.genDbLoggerFn(dbconf)
	if dbLoggerErr != nil {
		ia.levelLogger.Error("db logger init failed :", dbLoggerErr)
		return dbLoggerErr
	}

	lazyPtr.dbLogger = dblogger
	lazyPtr.closes.dbLoggerCloser = dbloggerCloser
	lazyPtr.dbConfig = dbconf

	targetDbPool, targetDbPoolCloser, targetdbErr := ia.genTargetDbPoolrFn(dbconf)
	if targetdbErr != nil {
		ia.levelLogger.Error("target db pool init failed :", targetdbErr)
		return targetdbErr
	}

	lazyPtr.closes.targetDbPoolCloser = targetDbPoolCloser
	lazyPtr.targetDbPool = targetDbPool

	cQ := collection.NewStdQueue[types.CollectFn[DB]](ia.applicationConfig.Queue.CollectSize)
	pQ := collection.NewStdQueue[types.CollectFnRet](ia.applicationConfig.Queue.PushSize)

	lazyPtr.closes.collectQCloser = cQ
	lazyPtr.closes.pushQCloser = pQ

	var collects []CollectThread[DB] = make([]CollectThread[DB], 0)
	var pushs    []PushThread = make([]PushThread, 0)
	

	for i := 0; i< ia.applicationConfig.Thread.CollectCount ; i ++ {
		collects = append(collects, newCollectThread(pQ, cQ, ia.levelLogger, lazyPtr.targetDbPool, ia.applCustomConfig))
	}

	for i := 0; i< ia.applicationConfig.Thread.PushCount ; i ++ {
		pushs = append(pushs, newPushThread(pQ, ia.levelLogger, lazyPtr.dbLogger))
	}

	lazyPtr.threads.collect = collects
	lazyPtr.threads.push = pushs

	geninterval, genIntervalErr := ia.genIntervalMap()
	if genIntervalErr != nil {
		ia.levelLogger.Error("interval map create failed :", genIntervalErr.Error())
		return genIntervalErr
	}

	var interval = newIntervalThread(geninterval, cQ, ia.levelLogger)
	lazyPtr.threads.interval = interval

	return nil
}

func (ia *implAgentApplication[DB, DBCONF])close() error {
	c := &ia.lazy.closes

	if c.collectQCloser != nil {
		c.collectQCloser.Close()
	}

	if c.pushQCloser != nil {
		c.pushQCloser.Close()
	}

	if c.dbLoggerCloser != nil {
		c.dbLoggerCloser.Close()
	}

	if c.targetDbPoolCloser != nil {
		c.targetDbPoolCloser.Close()
	}

	return nil
}

func (ia *implAgentApplication[DB, DBCONF]) Run(ctx context.Context) error {
	if err := ia.lazyInit(); err != nil {
		ia.levelLogger.Error("lazy init failed :", err)
		ia.close()
		return err
	}

	for _, t := range ia.lazy.threads.push {
		select {
		case <-ctx.Done():
			ia.close()
			return fmt.Errorf("AgnetApplication Stop before init thread")
		default:
		}
		go t.Run(ctx)
	}

	for _, t := range ia.lazy.threads.push {
		select {
		case <-ctx.Done():
			ia.close()
			return fmt.Errorf("AgnetApplication Stop before init thread")
		default:
		}
		go t.Run(ctx)
	}

	go ia.asyncRun(ctx)

	return nil
}