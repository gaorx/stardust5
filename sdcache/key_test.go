package sdcache

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestInt64Key(t *testing.T) {
	for _, prefix := range []string{"", "P/"} {
		key := Int64Key{Prefix: prefix}
		k0 := int64(333)
		// encode
		_, err := key.EncodeKey(nil)
		assert.Error(t, err)
		_, err = key.EncodeKey(strconv.FormatInt(k0, 10))
		assert.Error(t, err)
		s, err := key.EncodeKey(k0)
		assert.NoError(t, err)
		assert.Equal(t, prefix+strconv.FormatInt(k0, 10), s)

		// decode
		_, err = key.DecodeKey("")
		assert.Error(t, err)
		k1, err := key.DecodeKey(s)
		assert.NoError(t, err)
		assert.Equal(t, k0, k1)
	}
}

func TestStrKey(t *testing.T) {
	for _, prefix := range []string{"", "P/"} {
		key := StrKey{Prefix: prefix}
		k0 := "k333"
		// encode
		_, err := key.EncodeKey(nil)
		assert.Error(t, err)
		_, err = key.EncodeKey(333)
		assert.Error(t, err)
		s, err := key.EncodeKey(k0)
		assert.NoError(t, err)
		assert.Equal(t, prefix+k0, s)

		// decode
		_, err = key.DecodeKey("")
		if prefix == "" {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
		k1, err := key.DecodeKey(s)
		assert.NoError(t, err)
		assert.Equal(t, k0, k1)
	}
}
