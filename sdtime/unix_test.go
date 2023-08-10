package sdtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnix(t *testing.T) {
	nowS := NowUnixS()
	assert.Equal(t, nowS, ToUnixS(FromUnixS(nowS)))

	nowMs := NowUnixMs()
	assert.Equal(t, nowMs, ToUnixMs(FromUnixMs(nowMs)))
}
