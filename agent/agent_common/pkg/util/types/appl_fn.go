package types

import (
	"agent_common/pkg/logger"
	"context"
	"io"
)


type CollectFnRet struct {
	Tablename string
	Query     string
	Data [][]any
}

type IntervalRegister[DB any] struct {
	Name    string
	Fn      CollectFn[DB]
}

type CollectFn[DB any] func(targetConn DB, ctx context.Context) (*CollectFnRet, error)
type GenDbLoggerFn[DBCONF StrCloneToGetter] func(DBCONF) (logger.DbLogger, io.Closer, error)
type GenTargetDbPoolFn[DB any,DBCONF StrCloneToGetter] func(DBCONF) (TargetDbPool[DB], io.Closer, error)

type TargetDbPool[DB any] interface {
	GetDbConn(context.Context) (DB, error)
}