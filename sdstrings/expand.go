package sdstrings

import (
	"os"
)

func ExpandShellLike(s string, mapping func(string) string) string {
	return os.Expand(s, mapping)
}

func ExpandShellLikeMap(s string, data map[string]string) string {
	return ExpandShellLike(s, func(k string) string {
		return data[k]
	})
}
