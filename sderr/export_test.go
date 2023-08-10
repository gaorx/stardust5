package sderr

import (
	"fmt"
	"strings"
	"testing"

	stderrors "errors"

	"github.com/stretchr/testify/assert"
)

func TestIsAndAs(t *testing.T) {
	root1 := myErr1{1}
	root2 := &myErr2{2}
	err1 := Wrap(root1, "wrap1")
	err2 := Wrap(root2, "wrap2")

	// is
	assert.True(t, Is(err1, root1))
	assert.False(t, Is(err1, root2))
	assert.True(t, Is(err2, root2))
	assert.False(t, Is(err2, root1))

	// try (myErr1)
	root1a, ok := TryT[myErr1](err1)
	assert.True(t, ok)
	assert.Equal(t, root1, root1a)
	root1b, ok := TryT[*myErr2](err1)
	assert.False(t, ok)
	assert.Nil(t, root1b)
	// as (*myErr2)
	root2a, ok := TryT[*myErr2](err2)
	assert.True(t, ok)
	assert.Equal(t, root2, root2a)
	root2b, ok := TryT[myErr1](err2)
	assert.False(t, ok)
	assert.Equal(t, myErr1{}, root2b)

	// std errors 兼容性
	root3 := stderrors.New("xxx")
	wrap3 := Wrap(root3, "yyy")
	assert.True(t, stderrors.Is(wrap3, root3))
	assert.True(t, Cause(wrap3) == root3)
	assert.True(t, Unwrap(wrap3) == root3)
	assert.True(t, stderrors.Unwrap(wrap3) == root3)

	// 测试消息1
	root4 := New("xxx")
	wrap4a := Wrap(root4, "yyy")
	wrap4b := Wrap(wrap4a, "zzz")
	assert.True(t, strings.Contains(wrap4b.Error(), "xxx"))
	assert.True(t, strings.Contains(wrap4b.Error(), "yyy"))
	assert.True(t, strings.Contains(wrap4b.Error(), "zzz"))

	// 测试消息1
	root5 := stderrors.New("aaa")
	wrap5a := Wrap(root5, "bbb")
	wrap5b := Wrap(wrap5a, "ccc")
	assert.True(t, strings.Contains(wrap5b.Error(), "aaa"))
	assert.True(t, strings.Contains(wrap5b.Error(), "bbb"))
	assert.True(t, strings.Contains(wrap5b.Error(), "ccc"))
}

type myErr1 struct {
	Code int
}

func (e myErr1) Error() string {
	return fmt.Sprintf("myErr1(%d)", e.Code)
}

type myErr2 struct {
	Code int
}

func (e *myErr2) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("myErr2(%d)", e.Code)
}
