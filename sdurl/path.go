package sdurl

import (
	"strings"

	"github.com/gaorx/stardust5/sdstrings"
	"github.com/samber/lo"
)

func JoinPath(paths ...string) string {
	var segments []string
	for _, p := range paths {
		segments0 := sdstrings.SplitNonempty(p, "/", false)
		segments = append(segments, segments0...)
	}
	return "/" + strings.Join(lo.Map(segments, func(s string, _ int) string {
		return QueryEscape(s, EscapeEncodePath)
	}), "/")
}
