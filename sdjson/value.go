package sdjson

import (
	"encoding/json"
	"github.com/gaorx/stardust5/sderr"
)

type Value struct {
	v any
}

func V(v any) Value {
	return Value{unbox(v)}
}

func (v Value) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(v.v)
	if err != nil {
		return nil, sderr.Wrap(err, "sdjson marshal value to json")
	}
	return raw, nil
}

func (v *Value) UnmarshalJSON(raw []byte) error {
	var v0 any
	err := json.Unmarshal(raw, &v0)
	if err != nil {
		return sderr.Wrap(err, "sdjson unmarshal json to value")
	}
	v.v = v0
	return nil
}

func (v Value) Interface() any {
	return v.v
}

func (v Value) IsNil() bool {
	return v.v == nil
}

// Get / At

func (v Value) Get(k string) Value {
	return v.AsObjectDef(nil).Get(k)
}

func (v Value) Gets(k string, subKeys ...string) Value {
	if len(subKeys) <= 0 {
		return v.Get(k)
	} else {
		r := v.Get(k)
		for _, subKey := range subKeys {
			r = r.Get(subKey)
		}
		return r
	}
}

func (v Value) At(i int) Value {
	return v.AsArrayDef(nil).At(i)
}

// convert: ToXxx

func (v Value) ToBool() (bool, bool) {
	return ToBool(v)
}

func (v Value) ToString() (string, bool) {
	return ToString(v)
}

func (v Value) ToInt() (int, bool) {
	return ToInt[int](v)
}

func (v Value) ToInt8() (int8, bool) {
	return ToInt[int8](v)
}

func (v Value) ToInt16() (int16, bool) {
	return ToInt[int16](v)
}

func (v Value) ToInt32() (int32, bool) {
	return ToInt[int32](v)
}

func (v Value) ToInt64() (int64, bool) {
	return ToInt[int64](v)
}

func (v Value) ToUint() (uint, bool) {
	return ToUint[uint](v)
}

func (v Value) ToUint8() (uint8, bool) {
	return ToUint[uint8](v)
}

func (v Value) ToUint16() (uint16, bool) {
	return ToUint[uint16](v)
}

func (v Value) ToUint32() (uint32, bool) {
	return ToUint[uint32](v)
}

func (v Value) ToUint64() (uint64, bool) {
	return ToUint[uint64](v)
}

func (v Value) ToFloat32() (float32, bool) {
	return ToFloat[float32](v)
}

func (v Value) ToFloat64() (float64, bool) {
	return ToFloat[float64](v)
}

func (v Value) ToObject() (Object, bool) {
	return ToObject(v)
}

func (v Value) ToArray() (Array, bool) {
	return ToArray(v)
}

func (v Value) To(ptr any) bool {
	return converter.ToAny(v, false, ptr)
}

// convert: AsXxx

func (v Value) AsBool() (bool, bool) {
	return AsBool(v)
}

func (v Value) AsString() (string, bool) {
	return AsString(v)
}

func (v Value) AsInt() (int, bool) {
	return AsInt[int](v)
}

func (v Value) AsInt8() (int8, bool) {
	return AsInt[int8](v)
}

func (v Value) AsInt16() (int16, bool) {
	return AsInt[int16](v)
}

func (v Value) AsInt32() (int32, bool) {
	return AsInt[int32](v)
}

func (v Value) AsInt64() (int64, bool) {
	return AsInt[int64](v)
}

func (v Value) AsUint() (uint, bool) {
	return AsUint[uint](v)
}

func (v Value) AsUint8() (uint8, bool) {
	return AsUint[uint8](v)
}

func (v Value) AsUint16() (uint16, bool) {
	return AsUint[uint16](v)
}

func (v Value) AsUint32() (uint32, bool) {
	return AsUint[uint32](v)
}

func (v Value) AsUint64() (uint64, bool) {
	return AsUint[uint64](v)
}

func (v Value) AsFloat32() (float32, bool) {
	return AsFloat[float32](v)
}

func (v Value) AsFloat64() (float64, bool) {
	return AsFloat[float64](v)
}

func (v Value) AsObject() (Object, bool) {
	return AsObject(v)
}

func (v Value) AsArray() (Array, bool) {
	return AsArray(v)
}

func (v Value) As(ptr any) bool {
	return converter.ToAny(v, true, ptr)
}

// converter: AsXxxDef

func (v Value) AsBoolDef(def bool) bool {
	return AsBoolDef(v, def)
}

func (v Value) AsStringDef(def string) string {
	return AsStringDef(v, def)
}

func (v Value) AsIntDef(def int) int {
	return AsIntDef[int](v, def)
}

func (v Value) AsInt8Def(def int8) int8 {
	return AsIntDef[int8](v, def)
}

func (v Value) AsInt16Def(def int16) int16 {
	return AsIntDef[int16](v, def)
}

func (v Value) AsInt32Def(def int32) int32 {
	return AsIntDef[int32](v, def)
}

func (v Value) AsInt64Def(def int64) int64 {
	return AsIntDef[int64](v, def)
}

func (v Value) AsUintDef(def uint) uint {
	return AsUintDef[uint](v, def)
}

func (v Value) AsUint8Def(def uint8) uint8 {
	return AsUintDef[uint8](v, def)
}

func (v Value) AsUint16Def(def uint16) uint16 {
	return AsUintDef[uint16](v, def)
}

func (v Value) AsUint32Def(def uint32) uint32 {
	return AsUintDef[uint32](v, def)
}

func (v Value) AsUint64Def(def uint64) uint64 {
	return AsUintDef[uint64](v, def)
}

func (v Value) AsFloat32Def(def float32) float32 {
	return AsFloatDef[float32](v, def)
}

func (v Value) AsFloat64Def(def float64) float64 {
	return AsFloatDef[float64](v, def)
}

func (v Value) AsObjectDef(def Object) Object {
	return AsObjectDef(v, def)
}

func (v Value) AsArrayDef(def Array) Array {
	return AsArrayDef(v, def)
}

// helpers

func unbox(v any) any {
	switch v1 := v.(type) {
	case nil:
		return nil
	case Value:
		return v1.v
	case *Value:
		if v1 == nil {
			return nil
		} else {
			return v1.v
		}
	default:
		return v
	}
}
