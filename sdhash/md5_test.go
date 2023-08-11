package sdhash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMd5(t *testing.T) {
	assert.Equal(t, "5D41402ABC4B2A76B9719D911017C592", Md5([]byte("hello")).HexU())
}
