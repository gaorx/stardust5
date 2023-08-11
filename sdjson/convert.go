package sdjson

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

type Converter interface {
	ToBool(v any, as bool) (bool, bool)
	ToString(v any, as bool) (string, bool)
	ToInt(v any, as bool) (int64, bool)
	ToUint(v any, as bool) (uint64, bool)
	ToFloat(v any, as bool) (float64, bool)
	ToObject(v any, as bool) (Object, bool)
	ToArray(v any, as bool) (Array, bool)
	ToAny(v any, as bool, ptr any) bool
}

var converter = &mergedConverter{
	converters: []Converter{defaultConverter{}},
}

func Register(c Converter) {
	converter.add(c)
}

// to

func ToBool(v any) (bool, bool) {
	return converter.ToBool(v, false)
}

func ToString(v any) (string, bool) {
	return converter.ToString(v, false)
}

func ToInt[T int | int8 | int16 | int32 | int64](v any) (T, bool) {
	r, ok := converter.ToInt(v, false)
	if !ok {
		return 0, false
	}
	return T(r), true
}

func ToUint[T uint | uint8 | uint16 | uint32 | uint64](v any) (T, bool) {
	r, ok := converter.ToUint(v, false)
	if !ok {
		return 0, false
	}
	return T(r), true
}

func ToFloat[T float32 | float64](v any) (T, bool) {
	r, ok := converter.ToFloat(v, false)
	if !ok {
		return 0.0, false
	}
	return T(r), true
}

func ToObject(v any) (Object, bool) {
	r, ok := converter.ToObject(v, false)
	if !ok {
		return nil, false
	}
	return r, true
}

func ToArray(v any) (Array, bool) {
	r, ok := converter.ToArray(v, false)
	if !ok {
		return nil, false
	}
	return r, true
}

func To[T any](v any) (T, bool) {
	var r T
	ok := converter.ToAny(v, false, &r)
	return r, ok
}

// as

func AsBool(v any) (bool, bool) {
	return converter.ToBool(v, true)
}

func AsString(v any) (string, bool) {
	return converter.ToString(v, true)
}

func AsInt[T int | int8 | int16 | int32 | int64](v any) (T, bool) {
	r, ok := converter.ToInt(v, true)
	if !ok {
		return 0, false
	}
	return T(r), true
}

func AsUint[T uint | uint8 | uint16 | uint32 | uint64](v any) (T, bool) {
	r, ok := converter.ToUint(v, true)
	if !ok {
		return 0, false
	}
	return T(r), true
}

func AsFloat[T float32 | float64](v any) (T, bool) {
	r, ok := converter.ToFloat(v, true)
	if !ok {
		return 0.0, false
	}
	return T(r), true
}

func AsObject(v any) (Object, bool) {
	r, ok := converter.ToObject(v, true)
	if !ok {
		return nil, false
	}
	return r, true
}

func AsArray(v any) (Array, bool) {
	r, ok := converter.ToArray(v, true)
	if !ok {
		return nil, false
	}
	return r, true
}

func As[T any](v any) (T, bool) {
	var r T
	ok := converter.ToAny(v, true, &r)
	return r, ok
}

// as with default

func AsBoolDef(v any, def bool) bool {
	r, ok := AsBool(v)
	if !ok {
		return def
	}
	return r
}

func AsStringDef(v any, def string) string {
	r, ok := AsString(v)
	if !ok {
		return def
	}
	return r
}

func AsIntDef[T int | int8 | int16 | int32 | int64](v any, def T) T {
	r, ok := AsInt[T](v)
	if !ok {
		return def
	}
	return r
}

func AsUintDef[T uint | uint8 | uint16 | uint32 | uint64](v any, def T) T {
	r, ok := AsUint[T](v)
	if !ok {
		return def
	}
	return r
}

func AsFloatDef[T float32 | float64](v any, def T) T {
	r, ok := AsFloat[T](v)
	if !ok {
		return def
	}
	return r
}

func AsObjectDef(v any, def Object) Object {
	r, ok := AsObject(v)
	if !ok {
		return def
	}
	return r
}

func AsArrayDef(v any, def Array) Array {
	r, ok := AsArray(v)
	if !ok {
		return def
	}
	return r
}

func AsDef[T any](v any, def T) T {
	r, ok := As[T](v)
	if !ok {
		return def
	}
	return r
}

// merged converter

type mergedConverter struct {
	converters []Converter
}

func (mc *mergedConverter) add(c Converter) {
	if c != nil {
		mc.converters = append(mc.converters, c)
	}
}

func (mc *mergedConverter) ToBool(v any, as bool) (bool, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToBool(v, as)
		if ok {
			return v, ok
		}
	}
	return false, false
}

func (mc *mergedConverter) ToString(v any, as bool) (string, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToString(v, as)
		if ok {
			return v, ok
		}
	}
	return "", false
}

func (mc *mergedConverter) ToInt(v any, as bool) (int64, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToInt(v, as)
		if ok {
			return v, ok
		}
	}
	return 0, false
}

