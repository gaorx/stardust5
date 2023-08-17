package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdslices"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/samber/lo"
	"path"
	"regexp"
	"slices"
	"strings"
)

type Object struct {
	id   string
	tags []string
}

const oidPublic = "public"

var Public = Object{id: oidPublic}

var pattObjectPart = regexp.MustCompile(`^(\w|\$|\{|}|-|\.|:|\*|\?|\^)*$`)

func O(id string, tags ...string) Object {
	id = strings.TrimSpace(id)
	if !pattObjectPart.MatchString(id) {
		panic(sderr.NewWith("create access control object error 1", id))
	}

	var tags1 []string
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			if !pattObjectPart.MatchString(tag) {
				panic(sderr.NewWith("create access control object error 2", tag))
			}
			tags1 = append(tags1, tag)
		}
	}

	tags1 = sdslices.Ensure(tags1)
	slices.Sort(tags1)
	return Object{
		id:   id,
		tags: tags1,
	}
}

func (o Object) Id() string {
	return o.id
}

func (o Object) Tags() []string {
	return o.tags
}

func (o Object) HasTags() bool {
	return len(o.tags) > 0
}

func (o Object) MatchTag(tag string) bool {
	for _, tag0 := range o.tags {
		if tag0 == tag {
			return true
		} else {
			ok, err := path.Match(tag0, tag)
			if err != nil {
				continue
			}
			if ok {
				return true
			}
		}
	}
	return false
}

func (o Object) MatchTagsAny(tags []string) bool {
	for _, tag := range tags {
		if o.MatchTag(tag) {
			return true
		}
	}
	return false
}

func (o Object) TagCount() int {
	return len(o.tags)
}

func (o Object) IsEmpty() bool {
	return o.id == "" && len(o.tags) <= 0
}

func (o Object) IsPublic() bool {
	return o.id == oidPublic && len(o.tags) <= 0
}

func (o Object) IsExpanded() bool {
	if strings.Contains(o.id, "$") {
		return false
	}
	for _, tag := range o.tags {
		if strings.Contains(tag, "$") {
			return false
		}
	}
	return true
}

func (o Object) String() string {
	var buf strings.Builder
	buf.WriteString(o.id)
	if len(o.tags) > 0 {
		buf.WriteString("[")
		for i, tag := range o.tags {
			if i > 0 {
				buf.WriteString("|")
			}
			buf.WriteString(tag)
		}
		buf.WriteString("]")
	}
	return buf.String()
}

func (o Object) Expand(mapper any, others ...any) Object {
	normalize := func(m any) func(string) string {
		switch x := m.(type) {
		case nil:
			return nil
		case func(string) string:
			return x
		case map[string]string:
			return func(k string) string {
				return x[k]
			}
		default:
			panic(sderr.New("export access control object error"))
		}
	}
	var finals []func(string) string
	add := func(m any) {
		if m1 := normalize(m); m1 != nil {
			finals = append(finals, m1)
		}
	}
	add(mapper)
	lo.ForEach(others, func(x any, _ int) { add(x) })
	if len(finals) <= 0 {
		return o
	}

	merged := func(k string) string {
		for _, m := range finals {
			v := m(k)
			if v != "" {
				return v
			}
		}
		return ""
	}
	return O(
		sdstrings.ExpandShellLike(o.id, merged),
		lo.Map(o.tags, func(tag string, _ int) string {
			return sdstrings.ExpandShellLike(tag, merged)
		})...,
	)
}
