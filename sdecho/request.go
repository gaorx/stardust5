package sdecho

import (
	"encoding/json"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"io"
)

func (c Context) RequestBodyBytes() ([]byte, error) {
	reader := c.Request().Body
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, sderr.Wrap(err, "sdecho read request body error")
	}
	return r, nil
}

func (c Context) RequestBodyString() (string, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c Context) RequestBodyAsJsonValue() (sdjson.Value, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return sdjson.Value{}, err
	}
	v, err := sdjson.UnmarshalValue(b)
	if err != nil {
		return sdjson.Value{}, err
	}
	return v, nil
}

func (c Context) RequestBodyAsJsonObject() (sdjson.Object, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c Context) RequestBodyAsJsonArray() (sdjson.Array, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return nil, err
	}
	var a []any
	err = json.Unmarshal(b, &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (c Context) RequestBodyAs(v any) error {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
