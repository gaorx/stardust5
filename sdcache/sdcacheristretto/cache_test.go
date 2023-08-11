package sdcacheristretto

import (
	"github.com/dgraph-io/ristretto"
	"github.com/gaorx/stardust5/sdcache"
	"github.com/gaorx/stardust5/sdtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache(t *testing.T) {
	ttlSecs := int64(2)
	c, err := NewByRistrettoConfig(
		ristretto.Config{
			NumCounters: 100,
			MaxCost:     100,
			BufferItems: 64,
		},
		Config{
			TTL: sdtime.Seconds(ttlSecs),
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, c)

	// go
	sdcache.DoTestCommon(t, c)
	sdcache.DoTestExpiration(t, c, ttlSecs)
}