func (mc *mergedConverter) ToUint(v any, as bool) (uint64, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToUint(v, as)
		if ok {
			return v, ok
		}
	}
	return 0, false
}

func (mc *mergedConverter) ToFloat(v any, as bool) (float64, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToFloat(v, as)
		if ok {
			return v, ok
		}
	}
	return 0, false
}

func (mc *mergedConverter) ToObject(v any, as bool) (Object, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToObject(v, as)
		if ok {
			return v, ok
		}
	}
	return nil, false
}

func (mc *mergedConverter) ToArray(v any, as bool) (Array, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToArray(v, as)
		if ok {
			return v, ok
		}
	}
	return nil, false
}

func (mc *mergedConverter) ToAny(v any, as bool, ptr any) bool {
	for _, c := range mc.converters {
		ok := c.ToAny(v, as, ptr)
		if ok {
			return true
		}
	}
	return false
}

// default converter

type defaultConverter struct{}

func (c defaultConverter) ToBool(v any, as bool) (bool, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return false, false
	}

	// bool
	if v1, ok := v.(bool); ok {
		return v1, true
	}
	if as {
		switch v1 := v.(type) {
		// string
		case string:
			b, err := strconv.ParseBool(v1)
			if err != nil {
				return false, false
			}
			return b, true

		// number
		case json.Number:
			if isFloat(v1) {
				if v2, err := v1.Float64(); err != nil {
					return false, false
				} else {
					return v2 != 0.0, true
				}
			} else {
				if v2, err := v1.Int64(); err != nil {
					return false, false
				} else {
					return v2 != 0, true
				}
			}
		case int:
			return v1 != 0, true
		case int64:
			return v1 != 0, true
		case uint:
			return v1 != 0, true
		case uint64:
			return v1 != 0, true
		case float64:
			return v1 != 0.0, true
		case int8:
			return v1 != 0, true
		case int16:
			return v1 != 0, true
		case int32:
			return v1 != 0, true
		case uint8:
			return v1 != 0, true
		case uint16:
			return v1 != 0, true
		case uint32:
			return v1 != 0, true
		case float32:
			return v1 != 0.0, true
		}
	}
	return false, false
}

func (c defaultConverter) ToString(v any, as bool) (string, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return "", false
	}
	// string
	if v1, ok := v.(string); ok {
		return v1, true
	}

	if as {
		switch v1 := v.(type) {
		// bool
		case bool:
			return strconv.FormatBool(v1), true
		// number
		case json.Number:
			return v1.String(), true
		case int:
			return strconv.FormatInt(int64(v1), 10), true
		case int64:
			return strconv.FormatInt(v1, 10), true
		case uint:
			return strconv.FormatUint(uint64(v1), 10), true
		case uint64:
			return strconv.FormatUint(v1, 10), true
		case float64:
			return strconv.FormatFloat(v1, 'f', -1, 64), true
		case int8:
			return strconv.FormatInt(int64(v1), 10), true
		case int16:
			return strconv.FormatInt(int64(v1), 10), true
		case int32:
			return strconv.FormatInt(int64(v1), 10), true
		case uint8:
			return strconv.FormatUint(uint64(v1), 10), true
		case uint16:
			return strconv.FormatUint(uint64(v1), 10), true
		case uint32:
			return strconv.FormatUint(uint64(v1), 10), true
		case float32:
			return strconv.FormatFloat(float64(v1), 'f', -1, 32), true
		}
	}

	return "", false
}

func (c defaultConverter) ToInt(v any, as bool) (int64, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return 0, false
	}

	// number
	switch v1 := v.(type) {
	case json.Number:
		if isFloat(v1) {
			if v2, err := v1.Float64(); err != nil {
				return 0, false
			} else {
				return int64(v2), true
			}
		} else {
			if v2, err := v1.Int64(); err != nil {
				return 0, false
			} else {
				return v2, true
			}
		}
	case int:
		return int64(v1), true
	case int64:
		return v1, true
	case uint:
		return int64(v1), true
	case uint64:
		return int64(v1), true
	case float64:
		return int64(v1), true
	case int8:
		return int64(v1), true
	case int16:
		return int64(v1), true
	case int32:
		return int64(v1), true
	case uint8:
		return int64(v1), true
	case uint16:
		return int64(v1), true
	case uint32:
		return int64(v1), true
	case float32:
		return int64(v1), true
	}

	if as {
		switch v1 := v.(type) {
		// bool
		case bool:
			if v1 {
				return 1, true
			} else {
				return 0, true
			}
		// string
		case string:
			if isFloat(json.Number(v1)) {
				if v2, err := json.Number(v1).Float64(); err != nil {
					return 0, false
				} else {
					return int64(v2), true
				}
			} else {
				if v2, err := json.Number(v1).Int64(); err != nil {
					return 0, false
				} else {
					return v2, true
				}
			}
		}
	}

	return 0, false
}

