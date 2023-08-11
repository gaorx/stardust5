package sdhash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSha1(t *testing.T) {
	assert.Equal(t,
		"AAF4C61DDCC5E8A2DABEDE0F3B482CD9AEA9434D",
		Sha1([]byte("hello")).HexU(),
	)
}

func TestSha256(t *testing.T) {
	assert.Equal(t,
		"2CF24DBA5FB0A30E26E83B2AC5B9E29E1B161E5C1FA7425E73043362938B9824",
		Sha256([]byte("hello")).HexU(),
	)
}
