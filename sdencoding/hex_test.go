package sdencoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHex(t *testing.T) {
	b := []byte("hello 你好")
	testEncodingHelper(t, b, Hex)
}

func testEncodingHelper(t *testing.T, b []byte, e Encoding) {
	r1, err := e.Decode(e.Encode(b))
	assert.NoError(t, err)
	assert.Equal(t, b, r1)

	r2, err := e.DecodeString(e.EncodeString(b))
	assert.NoError(t, err)
	assert.Equal(t, b, r2)
}
