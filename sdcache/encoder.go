package sdcache

import (
	"encoding/json"
	"github.com/gaorx/stardust5/sderr"
)

type Encoder interface {
	EncodeValue(k, v any) ([]byte, error)
	DecodeValue(data []byte) (any, error)
}

type TextEncoder struct{}
type JsonEncoder[T any] struct{}

func (t TextEncoder) EncodeValue(k, v any) ([]byte, error) {
	if v == nil {
		return nil, sderr.New("nil value encode to string")
	}
	v1, ok := v.(string)
	if !ok {
		return nil, sderr.New("illegal value type")
	}
	return []byte(v1), nil
}

func (t TextEncoder) DecodeValue(data []byte) (any, error) {
	return string(data), nil
}

func (j JsonEncoder[T]) EncodeValue(k, v any) ([]byte, error) {
	if v == nil {
		return nil, sderr.New("nil value encode to json")
	}
	v1, ok := v.(T)
	if !ok {
		return nil, sderr.New("encode json type error")
	}
	return json.Marshal(v1)
}

func (j JsonEncoder[T]) DecodeValue(data []byte) (any, error) {
	var v T
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
