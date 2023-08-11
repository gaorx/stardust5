package sdblueprint

import (
	"github.com/gaorx/stardust5/sdjson"
	"reflect"
)

type Field interface {
	Id() string
	Comment() string
	Jsonable
	Type() reflect.Type
	Attributes
}

var _ Field = &field{}

type field struct {
	id      string
	comment string
	typ     reflect.Type
	attributes
}

func (f *field) Id() string {
	return f.id
}

func (f *field) Comment() string {
	return f.comment
}

func (f *field) Type() reflect.Type {
	return f.typ
}

func (f *field) ToJsonObject() sdjson.Object {
	if f == nil {
		return nil
	}
	return sdjson.Object{
		"id":         f.id,
		"comment":    f.comment,
		"type":       f.typ.String(),
		"attributes": f.attributes.ensure(),
	}
}
