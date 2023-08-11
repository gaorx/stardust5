package sdurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinPath(t *testing.T) {
	assert.Equal(t, "/a/b/c/d", JoinPath("a/b", "c/d/", ""))
	assert.Equal(t, "/a/%20b/c", JoinPath("a/ b", "c"))
}
