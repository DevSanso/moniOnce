package applnew

import (
	"agent_common/pkg/applnew/loader"
	"agent_common/pkg/applnew/logger"
	"agent_common/pkg/applnew/thread"
	apptype "agent_common/pkg/applnew/types"
	"agent_common/pkg/util/collection"
	"agent_common/pkg/util/types"
	"agent_common/pkg/util/writer"
	"context"
	"database/sql"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Application interface {
	Run(context.Context) error
}

type ApplicationInit[PUSH any, CONN io.Closer, FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]] interface {
	Init(data apptype.InitData[PUSH, CONN, FLAG, FLAGPTR]) error 
	Application
}

type implApplication[PUSH any, CONN io.Closer, FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]] struct {
	setting apptype.SettingData
	configDB loader.Configure[apptype.ApplConfData, apptype.AppSyncData, FLAG, *apptype.ApplConfData, * apptype.AppSyncData, FLAGPTR]
	config *apptype.ApplConfData

	collectConnPool apptype.CollectConnPool[CONN]
	dataPusher apptype.DataPusher[PUSH]

	closer struct {
		configDBCloser io.Closer

		cronLogCloser  io.Closer
		collectLogCloser  io.Closer
		pushLogCloser  io.Closer
		intervalLogCloser  io.Closer
		initLogCloser  io.Closer

		cronQCloser io.Closer
		pushQCloser io.Closer
		collectQCloser io.Closer

		collectConnPoolCloser io.Closer
		dataPusherCloser      io.Closer
	}

	queue struct {
		cronQ types.Queue[string]
		pushQ types.Queue[*PUSH]
		collectQ types.Queue[string]
	}

	thread struct {
		collectT []thread.CollectThread[PUSH, CONN]
		cronT    []thread.CronThread[PUSH,FLAG,FLAGPTR]
		pushT    []thread.PushThread[PUSH]
		intevalT thread.IntervalThread[FLAG,FLAGPTR]
	}

	loggers struct {
		initLogger logger.LevelLogger
		collectLogger logger.LevelLogger
		pushLogger logger.LevelLogger
		intervalLogger logger.LevelLogger
		cronLogger logger.LevelLogger
	}
}

func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) initSetting(settingPath string) error {
	initData, readErr := os.ReadFile(settingPath)
	if readErr != nil {
		return readErr
	}
	if err := toml.Unmarshal(initData, &i.setting); err != nil {
		return err
	}

	return nil
}


func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) initQueue() {
	i.queue.collectQ = collection.NewStdQueue[string](i.config.Queue.CollectSize)
	i.queue.cronQ = collection.NewStdQueue[string](i.config.Queue.CollectSize)
	i.queue.pushQ = collection.NewStdQueue[*PUSH](i.config.Queue.CollectSize)

	i.closer.collectQCloser = i.queue.collectQ
	i.closer.cronQCloser = i.queue.cronQ
	i.closer.pushQCloser = i.queue.pushQ
}

func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) initThread(data apptype.InitData[PUSH, CONN, FLAG, FLAGPTR]) {
	i.thread.intevalT = thread.NewIntervalThread(i.loggers.intervalLogger, i.configDB, i.queue.collectQ, i.queue.cronQ)
	i.thread.cronT = make([]thread.CronThread[PUSH, FLAG, FLAGPTR], 1)
	for n := 0; n < i.config.Thread.CronCount; n++ {
		t := thread.NewCronThread(i.queue.cronQ, i.queue.pushQ, data.CronM, i.loggers.cronLogger, i.configDB)
		i.thread.cronT = append(i.thread.cronT, t)
	}
	i.thread.collectT = make([]thread.CollectThread[PUSH, CONN], 1)
	for n := 0; n < i.config.Thread.CollectCount; n++ {
		t := thread.NewCollectThread(i.queue.cronQ, i.collectConnPool, i.queue.pushQ, data.CollectM, i.loggers.collectLogger)
		i.thread.collectT = append(i.thread.collectT, t)
	}
	i.thread.pushT = make([]thread.PushThread[PUSH], 1)
	for n := 0; n < i.config.Thread.PushCount; n++ {
		t := thread.NewPushThread(i.queue.pushQ, i.loggers.pushLogger, i.dataPusher)
		i.thread.pushT = append(i.thread.pushT, t)
	}
}


