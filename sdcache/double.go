package sdcache

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"time"
)

type Double[L1, L2 interface {
	comparable
	Cache
}] struct {
	L1 L1
	L2 L2
}

var _ Cache = Double[*mockCache, *mockCache]{}

func D[L1, L2 interface {
	comparable
	Cache
}](c1 L1, c2 L2) Double[L1, L2] {
	return Double[L1, L2]{L1: c1, L2: c2}
}

func (d Double[L1, L2]) Clear(ctx context.Context) error {
	c1, c2 := d.L1, d.L2
	if d.hasL1() && d.hasL2() {
		err2 := c2.Clear(ctx)
		err1 := c1.Clear(ctx)
		return combineErr(err2, err1)
	} else if d.hasL1() {
		return c1.Clear(ctx)
	} else if d.hasL2() {
		return c2.Clear(ctx)
	} else {
		return sderr.New("double cache is empty")
	}
}

func (d Double[L1, L2]) Get(ctx context.Context, k any) (any, error) {
	c1, c2 := d.L1, d.L2
	if d.hasL1() && d.hasL2() {
		return c1.GetOrPut(ctx, k, func(ctx context.Context, k any) (any, error) {
			return c2.Get(ctx, k)
		}, nil)
	} else if d.hasL1() {
		return c1.Get(ctx, k)
	} else if d.hasL2() {
		return c2.Get(ctx, k)
	} else {
		return nil, sderr.New("double cache is empty")
	}
}

func (d Double[L1, L2]) GetTTL(ctx context.Context, k any) (time.Duration, error) {
	c1, c2 := d.L1, d.L2
	if d.hasL1() {
		return c1.GetTTL(ctx, k)
	} else if d.hasL2() {
		return c2.GetTTL(ctx, k)
	} else {
		return 0, sderr.New("double cache is empty")
	}
}

func (d Double[L1, L2]) Put(ctx context.Context, k, v any, opts *PutOptions) error {
	c1, c2 := d.L1, d.L2
	if d.hasL1() && d.hasL2() {
		err2 := c2.Put(ctx, k, v, opts)
		err1 := c1.Put(ctx, k, v, nil)
		return combineErr(err1, err2)
	} else if d.hasL1() {
		return c1.Put(ctx, k, v, opts)
	} else if d.hasL2() {
		return c2.Put(ctx, k, v, nil)
	} else {
		return sderr.New("double cache is empty")
	}
}

func (d Double[L1, L2]) Delete(ctx context.Context, k any) error {
	c1, c2 := d.L1, d.L2
	if d.hasL1() && d.hasL2() {
		err2 := c2.Delete(ctx, k)
		err1 := c1.Delete(ctx, k)
		return combineErr(err1, err2)
	} else if d.hasL1() {
		return c1.Delete(ctx, k)
	} else if d.hasL2() {
		return c2.Delete(ctx, k)
	} else {
		return sderr.New("double cache is empty")
	}
}

func (d Double[L1, L2]) GetOrPut(ctx context.Context, k any, loader func(ctx context.Context, k any) (any, error), opts *PutOptions) (any, error) {
	c1, c2 := d.L1, d.L2
	if d.hasL1() && d.hasL2() {
		return c1.GetOrPut(ctx, k, func(ctx context.Context, k any) (any, error) {
			return c2.GetOrPut(ctx, k, loader, opts)
		}, nil)
	} else if d.hasL1() {
		return c1.GetOrPut(ctx, k, loader, nil)
	} else if d.hasL2() {
		return c2.GetOrPut(ctx, k, loader, opts)
	} else {
		if loader == nil {
			return nil, sderr.New("nil loader")
		}
		v, err := loader(ctx, k)
		if err != nil {
			return nil, sderr.Wrap(err, "load value for double cache error")
		}
		if v == nil {
			return nil, sderr.Wrap(ErrNotFound, "load nothing")
		}
		return v, nil
	}
}

func (d Double[L1, L2]) hasL1() bool {
	var empty1 L1
	return d.L1 != empty1
}

func (d Double[L1, L2]) hasL2() bool {
	var empty2 L2
	return d.L2 != empty2
}

func combineErr(err1, err2 error) error {
	if err1 != nil && err2 != nil {
		return sderr.Combine([]error{err1, err2})
	} else if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	} else {
		return nil
	}
}
