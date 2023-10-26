package sdjson

import (
	"encoding/json"
	"github.com/gaorx/stardust5/sdencoding"
	"github.com/samber/lo"
)

func MarshalBase64(v any) (string, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return sdencoding.Base64Url.EncodeString(j), nil
}

func MarshalBase64Def(v any, def string) string {
	if b64, err := MarshalBase64(v); err != nil {
		return def
	} else {
		return b64
	}
}

func UnmarshalBase64Typed[T any](b64 string) (T, error) {
	j, err := sdencoding.Base64Url.DecodeString(b64)
	if err != nil {
		return lo.Empty[T](), err
	}
	return UnmarshalTyped[T](j)
}

func UnmarshalBase64TypedDef[T any](b64 string, def T) T {
	if v, err := UnmarshalBase64Typed[T](b64); err != nil {
		return def
	} else {
		return v
	}
}

func UnmarshalBase64(b64 string, v any) error {
	j, err := sdencoding.Base64Url.DecodeString(b64)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, v)
}
