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
	dbLogger    logger.DbLogger
}

func newPushThread(recv types.Deque[types.CollectFnRet], log logger.LevelLogger, dbLogger  logger.DbLogger) PushThread {
	return PushThread{recv_channel: recv, levelLogger: log, dbLogger: dbLogger}
}

func (pt *PushThread)Run(ctx context.Context) error {
	isStop := false
	
	for !isStop {
		data,_ := pt.recv_channel.Pop()
		pt.levelLogger.Debug("pop ", data.Tablename," table data")
		err := pt.dbLogger.Exec(data.Query, data.Data)
		if err != nil {
			pt.levelLogger.Error("insert failed :", err.Error())
		}

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

