package sdhttpua

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFind(t *testing.T) {
	ual := Find(PlatformIs("Windows", "Linux"), IsMobile().Not())
	assert.True(t, len(ual) > 0)
	for _, ua := range ual {
		assert.True(t, ua.Platform == "Windows" || ua.Platform == "Linux")
		assert.False(t, ua.Mobile)
	}
}
