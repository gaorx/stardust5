package sdblueprint

import (
	"github.com/gaorx/stardust5/sdcodegen/sdgengo"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdstrings"
	"strconv"
)

type AttributeValue string

type Attributes interface {
	Has(k string) bool
	Lookup(k string) (AttributeValue, bool)
	Get(k string) AttributeValue
	First(keys []string) AttributeValue
}

var _ Attributes = attributes{}

type attributes map[string]string

func (attrs attributes) Has(k string) bool {
	_, ok := attrs[k]
	return ok
}

func (attrs attributes) Lookup(k string) (AttributeValue, bool) {
	v, ok := attrs[k]
	if !ok {
		return "", false
	}
	return AttributeValue(v), true
}

func (attrs attributes) Get(k string) AttributeValue {
	v, ok := attrs[k]
	if !ok {
		return ""
	}
	return AttributeValue(v)
}

func (attrs attributes) First(keys []string) AttributeValue {
	for _, k := range keys {
		v, ok := attrs[k]
		if ok {
			return AttributeValue(v)
		}
	}
	return ""
}

func (attrs attributes) addAttr(k, v string) {
	attrs[k] = v
}

func (attrs attributes) ensure() attributes {
	if attrs == nil {
		return attributes{}
	}
	return attrs
}

func (v AttributeValue) IfPresented(f func(v string)) {
	if v != "" && f != nil {
		f(string(v))
	}
}

func (v AttributeValue) AsStr() string {
	return string(v)
}

func (v AttributeValue) AsInt(def int) int {
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(string(v))
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return i
}

func (v AttributeValue) AsInt64(def int64) int64 {
	if v == "" {
		return def
	}
	i, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return i
}

func (v AttributeValue) AsUint(def uint) uint {
	if v == "" {
		return def
	}
	i, err := strconv.ParseUint(string(v), 10, 64)
	if err != nil {
		return def
	}
	return uint(i)
}

func (v AttributeValue) AsUint64(def uint64) uint64 {
	if v == "" {
		return def
	}
	i, err := strconv.ParseUint(string(v), 10, 64)
	if err != nil {
		return def
	}
	return i
}

func (v AttributeValue) AsFloat64(def float64) float64 {
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return f
}

func (v AttributeValue) AsBool(def bool) bool {
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(string(v))
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return b
}

func (v AttributeValue) AsSlice(sep string) []string {
	return sdstrings.SplitNonempty(string(v), sep, true)
}

func appendStructFieldTagsByAttrs(tags []sdgengo.FieldTag, attrs Attributes, keys ...string) []sdgengo.FieldTag {
	for _, k := range keys {
		v := attrs.Get(k).AsStr()
		if v != "" {
			tags = append(tags, sdgengo.FieldTag{K: k, V: v})
		}
	}
	return tags
}
