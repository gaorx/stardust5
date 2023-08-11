package sdblueprint

import (
	"github.com/gaorx/stardust5/sdregexp"
	"github.com/gaorx/stardust5/sdstrings"
	"reflect"
	"regexp"
	"strconv"
)

type structTag reflect.StructTag

func (tag structTag) Get(key string) string {
	v, _ := tag.Lookup(key)
	return v
}

const supportMultilineTag = true

func (tag structTag) Lookup(key string) (value string, ok bool) {
	isSpace := func(b byte) bool {
		if supportMultilineTag {
			return b == ' ' || b == '\t' || b == '\r' || b == '\n'
		} else {
			return b == ' '
		}
	}

	// 由于GO的reflect.StructTag.Lookup不支持多行tag解析，修改其代码支持多行tag解析
	tag1 := tag
	for tag1 != "" {
		// Skip leading space.
		i := 0
		for i < len(tag1) && isSpace(tag1[i]) {
			i++
		}
		tag1 = tag1[i:]
		if tag1 == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag1) && tag1[i] > ' ' && tag1[i] != ':' && tag1[i] != '"' && tag1[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag1) || tag1[i] != ':' || tag1[i+1] != '"' {
			break
		}
		name := string(tag1[:i])
		tag1 = tag1[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag1) && tag1[i] != '"' {
			if tag1[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag1) {
			break
		}
		qvalue := string(tag1[:i+1])
		tag1 = tag1[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}

func (tag structTag) split(k string) []string {
	return sdstrings.SplitNonempty(tag.Get(k), ",", true)
}

func (tag structTag) splitForFlags(k string) (string, []string) {
	l := tag.split(k)
	if len(l) <= 0 {
		return "", nil
	} else if len(l) == 1 {
		return l[0], nil
	} else {
		return l[0], l[1:]
	}
}

func (tag structTag) toAttrs(attrs attributes, keys ...string) {
	if len(keys) <= 0 {
		return
	}
	for _, k := range keys {
		if k == "" {
			continue
		}
		if v, ok := tag.Lookup(k); ok {
			attrs.addAttr(k, v)
		}
	}
}

func (tag structTag) toAttrsForFlags(attrs attributes, keys ...string) {
	if len(keys) <= 0 {
		return
	}
	for _, k := range keys {
		if k == "" {
			continue
		}
		first, flags := tag.splitForFlags(k)
		if first != "" {
			attrs.addAttr(k, first)
		}
		if len(flags) > 0 {
			for _, flag := range flags {
				if flag != "" {
					attrs.addAttr(k+"."+flag, "true")
				}
			}
		}
	}
}

var columnReferencePatt = regexp.MustCompile(`^ *(?P<col>\w+) *<-> *(?P<refTable>\w+) *\. *(?P<refCol>\w+) *$`)

func (tag structTag) toForeignKey(k string) ([]string, string, []string, bool) {
	indexes := tag.split(k)
	if len(indexes) <= 0 {
		return nil, "", nil, false
	}
	var colRefs columnReferences
	for _, index := range indexes {
		m := sdregexp.FindStringSubmatchGroup(columnReferencePatt, index)
		if len(m) <= 0 {
			return nil, "", nil, false
		}
		colRefs = append(colRefs, columnReference{col: m["col"], refTable: m["refTable"], refCol: m["refCol"]})
	}
	return colRefs.forFK()
}

func (tag structTag) comment() string {
	return tag.Get("comment")
}

func (tag structTag) group() string {
	return tag.Get("group")
}
