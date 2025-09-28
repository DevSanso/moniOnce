package pusher

import "context"

type IConnPusher interface {
	Exec(ctx context.Context, cmd string, args ...any) error
}