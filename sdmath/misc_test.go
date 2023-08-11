package sdmath

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalize(t *testing.T) {
	assert.Equal(t, 4.0, Normalize(2.0, Interval{1.0, 3.0}, Interval{2.0, 6.0}))
	assert.Panics(t, func() {
		Normalize(2.0, Interval{3.0, 3.0}, Interval{2.0, 6.0})
	})
}
