package sdjson

import (
	"encoding/json"
	"github.com/samber/lo"

	"github.com/gaorx/stardust5/sderr"
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
