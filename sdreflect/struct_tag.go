package sdreflect

import (
	"github.com/fatih/structtag"
	"github.com/samber/lo"
)

type StructTagMap struct {
	tag  string
	keys map[string]string
}

func ParseStructTag(tag string) (StructTagMap, bool) {
	if tag == "" {
		return StructTagMap{}, true
	}
	tags, err := structtag.Parse(tag)
	if err != nil {
		return StructTagMap{}, false
	}
	if tags.Len() <= 0 {
		return StructTagMap{}, true
	}
	m := StructTagMap{
		tag:  tag,
		keys: make(map[string]string),
	}
	for _, k := range tags.Tags() {
		m.keys[k.Key] = k.Value()
	}
	return m, true
}

func (m StructTagMap) String() string {
	return m.tag
}

func (m StructTagMap) Len() int {
	return len(m.keys)
}

func (m StructTagMap) Keys() []string {
	return lo.Keys(m.keys)
}

func (m StructTagMap) Lookup(k string) (string, bool) {
	v, ok := m.keys[k]
	return v, ok
}

func (m StructTagMap) Has(k string) bool {
	_, ok := m.Lookup(k)
	return ok
}

func (m StructTagMap) HasOne(keys ...string) bool {
	for _, k := range keys {
		if m.Has(k) {
			return true
		}
	}
	return false
}

func (m StructTagMap) HasAll(keys ...string) bool {
	for _, k := range keys {
		if !m.Has(k) {
			return false
		}
	}
	return true
}

func (m StructTagMap) Get(k string) string {
	v, _ := m.Lookup(k)
	return v
}

func (m StructTagMap) GetOrDefault(k string, def string) string {
	v, ok := m.Lookup(k)
	if !ok {
		return def
	}
	return v
}

func (m StructTagMap) First(keys ...string) string {
	if len(keys) <= 0 {
		return ""
	}
	for _, k := range keys {
		v, ok := m.Lookup(k)
		if ok {
			return v
		}
	}
	return ""
}

func (m StructTagMap) FirstOrDefault(keys []string, def string) string {
	if len(keys) <= 0 {
		return def
	}
	for _, k := range keys {
		v, ok := m.Lookup(k)
		if ok {
			return v
		}
	}
	return def
}
