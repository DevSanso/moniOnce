package thread

import (
	"agent_common/pkg/applnew/loader"
	"agent_common/pkg/applnew/logger"
	apptype "agent_common/pkg/applnew/types"
	"agent_common/pkg/util/types"
	"context"
	"time"
)


type CronThread[PUSH any,FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]] struct {
	recv types.Deque[string]
	pushSend types.Pusher[*PUSH]
	confLoader loader.Configure[apptype.ApplConfData, apptype.AppSyncData, FLAG, *apptype.ApplConfData, *apptype.AppSyncData, FLAGPTR]
	mapping map[string]apptype.CronFn[PUSH, FLAG, FLAGPTR]

	cronLogger logger.LevelLogger
}

func NewCronThread[PUSH any,FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]](
	recv types.Deque[string], 
	pusher types.Pusher[*PUSH],
	mapping map[string]apptype.CronFn[PUSH, FLAG, FLAGPTR],
	logger logger.LevelLogger,
	confLoader loader.Configure[apptype.ApplConfData, apptype.AppSyncData, FLAG, *apptype.ApplConfData, *apptype.AppSyncData, FLAGPTR]) CronThread[PUSH, FLAG, FLAGPTR] {
	return CronThread[PUSH, FLAG, FLAGPTR]{
		recv: recv,
		pushSend: pusher,
		mapping: mapping,
		confLoader: confLoader,
		cronLogger: logger,
	}
}

func (ct *CronThread[PUSH, FLAG, FLAGPTR])Run(ctx context.Context) error {
	isStop := false
	for !isStop {

		data,_ := ct.recv.Pop()
		f, ok := ct.mapping[data]
		
		if ok {
			if flag, err := ct.confLoader.LoadFlag(); err != nil {
				ct.cronLogger.Error("load flag is failed : ", err.Error())
			} else {
				ret, retErr := f(ctx, flag, ct.confLoader, ct.cronLogger)
				if retErr != nil {
					ct.cronLogger.Error("cronFn exec failed :", retErr.Error())
				}

				if ret != nil {
					ct.pushSend.Push(ret)
				}
				
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