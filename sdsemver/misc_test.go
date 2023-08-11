package sdsemver

import (
	"testing"

	"github.com/gaorx/stardust5/sdrand"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	for i := 0; i < 10000; i++ {
		major := sdrand.IntBetween(0, numLimit)
		minor := sdrand.IntBetween(0, numLimit)
		patch := sdrand.IntBetween(0, numLimit)
		s0 := New(major, minor, patch).String()
		vi, err := ToInt(s0)
		assert.NoError(t, err)
		s1, err := ToString(vi)
		assert.NoError(t, err)
		assert.True(t, s0 == s1)
	}
}
