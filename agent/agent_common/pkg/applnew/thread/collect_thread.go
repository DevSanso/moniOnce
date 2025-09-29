package thread

import (
	"agent_common/pkg/applnew/constants/context/key"
	"agent_common/pkg/applnew/logger"
	apptype "agent_common/pkg/applnew/types"
	"agent_common/pkg/util/types"
	"context"
	"io"
	"time"
)

type CollectThread[PUSH any, CONN io.Closer] struct {
	recv     types.Deque[apptype.IntervalData]
	pushSend types.Pusher[*PUSH]
	mapping  map[string]apptype.CollectFn[PUSH, CONN]

	collectConnPool apptype.CollectConnPool[CONN]
	collectLogger   logger.LevelLogger
}

func NewCollectThread[PUSH any, CONN io.Closer](recv types.Deque[apptype.IntervalData], conn apptype.CollectConnPool[CONN], pusher types.Pusher[*PUSH], mapping map[string]apptype.CollectFn[PUSH, CONN], logger logger.LevelLogger) CollectThread[PUSH, CONN] {
	return CollectThread[PUSH, CONN]{
		recv:            recv,
		pushSend:        pusher,
		mapping:         mapping,
		collectLogger:   logger,
		collectConnPool: conn,
	}
}

func (ct *CollectThread[PUSH, CONN]) Run(ctx context.Context) error {
	isStop := false
	for !isStop {
		data, _ := ct.recv.Pop()

		f, ok := ct.mapping[data.Name]

		if !ok {
			ct.collectLogger.Error("not support collect : ", data)
			continue
		}
		
		timeCtx, cancelFn := context.WithTimeout(ctx, time.Second * (time.Duration(data.Interval) - 1))
		
		conn, connErr := ct.collectConnPool.GetDbConn(timeCtx)
		if connErr != nil {
			ct.collectLogger.Error("conn get failed : ", connErr)
			cancelFn()
			continue
		}

		dataCtx := context.WithValue(context.WithValue(timeCtx, key.ContextNameKey, data.Name), key.ContextIntervalKey, data.Interval)

		ret, retErr := f(dataCtx, conn, ct.collectLogger)
		if retErr != nil {
			ct.collectLogger.Error("collectFn exec failed :", retErr.Error())
		}
		conn.Close()
		cancelFn()

		if ret != nil {
			ct.pushSend.Push(ret)
		} else {
			ct.collectLogger.Warn("collectFn exec is ret null, name: ", data)
		}

		if ok {
			select {
			case <-ctx.Done():
				isStop = true
			default:
			}

			time.Sleep(500 * time.Millisecond)
		}
	}

	return nil
}