func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) Init(data apptype.InitData[PUSH, CONN, FLAG, FLAGPTR]) error {
	if err := i.initSetting(data.SettingPath); err != nil {
		return err
	}
	if err := i.initLogger(); err != nil {
		return err
	}
	if err := i.connectConfig(); err != nil {
		i.loggers.initLogger.Error("connect config db failed :", err.Error())
		return err
	}
	if err := i.initConfig(); err != nil {
		i.loggers.initLogger.Error("init config failed :", err.Error())
		return err
	}
	if err := i.initCollectConnPool(data); err != nil {
		i.loggers.initLogger.Error("init collect conn pool : ", err.Error())
		return err
	}
	if err := i.initDataPusherPool(data); err != nil {
		i.loggers.initLogger.Error("init data pusher conn pool : ", err.Error())
		return err
	}

	i.initQueue()
	i.initThread(data)

	return nil
}

func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) AsyncThreads(ctx context.Context) error {
	go i.thread.intevalT.Run(ctx)
	for _, t := range i.thread.pushT {
		go t.Run(ctx)
	}
	for _, t := range i.thread.collectT {
		go t.Run(ctx)
	}
	for _, t := range i.thread.cronT {
		go t.Run(ctx)
	}

	return nil
}

func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) Close() error {
	i.closer.collectLogCloser.Close()
	i.closer.configDBCloser.Close()
	i.closer.cronLogCloser.Close()
	i.closer.initLogCloser.Close()
	i.closer.intervalLogCloser.Close()
	i.closer.pushLogCloser.Close()

	i.closer.pushQCloser.Close()
	i.closer.collectQCloser.Close()
	i.closer.cronQCloser.Close()

	i.closer.collectConnPoolCloser.Close()
	i.closer.configDBCloser.Close()
	
	return nil
}


func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) Run(ctx context.Context) error {
	i.AsyncThreads(ctx)
	
	isStop := false
	for !isStop {
		select {
		case <-ctx.Done():
			isStop = true
		default:
		}

		time.Sleep(time.Second * 5)
	}
	return nil
}

// connectConfig implements Application.
func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) connectConfig() error {
	db, dbErr := sql.Open(i.setting.ConfigDb.Drvier, i.setting.ConfigDb.Dsn)
	if dbErr != nil {
		return dbErr
	}
	i.closer.configDBCloser = db
	i.configDB = loader.NewSQLConfigure[apptype.ApplConfData, apptype.AppSyncData, FLAG, *apptype.ApplConfData, * apptype.AppSyncData, FLAGPTR](
		db, i.setting.ObjectId, i.loggers.initLogger)

	return dbErr
}

func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) initCollectConnPool(data apptype.InitData[PUSH, CONN, FLAG, FLAGPTR]) error {
	p, err := data.GetConnPoolFn(
		i.config.CollecbDbConfig.IP, 
		i.config.CollecbDbConfig.Port, 
		i.config.CollecbDbConfig.User, 
		i.config.CollecbDbConfig.Password, 
		i.config.CollecbDbConfig.Dbname)
	if err != nil {
		return err
	}
	i.collectConnPool = p
	i.closer.collectConnPoolCloser = i.collectConnPool
	return nil
}

func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) initDataPusherPool(data apptype.InitData[PUSH, CONN, FLAG, FLAGPTR]) error {
	p, err := data.GetPusherFn(
		i.config.DataPusherConfig.IP, 
		i.config.DataPusherConfig.Port, 
		i.config.DataPusherConfig.User, 
		i.config.DataPusherConfig.Password, 
		i.config.DataPusherConfig.Dbname)
	if err != nil {
		return err
	}
	i.dataPusher = p
	i.closer.dataPusherCloser = i.dataPusher
	return nil
}
// initConfig implements Application.
func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) initConfig() error {
	var err error
	i.config,err = i.configDB.LoadConfig()
	return err
}

