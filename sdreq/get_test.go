package sdreq

import (
	"context"
	"testing"

	"github.com/gaorx/stardust5/sdjson"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	statusCode, body, err := GetForJson[sdjson.Object](
		context.Background(),
		nil,
		"https://httpbin.org/get?k1=v1",
		QueryParam("k2", 2),
	)
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "v1", body.Get("args").Get("k1").AsStringDef(""))
	assert.Equal(t, "2", body.Get("args").Get("k2").AsStringDef(""))
}
