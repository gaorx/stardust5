package sdurl

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModify(t *testing.T) {
	url0 := "https://host1.com:3333/seg1/seg2?k2=v2&k1=v1"
	url1, err := Modify(url0)
	assert.NoError(t, err)
	assert.True(t, url0 == url1)

	// hostname
	url2, err := Modify(url0, SetHostname("host2.com"))
	assert.NoError(t, err)
	url2a, err := url.Parse(url2)
	assert.NoError(t, err)
	assert.Equal(t, "host2.com", url2a.Hostname())
	assert.Equal(t, "3333", url2a.Port())
	url2b, err := Modify(url2, SetHostname("host1.com"))
	assert.NoError(t, err)
	assert.True(t, url2b == url0)

	// port
	url3, err := Modify(url0, SetPort("4444"))
	assert.NoError(t, err)
	url3a, err := url.Parse(url3)
	assert.NoError(t, err)
	assert.Equal(t, "host1.com", url3a.Hostname())
	assert.Equal(t, "4444", url3a.Port())
	url3b, err := Modify(url3, SetPort("3333"))
	assert.NoError(t, err)
	assert.True(t, url3b == url0)

	// path
	url4, err := Modify(url0, SetPath("/abc"))
	assert.NoError(t, err)
	url4a, err := url.Parse(url4)
	assert.NoError(t, err)
	assert.Equal(t, "/abc", url4a.Path)
	url4b, err := Modify(url4, SetPath("/seg1/seg2"))
	assert.NoError(t, err)
	assert.True(t, url4b == url0)

	// TODO: params
}
