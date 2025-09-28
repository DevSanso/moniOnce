package types

import (
	"agent_common/pkg/applnew/loader"
	"agent_common/pkg/applnew/logger"
	"agent_common/pkg/util/types"
	"context"
	"io"
)

type CronFn[PUSH any, FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]] func(context.Context, *FLAG, loader.ConfigureUpdater[FLAG, FLAGPTR], logger.LevelLogger) (*PUSH, error)
type CollectFn[PUSH any, CONN io.Closer] func(context.Context, CONN, logger.LevelLogger) (*PUSH, error)

type DataPusher[PUSH any] interface {
	Push(*PUSH, context.Context, logger.LevelLogger) error
	io.Closer
}

type GenCollectConnPoolFn[CONN io.Closer] func(url string, args ...any) (CollectConnPool[CONN], error)
type GenPusherFn[PUSH any] func(url string, args ...any) (DataPusher[PUSH], error)