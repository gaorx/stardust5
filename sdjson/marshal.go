package sdjson

import (
	"encoding/json"
	"github.com/samber/lo"
)

// bytes

var (
	Unmarshal     = json.Unmarshal
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
)

// string

func UnmarshalString(s string, v any) error {
	return json.Unmarshal([]byte(s), v)
}

func MarshalString(v any) (string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func MarshalIndentString(v any, prefix, indent string) (string, error) {
	raw, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func MarshalStringDef(v any, def string) string {
	if r, err := MarshalString(v); err != nil {
		return def
	} else {
		return r
	}
}

func MarshalIndentStringDef(v any, prefix, indent, def string) string {
	if r, err := MarshalIndentString(v, prefix, indent); err != nil {
		return def
	} else {
		return r
	}
}

func MarshalPretty(v any) string {
	return MarshalIndentStringDef(v, "", "  ", "")
}

// typed

func UnmarshalTyped[T any](raw []byte) (T, error) {
	var t T
	if err := json.Unmarshal(raw, &t); err != nil {
		return lo.Empty[T](), err
	}
	return t, nil
}

func UnmarshalStringTyped[T any](j string) (T, error) {
	return UnmarshalTyped[T]([]byte(j))
}

func UnmarshalTypedDef[T any](raw []byte, def T) T {
	if t, err := UnmarshalTyped[T](raw); err != nil {
		return def
	} else {
		return t
	}
}

func UnmarshalStringTypedDef[T any](j string, def T) T {
	return UnmarshalTypedDef[T]([]byte(j), def)
}

// value

func UnmarshalValue(raw []byte) (Value, error) {
	return UnmarshalTyped[Value](raw)
}

func UnmarshalValueString(s string) (Value, error) {
	if v, err := UnmarshalValue([]byte(s)); err != nil {
		return V(nil), err
	} else {
		return v, nil
	}
}
