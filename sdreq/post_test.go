package sdreq

import (
	"testing"

	"github.com/gaorx/stardust5/sdjson"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	statusCode, body, err := PostForJson[sdjson.Object](
		nil,
		"https://httpbin.org/post?k1=v1",
		sdjson.Object{"req_k1": "req_v1"},
		QueryParam("k2", 2),
	)
	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.Equal(t, "v1", body.Get("args").Get("k1").AsStringDef(""))
	assert.Equal(t, "2", body.Get("args").Get("k2").AsStringDef(""))
	assert.Equal(t, "req_v1", body.Get("json").Get("req_k1").AsStringDef(""))
}
