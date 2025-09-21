package appl

import (
	"agent_common/pkg/logger"
	"agent_common/pkg/util/types"
	"context"
	"time"
)

type InternvalThread[DB any] struct {
	intervals    map[int][]types.IntervalRegister[DB]
	send_channel types.Pusher[types.CollectFn[DB]]

	levelLogger logger.LevelLogger
}

func newIntervalThread[DB any](registers map[int][]types.IntervalRegister[DB], queue types.Pusher[types.CollectFn[DB]], levelLogger logger.LevelLogger) InternvalThread[DB] {
	return InternvalThread[DB]{
		registers,
		queue,
		levelLogger,
	}
}

func (ih *InternvalThread[DB]) Run(ctx context.Context) error {
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
