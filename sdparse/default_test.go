package sdparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	assert.Equal(t, IntDef("33", 44), 33)
	assert.Equal(t, IntDef("error33", 44), 44)
}
