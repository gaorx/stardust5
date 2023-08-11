package sdcache

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTyped(t *testing.T) {
	type value struct {
		Id   string `json:"_id"`
		Name string `json:"_name"`
	}

	ttlSecs := int64(2)
	c := newMockCache(sdtime.Seconds(ttlSecs))
	tc := T[string, *value](c)
	err := tc.Clear(context.Background())
	assert.NoError(t, err)
	errLoad := sderr.New("load error")
	loadCounter := 1
	loader := func(ctx context.Context, k string) (*value, error) {
		loadCounter += 1
		if k == "k1" {
			return &value{Id: "k1", Name: "v1"}, nil
		} else {
			return nil, errLoad
		}
	}
	v2, err := tc.GetOrPut(context.Background(), "k2", loader, nil)
	assert.ErrorIs(t, err, errLoad)
	assert.Nil(t, v2)
	assert.Equal(t, 2, loadCounter)
	v1, err := tc.GetOrPut(context.Background(), "k1", loader, nil)
	assert.NoError(t, err)
	assert.NotNil(t, v1)
	assert.Equal(t, value{Id: "k1", Name: "v1"}, *v1)
	assert.Equal(t, 3, loadCounter)
}
