package sdcacheristretto

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/gaorx/stardust5/sdcache"
	"github.com/gaorx/stardust5/sderr"
	"time"
)

type Cache struct {
	cache  *ristretto.Cache
	config Config
}

type Config struct {
	Key  sdcache.Key
	TTL  time.Duration
	Cost int64
}

var _ sdcache.Cache = &Cache{}

func New(c *ristretto.Cache, config Config) (*Cache, error) {
	if c == nil {
		return nil, sderr.New("nil ristretto cache")
	}
	return &Cache{cache: c, config: config.trim()}, nil
}

func NewByRistrettoConfig(ristrettoConfig ristretto.Config, config Config) (*Cache, error) {
	c, err := ristretto.NewCache(&ristrettoConfig)
	if err != nil {
		return nil, sderr.Wrap(err, "new ristretto cache error")
	}
	return New(c, config)
}

func (c *Cache) Ristretto() *ristretto.Cache {
	return c.cache
}

func (c *Cache) SetKey(key sdcache.Key) *Cache {
	if key != nil {
		c.config.Key = key
	}
	return c
}

func (c *Cache) SetTTL(ttl time.Duration) *Cache {
	if ttl >= 0 {
		c.config.TTL = ttl
	}
	return c
}

func (c *Cache) SetCost(cost int64) *Cache {
	if cost >= 0 {
		c.config.Cost = cost
	}
	return c
}

func (c *Cache) Clear(ctx context.Context) error {
	c.cache.Clear()
	return nil
}

func (c *Cache) Get(ctx context.Context, k any) (any, error) {
	ristrettoKey, err := c.encodeKey(k)
	if err != nil {
		return nil, sderr.Wrap(err, "encode ristretto key error")
	}
	v, ok := c.cache.Get(ristrettoKey)
	if !ok {
		return nil, sderr.Wrap(sdcache.ErrNotFound, "get ristretto key error")
	}
	return v, nil
}

func (c *Cache) GetTTL(ctx context.Context, k any) (time.Duration, error) {
	ristrettoKey, err := c.encodeKey(k)
	if err != nil {
		return 0, sderr.Wrap(err, "encode ristretto key error")
	}
	ttl, ok := c.cache.GetTTL(ristrettoKey)
	if !ok {
		return 0, sderr.Wrap(sdcache.ErrNotFound, "get ristretto key ttl error")
	}
	return ttl, nil
}

func (c *Cache) Put(ctx context.Context, k, v any, opts *sdcache.PutOptions) error {
	ttl, cost := c.getTTL(opts), c.getCost(opts)

	ristrettoKey, err := c.encodeKey(k)
	if err != nil {
		return sderr.Wrap(err, "encode ristretto key error")
	}

	if ttl > 0 {
		_ = c.cache.SetWithTTL(ristrettoKey, v, cost, ttl)
		c.cache.Wait()
	} else {
		_ = c.cache.Set(ristrettoKey, v, cost)
		c.cache.Wait()
	}

	return nil
}

func (c *Cache) Delete(ctx context.Context, k any) error {
	ristrettoKey, err := c.encodeKey(k)
	if err != nil {
		return sderr.Wrap(err, "encode ristretto key error")
	}
	c.cache.Del(ristrettoKey)
	c.cache.Wait()
	return nil
}

func (c *Cache) GetOrPut(ctx context.Context, k any, loader func(ctx context.Context, k any) (any, error), opts *sdcache.PutOptions) (any, error) {
	if loader == nil {
		return nil, sderr.New("nil loader")
	}

	if c == nil {
		v, err := loader(ctx, k)
		if err != nil {
			return nil, sderr.Wrap(err, "load value for ristretto error")
		}
		if v == nil {
			return nil, sderr.Wrap(sdcache.ErrNotFound, "load nothing")
		}
		return v, nil
	}

	ristrettoKey, err := c.encodeKey(k)
	if err != nil {
		return nil, sderr.Wrap(err, "encode ristretto key error")
	}

	v, ok := c.cache.Get(ristrettoKey)
	if ok {
		return v, nil
	}

	v, err = loader(ctx, k)
	if err != nil {
		return nil, sderr.Wrap(err, "load value for ristretto error")
	}
	if v == nil {
		return nil, sderr.Wrap(sdcache.ErrNotFound, "load nothing")
	}

	ttl, cost := c.getTTL(opts), c.getCost(opts)
	var set bool
	if ttl > 0 {
		set = c.cache.SetWithTTL(ristrettoKey, v, cost, ttl)
		c.cache.Wait()
	} else {
		set = c.cache.Set(ristrettoKey, v, cost)
		c.cache.Wait()
	}
	if !set {
		return nil, sderr.New("put ristretto key error")
	}
	return v, nil
}

func (c *Cache) encodeKey(k any) (any, error) {
	if c.config.Key != nil {
		return c.config.Key.EncodeKey(k)
	} else {
		return k, nil
	}
}

func (c *Cache) getCost(opts *sdcache.PutOptions) int64 {
	if opts == nil || opts.Cost < 0 {
		return c.config.Cost
	}
	return opts.Cost
}

func (c *Cache) getTTL(opts *sdcache.PutOptions) time.Duration {
	if opts == nil || opts.TTL < 0 {
		return c.config.TTL
	}
	return opts.TTL
}

func (config Config) trim() Config {
	if config.TTL < 0 {
		config.TTL = 0
	}
	if config.Cost < 0 {
		config.Cost = 0
	}
	return config
}
