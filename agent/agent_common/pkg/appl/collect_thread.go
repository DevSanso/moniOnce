package appl

import (
	"agent_common/pkg/constant"
	"agent_common/pkg/logger"
	"agent_common/pkg/util/types"
	"context"
	"time"
)

type CollectThread[DB any] struct {
	send_channel  types.Pusher[types.CollectFnRet]
	recv_channel  types.Deque[types.CollectFn[DB]]
	levelLogger   logger.LevelLogger
	tartgetDbPool types.TargetDbPool[DB]

	applCustomConfig map[string]string
}

func newCollectThread[DB any](send types.Pusher[types.CollectFnRet], recv types.Deque[types.CollectFn[DB]],
	 log logger.LevelLogger, target types.TargetDbPool[DB], applCustomConfig map[string]string) CollectThread[DB] {

	return CollectThread[DB]{send_channel: send, recv_channel: recv, levelLogger: log, tartgetDbPool: target, applCustomConfig : applCustomConfig}
}

func (ct *CollectThread[DB]) GetConnCtx() (context.Context, context.CancelFunc) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second * 10)
	retCtx := context.WithValue(ctx, constant.ApplCustomConfigCtxKey, ct.applCustomConfig)
	return retCtx, cancelFn
}

func (ct *CollectThread[DB]) Run(ctx context.Context) error {
	isStop := false

	for !isStop {
		fn, _ := ct.recv_channel.Pop()

		connCtx, ctxCancelFn := ct.GetConnCtx()
		
		conn, connErr := ct.tartgetDbPool.GetDbConn(connCtx)

		if connErr != nil {
			ctxCancelFn()
			ct.levelLogger.Error("Get Connection Failed : ", connErr.Error())
		}

		datas, fnErr := fn(conn, connCtx)
		ctxCancelFn()

		if fnErr == nil {
			if datas != nil {
				ct.send_channel.Push(*datas)
			} else {
				ct.levelLogger.Debug("Collect Ret is Null")
			}
			
		} else {
			ct.levelLogger.Error("Get Connection Failed : ", fnErr.Error())
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
