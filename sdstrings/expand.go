package sdstrings

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
	"os"
)

type ExpandMapper = func(string) string

func ExpandShellLike(s string, mapper ExpandMapper) string {
	return os.Expand(s, mapper)
}

func ExpandShellLikeV(s string, vars map[string]string) string {
	return ExpandShellLike(s, func(k string) string {
		return vars[k]
	})
}

func ExpandShellLikeN(s string, mapper any, others ...any) string {
	return ExpandShellLike(s, MergeExpandMappers(mapper, others))
}

func MergeExpandMappers(mapper any, others ...any) ExpandMapper {
	// normalize mapper
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

	// merge mappers
	var finals []func(string) string
	add := func(m any) {
		if m1 := normalize(m); m1 != nil {
			finals = append(finals, m1)
		}
	}
	add(mapper)
	lo.ForEach(others, func(x any, _ int) { add(x) })
	merged := func(k string) string {
		for _, m := range finals {
			v := m(k)
			if v != "" {
				return v
			}
		}
		return ""
	}
	return merged
}
