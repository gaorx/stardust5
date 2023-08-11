package sdencoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestText(t *testing.T) {
	s := "你好，世界!"
	testTextHelper(t, s, GB2312)
	testTextHelper(t, s, GBK)
	testTextHelper(t, s, GB18030)
}

func testTextHelper(t *testing.T, s string, e TextEncoding) {
	encoded, err := e.Encode(s)
	assert.NoError(t, err)
	r, err := e.Decode(encoded)
	assert.NoError(t, err)
	assert.Equal(t, s, r)
}
