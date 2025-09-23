package types

import (
	"agent_common/pkg/applnew/loader"
	"agent_common/pkg/applnew/logger"
	"agent_common/pkg/util/types"
	"context"
)

type CronFn[PUSH any, FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]] func(context.Context, *FLAG, loader.ConfigureUpdater[FLAG, FLAGPTR], logger.LevelLogger) (*PUSH, error)
type CollectFn[PUSH any] func(context.Context, logger.LevelLogger) (*PUSH, error)

type DataPusher[PUSH any] interface {
	Push(*PUSH, context.Context, logger.LevelLogger) error
}