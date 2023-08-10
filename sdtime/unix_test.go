package sdtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnix(t *testing.T) {
	nowS := NowUnixS()
	assert.Equal(t, nowS, ToUnixS(FromUnixS(nowS)))

	nowMs := NowUnixMS()
	assert.Equal(t, nowMs, ToUnixMS(FromUnixMS(nowMs)))
}
