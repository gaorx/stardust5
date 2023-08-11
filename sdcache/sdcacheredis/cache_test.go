package sdcacheredis

import (
	"github.com/gaorx/stardust5/sdcache"
	"github.com/gaorx/stardust5/sdredis"
	"github.com/gaorx/stardust5/sdtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache(t *testing.T) {
	ttlSecs := int64(2)
	c, err := Dial(sdredis.Address{
		Addrs:    []string{"localhost:6379"},
		Password: "123456",
	}, Config{
		Key:     sdcache.StrKey{},
		Encoder: sdcache.TextEncoder{},
		TTL:     sdtime.Seconds(ttlSecs),
	})
	assert.NoError(t, err)

	// go
	sdcache.DoTestCommon(t, c)
	sdcache.DoTestExpiration(t, c, ttlSecs)
}
