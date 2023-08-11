package sdcache

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDouble(t *testing.T) {
	ttlSecs1, ttlSecs2 := int64(2), int64(3)
	l1 := newMockCache(sdtime.Seconds(ttlSecs1))
	l2 := newMockCache(sdtime.Seconds(ttlSecs2))

	{
		c := D[*mockCache, *mockCache](l1, l2)
		DoTestCommon(t, c)
		DoTestExpiration(t, c, ttlSecs2)
	}

	{
		c := D[*mockCache, *mockCache](l1, nil)
		DoTestCommon(t, c)
		DoTestExpiration(t, c, ttlSecs1)
	}

	{
		c := D[*mockCache, *mockCache](nil, l2)
		DoTestCommon(t, c)
		DoTestExpiration(t, c, ttlSecs2)
	}

	{
		c := D[*mockCache, *mockCache](nil, nil)
		var errLoad = sderr.New("load error")
		loadCounter := 1
		loader := func(ctx context.Context, k any) (any, error) {
			loadCounter += 1
			if k == "k1" {
				return "v1", nil
			} else {
				return nil, sderr.WithStack(errLoad)
			}
		}
		v2, err := c.GetOrPut(context.Background(), "k2", loader, nil)
		assert.ErrorIs(t, err, errLoad)
		assert.Nil(t, v2)
		assert.Equal(t, 2, loadCounter)
		v1, err := c.GetOrPut(context.Background(), "k1", loader, nil)
		assert.NoError(t, err)
		assert.Equal(t, "v1", v1)
		assert.Equal(t, 3, loadCounter)
		v1, err = c.GetOrPut(context.Background(), "k1", loader, nil)
		assert.NoError(t, err)
		assert.Equal(t, "v1", v1)
		assert.Equal(t, 4, loadCounter)

	}
}
