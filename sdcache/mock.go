package sdcache

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"time"
)

type mockEntry struct {
	expireAt time.Time
	value    any
}
type mockCache struct {
	cache map[string]mockEntry
	ttl   time.Duration
}

var _ Cache = &mockCache{}

func newMockCache(ttl time.Duration) *mockCache {
	if ttl < 0 {
		ttl = 0
	}
	return &mockCache{
		cache: make(map[string]mockEntry),
		ttl:   ttl,
	}
}

func (m *mockCache) Clear(ctx context.Context) error {
	m.cache = make(map[string]mockEntry)
	return nil
}

func (m *mockCache) Get(ctx context.Context, k any) (any, error) {
	k1 := k.(string)
	entry, ok := m.cache[k1]
	if ok {
		if entry.expireAt.IsZero() {
			return entry.value, nil
		} else {
			if time.Now().Before(entry.expireAt) {
				return entry.value, nil
			} else {
				delete(m.cache, k1)
				return nil, sderr.WithStack(ErrNotFound)
			}
		}
	} else {
		return nil, sderr.WithStack(ErrNotFound)
	}
}

func (m *mockCache) GetTTL(ctx context.Context, k any) (time.Duration, error) {
	k1 := k.(string)
	entry, ok := m.cache[k1]
	if ok {
		if entry.expireAt.IsZero() {
			return 0, nil
		} else {
			now := time.Now()
			if now.Before(entry.expireAt) {
				return entry.expireAt.Sub(now), nil
			} else {
				delete(m.cache, k1)
				return 0, sderr.WithStack(ErrNotFound)
			}
		}
	} else {
		return 0, sderr.WithStack(ErrNotFound)
	}
}

func (m *mockCache) Put(ctx context.Context, k, v any, opts *PutOptions) error {
	k1, ttl := k.(string), m.getTTL(opts)
	if ttl > 0 {
		m.cache[k1] = mockEntry{expireAt: time.Now().Add(ttl), value: v}
	} else {
		m.cache[k1] = mockEntry{expireAt: time.Time{}, value: v}
	}
	return nil
}

func (m *mockCache) Delete(ctx context.Context, k any) error {
	k1 := k.(string)
	delete(m.cache, k1)
	return nil
}

func (m *mockCache) GetOrPut(ctx context.Context, k any, loader func(ctx context.Context, k any) (any, error), opts *PutOptions) (any, error) {
	if loader == nil {
		return nil, sderr.New("no loader")
	}

	if m == nil {
		v, err := loader(ctx, k)
		if err != nil {
			return nil, sderr.Wrap(err, "load value for mock error")
		}
		if v == nil {
			return nil, sderr.Wrap(ErrNotFound, "load nothing")
		}
		return v, nil
	}

	k1, ttl := k.(string), m.getTTL(opts)

	entry, ok := m.cache[k1]
	if ok {
		if entry.expireAt.IsZero() {
			return entry.value, nil
		} else {
			if time.Now().Before(entry.expireAt) {
				return entry.value, nil
			} else {
				delete(m.cache, k1)
			}
		}
	}

	v, err := loader(ctx, k1)
	if err != nil {
		return nil, sderr.Wrap(err, "load value for mock error")
	}
	if v == nil {
		return nil, sderr.Wrap(ErrNotFound, "load nothing")
	}
	if ttl > 0 {
		m.cache[k1] = mockEntry{expireAt: time.Now().Add(ttl), value: v}
	} else {
		m.cache[k1] = mockEntry{expireAt: time.Time{}, value: v}
	}
	return v, nil
}

func (m *mockCache) getTTL(opts *PutOptions) time.Duration {
	if opts == nil || opts.TTL < 0 {
		return m.ttl
	}
	return opts.TTL
}