func (c defaultConverter) ToUint(v any, as bool) (uint64, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return 0, false
	}

	// number
	switch v1 := v.(type) {
	case json.Number:
		if isFloat(v1) {
			if v2, err := v1.Float64(); err != nil {
				return 0, false
			} else {
				return uint64(v2), true
			}
		} else {
			if v2, err := strconv.ParseUint(string(v1), 10, 64); err != nil {
				return 0, false
			} else {
				return v2, true
			}
		}
	case int:
		return uint64(v1), true
	case int64:
		return uint64(v1), true
	case uint:
		return uint64(v1), true
	case uint64:
		return v1, true
	case float64:
		return uint64(v1), true
	case int8:
		return uint64(v1), true
	case int16:
		return uint64(v1), true
	case int32:
		return uint64(v1), true
	case uint8:
		return uint64(v1), true
	case uint16:
		return uint64(v1), true
	case uint32:
		return uint64(v1), true
	case float32:
		return uint64(v1), true
	}

	if as {
		switch v1 := v.(type) {
		// bool
		case bool:
			if v1 {
				return 1, true
			} else {
				return 0, true
			}
		// string
		case string:
			if isFloat(json.Number(v1)) {
				if v2, err := json.Number(v1).Float64(); err != nil {
					return 0, false
				} else {
					return uint64(v2), true
				}
			} else {
				if v2, err := strconv.ParseUint(v1, 10, 64); err != nil {
					return 0, false
				} else {
					return v2, true
				}
			}
		}
	}

	return 0, false
}

func (c defaultConverter) ToFloat(v any, as bool) (float64, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return 0, false
	}

	// number
	switch v1 := v.(type) {
	case json.Number:
		if v2, err := v1.Float64(); err != nil {
			return 0, false
		} else {
			return v2, true
		}
	case float64:
		return v1, true
	case float32:
		return float64(v1), true
	case int:
		return float64(v1), true
	case int64:
		return float64(v1), true
	case uint:
		return float64(v1), true
	case uint64:
		return float64(v1), true
	case int8:
		return float64(v1), true
	case int16:
		return float64(v1), true
	case int32:
		return float64(v1), true
	case uint8:
		return float64(v1), true
	case uint16:
		return float64(v1), true
	case uint32:
		return float64(v1), true
	}

	if as {
		switch v1 := v.(type) {
		// bool
		case bool:
			if v1 {
				return 1.0, true
			} else {
				return 0.0, true
			}
		// string
		case string:
			if v2, err := json.Number(v1).Float64(); err != nil {
				return 0.0, false
			} else {
				return v2, true
			}
		}
	}

	return 0.0, false
}

func (c defaultConverter) ToObject(v any, as bool) (Object, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return nil, false
	}

	// map like
	if v1, ok := v.(map[string]any); ok {
		return v1, true
	} else if v1, ok := v.(Object); ok {
		return v1, true
	} else if rv := reflect.ValueOf(v); rv.Type().Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String {
		return genericMapToObject(rv)
	}

	if as {
		// struct
		rv := reflect.ValueOf(v)
		rt := rv.Type()
		if rt.Kind() == reflect.Struct {
			return structToObject(rv)
		} else if rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Struct {
			return structToObject(rv)
		}
	}

	return nil, false
}

func (c defaultConverter) ToArray(v any, as bool) (Array, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return nil, false
	}

	// slice like
	if v1, ok := v.([]any); ok {
		return v1, true
	} else if v1, ok := v.(Array); ok {
		return v1, true
	} else if rv := reflect.ValueOf(v); rv.Type().Kind() == reflect.Slice || rv.Type().Kind() == reflect.Array {
		return genericSliceToArray(rv)
	}

	return nil, false
}

func (c defaultConverter) ToAny(v any, as bool, ptr any) bool {
	raw, err := json.Marshal(v)
	if err != nil {
		return false
	}
	err = json.Unmarshal(raw, ptr)
	if err != nil {
		return false
	}
	return true
}

func genericMapToObject(rv reflect.Value) (map[string]any, bool) {
	if rv.IsNil() {
		return nil, true
	}
	l := rv.Len()
	m := make(map[string]any, l)
	iter := rv.MapRange()
	for iter.Next() {
		k := iter.Key().Interface().(string)
		v := iter.Value().Interface()
		m[k] = v
	}
	return m, true
}

func structToObject(rv reflect.Value) (map[string]any, bool) {
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return nil, true
	}
	raw, err := json.Marshal(rv.Interface())
	if err != nil {
		return nil, false
	}
	var m map[string]any
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, false
	}
	return m, true
}

func genericSliceToArray(rv reflect.Value) ([]any, bool) {
	if rv.Kind() == reflect.Slice && rv.IsNil() {
		return nil, true
	}
	l := rv.Len()
	a := make([]any, 0, l)
	for i := 0; i < l; i++ {
		a = append(a, rv.Index(i).Interface())
	}
	return a, true
}

func isFloat(n json.Number) bool {
	return strings.Contains(n.String(), ".")
}
