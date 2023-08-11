package sdstrings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitNonempty(t *testing.T) {
	assert.Empty(t, SplitNonempty("", ",", false))
	assert.Empty(t, SplitNonempty(",", ",", false))
	assert.Empty(t, SplitNonempty(",,", ",", false))
	assert.Equal(t, []string{"a", "b"}, SplitNonempty(",a,,b", ",", false))
	assert.Equal(t, []string{"a", "b"}, SplitNonempty(", a , , b", ",", true))
}

func TestSplit2s(t *testing.T) {
	var s1, s2 string

	s1, s2 = Split2s("a.b", ".")
	assert.Equal(t, "a", s1)
	assert.Equal(t, "b", s2)

	s1, s2 = Split2s("", ".")
	assert.Equal(t, "", s1)
	assert.Equal(t, "", s2)

	s1, s2 = Split2s("a", ".")
	assert.Equal(t, "a", s1)
	assert.Equal(t, "", s2)

	s1, s2 = Split2s("a.b.c", ".")
	assert.Equal(t, "a", s1)
	assert.Equal(t, "b.c", s2)
}

func TestSplit3s(t *testing.T) {
	var s1, s2, s3 string

	s1, s2, s3 = Split3s("a.b.c", ".")
	assert.Equal(t, "a", s1)
	assert.Equal(t, "b", s2)
	assert.Equal(t, "c", s3)

	s1, s2, s3 = Split3s("", ".")
	assert.Equal(t, "", s1)
	assert.Equal(t, "", s2)
	assert.Equal(t, "", s3)

	s1, s2, s3 = Split3s("a.b", ".")
	assert.Equal(t, "a", s1)
	assert.Equal(t, "b", s2)
	assert.Equal(t, "", s3)

	s1, s2, s3 = Split3s("a.b.c.d", ".")
	assert.Equal(t, "a", s1)
	assert.Equal(t, "b", s2)
	assert.Equal(t, "c.d", s3)
}
