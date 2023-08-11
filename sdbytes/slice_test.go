package sdbytes

import (
	"github.com/gaorx/stardust5/sdencoding"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlice(t *testing.T) {
	a := []byte("hello")
	assert.Equal(t, a, lo.Must(sdencoding.Hex.DecodeString(Slice(a).HexL())))
	assert.Equal(t, a, lo.Must(sdencoding.Hex.DecodeString(Slice(a).HexU())))
	assert.Equal(t, a, lo.Must(sdencoding.Base64Std.DecodeString(Slice(a).Base64Std())))
	assert.Equal(t, a, lo.Must(sdencoding.Base64Url.DecodeString(Slice(a).Base64Url())))
}
