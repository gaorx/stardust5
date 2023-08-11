package sdconcur

import (
	"github.com/samber/lo"
	"sync"
)

func DoFuncs(n int, actions []func()) error {
	numActions := len(actions)
	if numActions == 0 {
		return nil
	}
	if n <= 0 {
		var wg sync.WaitGroup
		for _, f := range actions {
			wg.Add(1)
			go func(f func()) {
				defer wg.Done()
				_ = lo.Try0(f)
			}(f)
		}
		wg.Wait()
		return nil
	} else {
		if n > numActions {
			n = numActions
		}
		pool, err := NewPool(n, &PoolOptions{
			PreAlloc: true,
		})
		if err != nil {
			return err
		}
		defer func() { _ = pool.Close() }()
		var wg sync.WaitGroup
		for _, f := range actions {
			f1 := f
			wg.Add(1)
			err := pool.Submit(func() {
				defer wg.Done()
				_ = lo.Try0(f1)
			})
			if err != nil {
				return err
			}
		}
		wg.Wait()
		return nil
	}
}

func Do[T any](n int, l []T, action func(int, T)) error {
	return DoFuncs(n, Bind(l, action))
}

func Bind[T any](l []T, action func(int, T)) []func() {
	if action == nil {
		action = func(int, T) {}
	}
	actions := make([]func(), 0, len(l))
	for i, v := range l {
		i0, v0 := i, v
		actions = append(actions, func() {
			action(i0, v0)
		})
	}
	return actions
}
