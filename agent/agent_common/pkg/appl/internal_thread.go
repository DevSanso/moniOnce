package appl

import (
	"agent_common/pkg/logger"
	"agent_common/pkg/util/types"
	"context"
	"time"
)

type InternvalThread[CONN any] struct {
	intervals    map[int][]types.IntervalRegister[CONN]
	send_channel types.Pusher[types.CollectFn[CONN]]

	levelLogger logger.LevelLogger
}

func newIntervalThread[CONN any](registers map[int][]types.IntervalRegister[CONN], queue types.Pusher[types.CollectFn[CONN]], levelLogger logger.LevelLogger) InternvalThread[CONN] {
	return InternvalThread[CONN]{
		registers,
		queue,
		levelLogger,
	}
}

func (ih *InternvalThread[CONN]) Run(ctx context.Context) error {
	isStop := false

	for !isStop {
		nowSec := time.Now().Second() % 60

		list := ih.intervals[nowSec]

		for _, fn := range list {
			ih.send_channel.Push(fn.Fn)
			ih.levelLogger.Debug("push collect fn : ", fn.Name)
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
