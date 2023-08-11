package sdcache

import (
	"github.com/gaorx/stardust5/sdtime"
	"testing"
)

func TestMock(t *testing.T) {
	ttlSecs := int64(2)
	c := newMockCache(sdtime.Seconds(ttlSecs))
	DoTestCommon(t, c)
	DoTestExpiration(t, c, ttlSecs)
}
