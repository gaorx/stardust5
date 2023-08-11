package sdcacheredis

import (
	"context"
	"github.com/gaorx/stardust5/sdcache"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdredis"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache struct {
	client redis.UniversalClient
	config Config
}

type Config struct {
	Key     sdcache.Key
	Encoder sdcache.Encoder
	TTL     time.Duration
}

var _ sdcache.Cache = &Cache{}

func New(client redis.UniversalClient, config Config) (*Cache, error) {
	if client == nil {
		return nil, sderr.New("nil redis client")
	}
	return &Cache{client: client, config: config.trim()}, nil
}

func Dial(addr sdredis.Address, config Config) (*Cache, error) {
	client, err := sdredis.Dial(addr)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return New(client, config)
}

func (c *Cache) Client() redis.UniversalClient {
	return c.client
}

func (c *Cache) Config() Config {
	return c.config
}

func (c *Cache) SetKey(key sdcache.Key) *Cache {
	if key != nil {
		c.config.Key = key
	}
	return c
}

func (c *Cache) SetEncoder(encoder sdcache.Encoder) *Cache {
	if encoder != nil {
		c.config.Encoder = encoder
	}
	return c
}

func (c *Cache) SetTTL(ttl time.Duration) *Cache {
	if ttl >= 0 {
		c.config.TTL = ttl
	}
	return c
}

func (c *Cache) Clear(ctx context.Context) error {
	if err := c.checkConfig(true, false); err != nil {
		return err
	}
	client, prefix := c.client, c.config.Key.PrefixForClear()
	if prefix != "" {
		var cursor uint64 = 0
		for {
			keys, nextCursor, err := client.Scan(ctx, cursor, prefix+"*", 1).Result()
			if err != nil {
				return sderr.Wrap(err, "scan redis key error")
			}
			if len(keys) <= 0 {
				break
			}
			err = client.Del(ctx, keys...).Err()
			if err != nil {
				return sderr.Wrap(err, "delete redis key error")
			}
			cursor = nextCursor
		}
	} else {
		err := client.FlushAll(ctx).Err()
		if err != nil {
			return sderr.Wrap(err, "clear redis data error")
		}
	}
	return nil
}

func (c *Cache) Get(ctx context.Context, k any) (any, error) {
	if err := c.checkConfig(true, true); err != nil {
		return nil, err
	}
	client, key, encoder := c.client, c.config.Key, c.config.Encoder
	redisKey, err := key.EncodeKey(k)
	if err != nil {
		return nil, sderr.Wrap(err, "encode redis key error")
	}
	redisVal, err := client.Get(ctx, redisKey).Bytes()
	if err != nil {
		if sderr.Is(err, redis.Nil) {
			return nil, sderr.Wrap(sdcache.ErrNotFound, "get redis key error")
		} else {
			return nil, sderr.Wrap(err, "get redis value error")
		}
	}
	v, err := encoder.DecodeValue(redisVal)
	if err != nil {
		return nil, sderr.Wrap(err, "decode redis value error")
	}
	return v, nil
}

func (c *Cache) GetTTL(ctx context.Context, k any) (time.Duration, error) {
	if err := c.checkConfig(true, false); err != nil {
		return 0, err
	}
	client, key := c.client, c.config.Key
	redisKey, err := key.EncodeKey(k)
	if err != nil {
		return 0, sderr.Wrap(err, "encode redis key error")
	}
	ttl, err := client.TTL(ctx, redisKey).Result()
	if err != nil {
		return 0, sderr.Wrap(err, "get redis key ttl error")
	}
	if ttl == -1 {
		return 0, nil
	} else if ttl == -2 {
		return 0, sderr.Wrap(sdcache.ErrNotFound, "get redis key ttl error")
	} else {
		return ttl, nil
	}
}

func (c *Cache) Put(ctx context.Context, k, v any, opts *sdcache.PutOptions) error {
	if err := c.checkConfig(true, true); err != nil {
		return err
	}
	client, key, encoder := c.client, c.config.Key, c.config.Encoder
	redisKey, err := key.EncodeKey(k)
	if err != nil {
		return sderr.Wrap(err, "encode redis key error")
	}
	redisVal, err := encoder.EncodeValue(k, v)
	if err != nil {
		return sderr.Wrap(err, "encode redis value error")
	}
	err = client.Set(ctx, redisKey, redisVal, c.getTTL(opts)).Err()
	if err != nil {
		return sderr.Wrap(err, "set redis value error")
	}
	return nil
}

func (c *Cache) Delete(ctx context.Context, k any) error {
	if err := c.checkConfig(true, false); err != nil {
		return err
	}
	client, key := c.client, c.config.Key
	redisKey, err := key.EncodeKey(k)
	if err != nil {
		return sderr.Wrap(err, "encode redis key error")
	}
	err = client.Del(ctx, redisKey).Err()
	if err != nil {
		return sderr.Wrap(err, "delete redis key error")
	}
	return nil
}

func (c *Cache) GetOrPut(ctx context.Context, k any, loader func(ctx context.Context, k any) (any, error), opts *sdcache.PutOptions) (any, error) {
	if loader == nil {
		return nil, sderr.New("nil loader")
	}
	if c == nil {
		v, err := loader(ctx, k)
		if err != nil {
			return nil, sderr.Wrap(err, "load value for redis error")
		}
		if v == nil {
			return nil, sderr.Wrap(sdcache.ErrNotFound, "load nothing")
		}
		return v, nil
	}
	if err := c.checkConfig(true, true); err != nil {
		return nil, err
	}
	client, key, encoder := c.client, c.config.Key, c.config.Encoder
	redisKey, err := key.EncodeKey(k)
	if err != nil {
		return nil, sderr.Wrap(err, "encode redis key error")
	}
	redisVal, err := client.Get(ctx, redisKey).Bytes()
	if err != nil {
		if sderr.Is(err, redis.Nil) {
			v, err := loader(ctx, k)
			if err != nil {
				return nil, sderr.Wrap(err, "load value for redis error")
			}
			if v == nil {
				return nil, sderr.Wrap(sdcache.ErrNotFound, "load nothing")
			}
			redisVal, err := encoder.EncodeValue(k, v)
			if err != nil {
				return nil, sderr.Wrap(err, "encode redis value error")
			}
			err = client.Set(ctx, redisKey, redisVal, c.getTTL(opts)).Err()
			if err != nil {
				return nil, sderr.Wrap(err, "set redis value error")
			}
			return v, nil
		} else {
			return nil, sderr.Wrap(err, "get redis value error")
		}
	} else {
		v, err := encoder.DecodeValue(redisVal)
		if err != nil {
			return nil, sderr.Wrap(err, "decode redis value error")
		}
		return v, nil
	}
}

func (c *Cache) checkConfig(forKey, forEncoder bool) error {
	if forKey {
		if c.config.Key == nil {
			return sderr.New("nil key")
		}
	}
	if forEncoder {
		if c.config.Encoder == nil {
			return sderr.New("nil encoder")
		}
	}
	return nil
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
	return config
}
