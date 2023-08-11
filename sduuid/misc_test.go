package sduuid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	uuids := make(map[string]int)
	for i := 0; i < 100; i++ {
		id1 := NewV1().HexL()
		assert.Equal(t, 32, len(id1))
		_, ok := uuids[id1]
		assert.False(t, ok)
		uuids[id1] = 0
	}

	for i := 0; i < 100; i++ {
		id4 := NewV4().HexL()
		assert.Equal(t, 32, len(id4))
		_, ok := uuids[id4]
		assert.False(t, ok)
		uuids[id4] = 0
	}
}
