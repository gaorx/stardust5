package sdstrings

import (
	"github.com/samber/lo"
	"strings"
)

func JoinFunc[T any](collection []T, sep string, iteratee func(item T, index int) string) string {
	l := lo.Map(collection, iteratee)
	return strings.Join(l, sep)
}
