package sdcache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextEncoder(t *testing.T) {
	encoder := TextEncoder{}
	_, err := encoder.EncodeValue(nil, nil)
	assert.Error(t, err)
	_, err = encoder.EncodeValue(nil, 3)
	assert.Error(t, err)
	data, err := encoder.EncodeValue(nil, "abc")
	assert.NoError(t, err)
	assert.Equal(t, []byte("abc"), data)
	v, err := encoder.DecodeValue(data)
	assert.NoError(t, err)
	assert.Equal(t, "abc", v)
}

func TestJsonEncoder(t *testing.T) {
	type user struct {
		Id   int64  `json:"_id"`
		Name string `json:"_name"`
	}

	// struct
	encoder1 := JsonEncoder[user]{}
	u0 := user{Id: 333, Name: "user333"}
	_, err := encoder1.EncodeValue("k", nil)
	assert.Error(t, err)
	_, err = encoder1.EncodeValue("k", &u0)
	assert.Error(t, err)
	data, err := encoder1.EncodeValue("k", u0)
	assert.NoError(t, err)
	_, err = encoder1.DecodeValue(nil)
	assert.Error(t, err)
	u1, err := encoder1.DecodeValue(data)
	assert.NoError(t, err)
	assert.IsType(t, user{}, u1)
	assert.True(t, u0 == u1.(user))

	// struct pointer
	encoder2 := JsonEncoder[*user]{}
	_, err = encoder2.EncodeValue("k", nil)
	assert.Error(t, err)
	_, err = encoder2.EncodeValue("k", u0)
	assert.Error(t, err)
	data, err = encoder2.EncodeValue("k", &u0)
	assert.NoError(t, err)
	_, err = encoder2.DecodeValue(nil)
	assert.Error(t, err)
	u1, err = encoder2.DecodeValue(data)
	assert.NoError(t, err)
	assert.IsType(t, &user{}, u1)
	assert.True(t, u0 == *(u1.(*user)))
}
