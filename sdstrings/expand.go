package sdstrings

import (
	"os"
)

func ExpandShellLike(s string, mapping func(string) string) string {
	return os.Expand(s, mapping)
}
