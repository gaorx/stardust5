package sdlocal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIP(t *testing.T) {
	// IP4
	ips, err := IPs(Is4())
	assert.NoError(t, err)
	for _, ip := range ips {
		assert.True(t, len(ip.To4()) > 0)
	}

	// !IP4
	ips, err = IPs(Is4().Not())
	assert.NoError(t, err)
	for _, ip := range ips {
		assert.True(t, len(ip.To4()) <= 0)
	}

	// Loopback
	ips, err = IPs(IsLoopback())
	assert.NoError(t, err)
	for _, ip := range ips {
		assert.True(t, ip.IsLoopback())
	}

	// Private
	ips, err = IPs(IsPrivate())
	assert.NoError(t, err)
	for _, ip := range ips {
		assert.True(t, ip.IsPrivate())
	}
}
