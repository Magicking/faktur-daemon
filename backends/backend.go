package backends

import (
	"context"
)

type Backend interface {
	Init(ctx context.Context) error
	Run(ctx context.Context, ctl chan struct{}) error
}
