package sdjwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEncodeAndDecode(t *testing.T) {
	type user struct {
		UID        string `json:"uid"`
		Expiration int64  `json:"expiration"`
	}

	const secret = "QphlY11dKQ24IoZr"
	u0 := user{
		UID:        "3939939",
		Expiration: time.Now().UnixMilli(),
	}
	token, err := Encode(secret, u0)
	assert.NoError(t, err)
	u1, err := Decode[user](secret, token)
	assert.NoError(t, err)
	assert.Equal(t, u0, u1)
}
