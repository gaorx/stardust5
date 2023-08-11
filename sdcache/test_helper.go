package sdcache

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func DoTestCommon(t *testing.T, c Cache) {
	// clear
	err := c.Clear(context.Background())
	assert.NoError(t, err)

	// Get
	v1, err := c.Get(context.Background(), "k1")
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, v1)

	// Put
	err = c.Put(context.Background(), "k1", "v1", nil)
	assert.NoError(t, err)

	// Get
	v1, err = c.Get(context.Background(), "k1")
	assert.NoError(t, err)
	assert.Equal(t, "v1", v1)

	// Delete
	err = c.Delete(context.Background(), "k1")
	assert.NoError(t, err)

	// Get
	v1, err = c.Get(context.Background(), "k1")
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, v1)

	// GetOrPut
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
	v2, err = c.Get(context.Background(), "k2")
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, v2)
	v1, err = c.GetOrPut(context.Background(), "k1", loader, nil)
	assert.NoError(t, err)
	assert.Equal(t, "v1", v1)
	assert.Equal(t, 3, loadCounter)
	v1, err = c.Get(context.Background(), "k1")
	assert.NoError(t, err)
	assert.Equal(t, "v1", v1)
	v1, err = c.GetOrPut(context.Background(), "k1", loader, nil)
	assert.NoError(t, err)
	assert.Equal(t, "v1", v1)
	assert.Equal(t, 3, loadCounter)
}

func DoTestExpiration(t *testing.T, c Cache, ttlSecs int64) {
	assert.True(t, ttlSecs > 1)

	// clear
	err := c.Clear(context.Background())
	assert.NoError(t, err)

	// Get
	v1, err := c.Get(context.Background(), "k1")
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, v1)

	// Put
	err = c.Put(context.Background(), "k1", "v1", &PutOptions{TTL: sdtime.Seconds(ttlSecs + 1)})
	assert.NoError(t, err)
	sdtime.SleepS(1)
	v1, err = c.Get(context.Background(), "k1")
	assert.NoError(t, err)
	assert.Equal(t, "v1", v1)
	sdtime.SleepS(ttlSecs + 1)
	_, err = c.Get(context.Background(), "k1")
	assert.ErrorIs(t, err, ErrNotFound)

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
	v2, err = c.Get(context.Background(), "k2")
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, v2)
	v1, err = c.GetOrPut(context.Background(), "k1", loader, nil)
	assert.NoError(t, err)
	assert.Equal(t, "v1", v1)
	assert.Equal(t, 3, loadCounter)
	sdtime.SleepS(1)
	v1, err = c.GetOrPut(context.Background(), "k1", loader, nil)
	assert.NoError(t, err)
	assert.Equal(t, "v1", v1)
	assert.Equal(t, 3, loadCounter)
	sdtime.SleepS(ttlSecs)
	v1, err = c.Get(context.Background(), "k1")
	assert.ErrorIs(t, err, ErrNotFound)
	v1, err = c.GetOrPut(context.Background(), "k1", loader, &PutOptions{TTL: sdtime.Seconds(2)})
	assert.NoError(t, err)
	assert.Equal(t, "v1", v1)
	assert.Equal(t, 4, loadCounter)
}
