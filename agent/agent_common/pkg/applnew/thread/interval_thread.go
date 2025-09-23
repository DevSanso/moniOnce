package thread

import (
	"agent_common/pkg/applnew/loader"
	"agent_common/pkg/applnew/logger"
	apptype "agent_common/pkg/applnew/types"
	"agent_common/pkg/util/types"
	"context"
	"maps"
	"time"
)

type IntervalThread[FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]] struct {
	intervalLogger logger.LevelLogger
	confLoader loader.Configure[apptype.ApplConfData, apptype.AppSyncData, FLAG, *apptype.ApplConfData, *apptype.AppSyncData, FLAGPTR]

	collectSend    types.Pusher[string]
	cronSend     types.Pusher[string]

	intervals map[string]int
	crons     map[string]int

	isCollectAndCronStop bool
}

func NewIntervalThread[FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]](
	logger logger.LevelLogger, 
	confLoader loader.Configure[apptype.ApplConfData, apptype.AppSyncData, FLAG, *apptype.ApplConfData, *apptype.AppSyncData, FLAGPTR],
	collectSend    types.Pusher[string],
	cronSend     types.Pusher[string]) IntervalThread[FLAG, FLAGPTR] {
	
	return IntervalThread[FLAG, FLAGPTR]{
		intervalLogger: logger,
		confLoader: confLoader,
		collectSend: collectSend,
		cronSend: cronSend,

		intervals: make(map[string]int),
		crons : make(map[string]int),
		isCollectAndCronStop: false,
	}
}

func(it *IntervalThread[FLAG, FLAGPTR])syncSetting() error {
	syncData, syncErr := it.confLoader.LoadSync()
	if syncErr != nil {
		return syncErr
	}

	it.isCollectAndCronStop = syncData.Stop

	maps.Copy(it.intervals, syncData.Intervals)
	maps.Copy(it.crons, syncData.Cron)
	return nil
}

func(it *IntervalThread[FLAG, FLAGPTR])Run(ctx context.Context) error {
	isStop := false
	oldMin := 0 
	oldStopSync := false

	for !isStop {
		nowSec := time.Now().Second() % 60
		nowMin := time.Now().Minute() % 60

		if nowMin != oldMin {
			oldMin = nowMin
			if syncErr := it.syncSetting(); syncErr != nil {
				it.intervalLogger.Error("sync failed :", syncErr.Error())
			}

			if it.isCollectAndCronStop {
				if it.isCollectAndCronStop != oldStopSync{
					it.intervalLogger.Info("reset stop agent")
				}

				oldStopSync = it.isCollectAndCronStop
				time.Sleep(1 * time.Second)
				continue
			} else {
				if it.isCollectAndCronStop != oldStopSync{
					it.intervalLogger.Info("reset start agent")
				}
				oldStopSync = it.isCollectAndCronStop
			}
			
			it.intervalLogger.Debug("sync interval")
		}

		for name, interval := range it.intervals {
			if nowSec % interval == 0 {
				it.collectSend.Push(name)
			}
		}

		for name, interval := range it.crons {
			if nowSec % interval == 0 {
				it.cronSend.Push(name)
			}
		}

		select {
		case <-ctx.Done():
			isStop = true
		default:
		}

		time.Sleep(500 * time.Millisecond)
	}

	return nil
	
}