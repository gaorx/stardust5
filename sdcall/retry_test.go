package sdcall

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetry(t *testing.T) {
	{
		c := 0
		err := Retry(0, func() error {
			c++
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, c)
	}

	{
		c := 0
		err := Retry(0, func() error {
			c++
			return sderr.New("xx")
		})
		assert.Error(t, err)
		assert.Equal(t, 1, c)
	}

	{
		c := 0
		err := Retry(10, func() error {
			c++
			if c >= 3 {
				return nil
			} else {
				return sderr.New("xx")
			}
		})
		assert.NoError(t, err)
		assert.Equal(t, 3, c)
	}

	{
		c := 0
		err := Retry(10, func() error {
			c++
			if c >= 30 {
				return nil
			} else {
				return sderr.New("xx")
			}
		})
		assert.Error(t, err)
		assert.Equal(t, 11, c)
	}
}
