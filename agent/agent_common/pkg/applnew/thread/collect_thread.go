package thread

import (
	"agent_common/pkg/applnew/logger"
	apptype "agent_common/pkg/applnew/types"
	"agent_common/pkg/util/types"
	"context"
	"time"
)

type CollectThread[PUSH any] struct {
	recv types.Deque[string]
	pushSend types.Pusher[*PUSH]
	mapping map[string]apptype.CollectFn[PUSH]

	collectLogger logger.LevelLogger
}

func NewCollectThread[PUSH any](recv types.Deque[string], pusher types.Pusher[*PUSH], mapping map[string]apptype.CollectFn[PUSH], logger logger.LevelLogger) CollectThread[PUSH] {
	return CollectThread[PUSH]{
		recv: recv,
		pushSend: pusher,
		mapping: mapping,
		collectLogger: logger,
	}
}

func (ct *CollectThread[PUSH])Run(ctx context.Context) error {
	isStop := false
	for !isStop {
		data,_ := ct.recv.Pop()

		f, ok := ct.mapping[data]
		ret, retErr := f(ctx, ct.collectLogger)
		if retErr != nil {
			ct.collectLogger.Error("collectFn exec failed :", retErr.Error())
		}

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