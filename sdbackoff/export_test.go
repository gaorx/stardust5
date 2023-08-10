package sdbackoff

import (
	"testing"

	"github.com/gaorx/stardust5/sderr"
	"github.com/stretchr/testify/assert"
)

func TestZero(t *testing.T) {
	var n int

	n = 0
	Retry(Zero(), func() error {
		n++
		return nil
	})
	assert.Equal(t, 1, n)

	n = 0
	once := false
	Retry(Zero(), func() error {
		n++
		if once {
			return nil
		}
		once = true
		return sderr.New("some error")
	})
	assert.Equal(t, 2, n)
}

func TestStop(t *testing.T) {
	var n int

	n = 0
	Retry(Stop(), func() error {
		n++
		return nil
	})
	assert.Equal(t, 1, n)

	n = 0
	Retry(Stop(), func() error {
		n++
		return sderr.New("some error")
	})
	assert.Equal(t, 1, n)
}
