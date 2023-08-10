package sdtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnix(t *testing.T) {
	nowS := NowUnixS()
	assert.Equal(t, nowS, ToUnixS(FromUnixS(nowS)))

	nowMs := NowUnixMillis()
	assert.Equal(t, nowMs, ToUnixMillis(FromUnixMillis(nowMs)))
}
