package types

import (
	"context"
	"io"
)

type CollectConnPool[CONN any] interface {
	GetDbConn(ctx context.Context) (CONN, error)
	io.Closer
}