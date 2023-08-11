package sdcrypto

import (
	"bytes"
	"testing"

	"github.com/gaorx/stardust5/sdrand"
	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {
	data0 := []byte(sdrand.String(1303, sdrand.LowerCaseLettersCharset))
	key := []byte(sdrand.String(16, sdrand.AlphanumericCharset))

	data1, err := AES.Encrypt(key, data0)
	assert.NoError(t, err)
	data2, err := AES.Decrypt(key, data1)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(data0, data2))

	data3, err := AESCRC32.Encrypt(key, data0)
	assert.NoError(t, err)
	data4, err := AESCRC32.Decrypt(key, data3)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(data0, data4))
}
