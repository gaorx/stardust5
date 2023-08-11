package sdcache

import (
	"context"
	"time"
)

type Cache interface {
	Clear(ctx context.Context) error
	Get(ctx context.Context, k any) (any, error)
	GetTTL(ctx context.Context, k any) (time.Duration, error)
	Put(ctx context.Context, k, v any, opts *PutOptions) error
	Delete(ctx context.Context, k any) error
	GetOrPut(ctx context.Context, k any, loader func(ctx context.Context, k any) (any, error), opts *PutOptions) (any, error)
}

type PutOptions struct {
	TTL  time.Duration
	Cost int64
}
