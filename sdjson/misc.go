package sdjson

import (
	"encoding/json"
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
)

func StructToObject(v any) (Object, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return nil, sderr.Wrap(err, "marshal json error")
	}
	var v1 map[string]any
	err = json.Unmarshal(j, &v1)
	if err != nil {
		return nil, sderr.Wrap(err, "unmarshal to object error")
	}
	return v1, nil
}

func ObjectToStruct[T any](o Object) (T, error) {
	j, err := json.Marshal(o)
	if err != nil {
		return lo.Empty[T](), sderr.Wrap(err, "marshal json object error")
	}
	var v1 T
	err = json.Unmarshal(j, &v1)
	if err != nil {
		return lo.Empty[T](), sderr.Wrap(err, "unmarshal to struct error")
	}
	return v1, nil
}

func ToPrimitive(v any) (any, error) {
	if v == nil {
		return nil, nil
	}
	switch v1 := v.(type) {
	case string, bool:
		return v1, nil
	case int, int8, int16, int32, int64:
		return v1, nil
	case uint, uint8, uint16, uint32, uint64:
		return v1, nil
	case float32, float64:
		return v1, nil
	case json.Number:
		return v1, nil
	default:
		return MarshalString(v)
	}
}

func ToPrimitiveDef(v any, def any) any {
	if v1, err := ToPrimitive(v); err != nil {
		return def
	} else {
		return v1
	}
}

func TryPrimitive(v any) any {
	if v1, err := ToPrimitive(v); err != nil {
		return v
	} else {
		return v1
	}
}