// initLogger implements Application.
func (i *implApplication[PUSH, CONN, FLAG, FLAGPTR]) initLogger() error {
	_, err := os.Stat(i.setting.LogConfig.Dir)
	if os.IsNotExist(err) {
		return err
	}

	const cronName = "cron.log"
	const collectName = "collect.log"
	const pushName = "push.log"
	const intervalName = "interval.log"
	const initName     = "init.log"

	var cronWriter io.WriteCloser = nil
	var collectWriter io.WriteCloser = nil
	var pushWriter io.WriteCloser = nil
	var intervalWriter io.WriteCloser = nil
	var initWriter io.WriteCloser = nil

	var writerErr error = nil
	initWriter, writerErr = writer.NewSizeLimitedWriter(i.setting.LogConfig.Dir, initName, int(i.setting.LogConfig.Size))
	if writerErr != nil {
		return writerErr
	} else {
		var loggerErr error = nil
		i.closer.initLogCloser = initWriter
		i.loggers.initLogger, loggerErr = logger.NewSlogLogger(initWriter, logger.LogLevel(i.setting.LogConfig.Level))

		if loggerErr != nil {
			return loggerErr
		}
	}

	cronWriter, writerErr = writer.NewSizeLimitedWriter(i.setting.LogConfig.Dir, cronName, int(i.setting.LogConfig.Size))
	if writerErr != nil {
		i.loggers.initLogger.Error("init failed cronwriter")
		return writerErr
	} else {
		var loggerErr error = nil
		i.closer.cronLogCloser = cronWriter
		i.loggers.cronLogger, loggerErr = logger.NewSlogLogger(cronWriter, logger.LogLevel(i.setting.LogConfig.Level))
		if loggerErr != nil {
			i.loggers.initLogger.Error("init failed cronlogger")
			return loggerErr
		}
	}
	collectWriter, writerErr = writer.NewSizeLimitedWriter(i.setting.LogConfig.Dir, collectName, int(i.setting.LogConfig.Size))
	if writerErr != nil {
		i.loggers.initLogger.Error("init failed collectWriter")
		return writerErr
	} else {
		var loggerErr error = nil
		i.closer.collectLogCloser = collectWriter
		i.loggers.collectLogger, loggerErr = logger.NewSlogLogger(collectWriter, logger.LogLevel(i.setting.LogConfig.Level))
		if loggerErr != nil {
			i.loggers.initLogger.Error("init failed collectlogger")
			return loggerErr
		}
	}
	pushWriter, writerErr = writer.NewSizeLimitedWriter(i.setting.LogConfig.Dir, pushName, int(i.setting.LogConfig.Size))
	if writerErr != nil {
		i.loggers.initLogger.Error("init failed pushwriter")
		return writerErr
	} else {
		var loggerErr error = nil
		i.closer.pushLogCloser = pushWriter
		i.loggers.pushLogger, loggerErr = logger.NewSlogLogger(pushWriter, logger.LogLevel(i.setting.LogConfig.Level))
		if loggerErr != nil {
			i.loggers.initLogger.Error("init failed pushlogger")
			return loggerErr
		}
	}
	intervalWriter, writerErr = writer.NewSizeLimitedWriter(i.setting.LogConfig.Dir, intervalName, int(i.setting.LogConfig.Size))
	if writerErr != nil {
		i.loggers.initLogger.Error("init failed intervalwriter")
		return writerErr
	} else {
		var loggerErr error = nil
		i.closer.intervalLogCloser = intervalWriter
		i.loggers.intervalLogger, loggerErr = logger.NewSlogLogger(intervalWriter, logger.LogLevel(i.setting.LogConfig.Level))
		if loggerErr != nil {
			i.loggers.initLogger.Error("init failed intervallogger")
			return loggerErr
		}
	}

	return nil

}

func NewApplication[PUSH any, CONN io.Closer, FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]]() ApplicationInit[PUSH, CONN, FLAG, FLAGPTR] {
	return &implApplication[PUSH, CONN, FLAG, FLAGPTR]{}
}
