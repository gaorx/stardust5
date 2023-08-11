package sdredis

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"github.com/redis/go-redis/v9"
)

func ForEachShards(ctx context.Context, client redis.UniversalClient, action func(context.Context, *redis.Client) error) error {
	if c1, ok := client.(*redis.Client); ok {
		err := action(ctx, c1)
		return sderr.Wrap(err, "for each shard error")
	} else if c1, ok := client.(*redis.Ring); ok {
		err := c1.ForEachShard(ctx, action)
		return sderr.Wrap(err, "for each shard error (ring)")
	} else if c1, ok := client.(*redis.ClusterClient); ok {
		err := c1.ForEachShard(ctx, action)
		return sderr.Wrap(err, "for each shard error (cluster)")
	} else {
		panic(sderr.New("for each shards error"))
	}
}

func IsNotFound(err error) bool {
	return sderr.Is(err, redis.Nil)
}
