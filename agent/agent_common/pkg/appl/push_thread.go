package appl

import (
	"agent_common/pkg/logger"
	"agent_common/pkg/util/types"
	"context"
	"time"
)

type PushThread struct {
	recv_channel types.Deque[types.CollectFnRet]
	levelLogger logger.LevelLogger
	dbLogger    logger.ILogger
}

func newPushThread(recv types.Deque[types.CollectFnRet], log logger.LevelLogger, dbLogger  logger.ILogger) PushThread {
	return PushThread{recv_channel: recv, levelLogger: log, dbLogger: dbLogger}
}

func (pt *PushThread)Run(ctx context.Context) error {
	isStop := false
	
	for !isStop {
		data,_ := pt.recv_channel.Pop()
		pt.levelLogger.Debug("pop ", data.Key," table data")
		err := pt.dbLogger.Log(data.Key, data.Data)
		if err != nil {
			pt.levelLogger.Error("insert failed :", err.Error())
		}

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

