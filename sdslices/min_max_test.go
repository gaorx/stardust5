package sdslices

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinMax(t *testing.T) {
	assert.Equal(t, 2, MinV(4, 3, 5, 2))
	assert.Equal(t, 5, MaxV(4, 3, 5, 2))
}
