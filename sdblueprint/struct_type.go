package sdblueprint

import (
	"reflect"
)

type structType struct {
	reflect.Type
}

type markedField struct {
	name string
	mark reflect.Type
	tag  structTag
}

func structTypeOf(v any) (structType, bool) {
	switch t := v.(type) {
	case nil:
		return structType{}, false
	case reflect.Type:
		if t.Kind() == reflect.Struct {
			return structType{t}, true
		} else if t.Kind() == reflect.Pointer {
			return structTypeOf(t.Elem())
		} else {
			return structType{t}, false
		}
	case reflect.Value:
		return structTypeOf(t.Type())
	default:
		return structTypeOf(reflect.TypeOf(t))
	}
}

func (st structType) isZero() bool {
	return st.Type == nil
}

func (st structType) forEachField(action func(reflect.StructField)) {
	n := st.NumField()
	for i := 0; i < n; i++ {
		sf := st.Field(i)
		action(sf)
	}
}

func (st structType) findStructMarkIn(markSet markSet) []markedField {
	var r []markedField
	st.forEachField(func(sf reflect.StructField) {
		mark, ok := getFieldMark(sf.Type, markSet)
		if ok {
			r = append(r, markedField{name: sf.Name, mark: mark, tag: structTag(sf.Tag)})
		}
	})
	return r
}

func newFieldMark(sf reflect.StructField) markedField {
	return markedField{name: sf.Name, mark: sf.Type, tag: structTag(sf.Tag)}
}

func (mark markedField) getId(st *structType) string {
	if st != nil {
		return selectNotEmpty[string](mark.tag.Get("id"), mark.name, st.Name())
	} else {
		return selectNotEmpty[string](mark.tag.Get("id"), mark.name, "")
	}
}
