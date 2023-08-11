package sdreflect

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/oleiade/reflections"
)

var (
	StructFields        = reflections.Fields
	StructHasField      = reflections.HasField
	StructFieldKind     = reflections.GetFieldKind
	StructFieldType     = reflections.GetFieldType
	StructFieldTag      = reflections.GetFieldTag
	StructGetFieldValue = reflections.GetField
	StructSetFieldValue = reflections.SetField
)

func StructToMap(obj any) map[string]any {
	m, err := reflections.Items(obj)
	if err != nil {
		return nil
	}
	return m
}

func StructSelectFields(dest, src any, fields []string) error {
	var fields1 []string
	for _, f := range fields {
		ok, err := StructHasField(src, f)
		if err != nil {
			return sderr.WrapWith(err, "has field error", f)
		}
		if ok {
			fields1 = append(fields1, f)
		}
	}
	for _, f := range fields1 {
		v, err := StructGetFieldValue(src, f)
		if err != nil {
			return sderr.WrapWith(err, "read field value error", f)
		}
		err = StructSetFieldValue(dest, f, v)
		if err != nil {
			return sderr.WrapWith(err, "set field value error", f)
		}
	}
	return nil
}
