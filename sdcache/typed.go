package sdcache

import (
	"context"
	"time"
)

type Typed[K, V any] struct {
	C Cache
}

func T[K, V any](c Cache) Typed[K, V] {
	return Typed[K, V]{C: c}
}

func (t Typed[K, V]) Clear(ctx context.Context) error {
	return t.C.Clear(ctx)
}

func (t Typed[K, V]) Get(ctx context.Context, k K) (V, error) {
	v, err := t.C.Get(ctx, k)
	if err != nil {
		var empty V
		return empty, err
	}
	return v.(V), nil
}

func (t Typed[K, V]) GetTTL(ctx context.Context, k K) (time.Duration, error) {
	return t.C.GetTTL(ctx, k)
}

func (t Typed[K, V]) Put(ctx context.Context, k K, v V, opts *PutOptions) error {
	return t.C.Put(ctx, k, v, opts)
}

func (t Typed[K, V]) Delete(ctx context.Context, k K) error {
	return t.C.Delete(ctx, k)
}

func (t Typed[K, V]) GetOrPut(ctx context.Context, k K, loader func(ctx context.Context, k K) (V, error), opts *PutOptions) (V, error) {
	v, err := t.C.GetOrPut(ctx, k, func(ctx context.Context, k0 any) (any, error) {
		v0, err := loader(ctx, k0.(K))
		if err != nil {
			return nil, err
		}
		return v0, nil
	}, opts)
	if err != nil {
		var empty V
		return empty, err
	}
	return v.(V), nil
}
