package thread

import (
	"agent_common/pkg/applnew/logger"
	appltype "agent_common/pkg/applnew/types"
	"agent_common/pkg/util/types"
	"context"
	"time"
)

type PushThread[PUSH any] struct {
	objectId int
	recv types.Deque[*PUSH]
	pushLogger logger.LevelLogger
	pusher appltype.DataPusher[PUSH]
}

func NewPushThread[PUSH any](objectId int, recv types.Deque[*PUSH], logger logger.LevelLogger, pusher appltype.DataPusher[PUSH]) PushThread[PUSH] {
	return PushThread[PUSH]{
		recv : recv,
		pushLogger: logger,
		pusher: pusher,
		objectId: objectId,
	}
}

func(pt *PushThread[PUSH])Run(ctx context.Context) error {
	isStop := false

	for !isStop {
		data,_ :=  pt.recv.Pop()

		if err := pt.pusher.Push(pt.objectId, data, ctx, pt.pushLogger); err != nil {
			pt.pushLogger.Error("Push error : ", err.Error())
		}

		select {
		case <-ctx.Done():
			isStop = true
		default:
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
	
}