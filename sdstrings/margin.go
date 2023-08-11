package sdstrings

import (
	"strings"
	"unicode"
)

func TrimMargin(s string, marginPrefix string) string {
	lines := strings.Split(s, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		line = strings.TrimLeftFunc(line, func(c rune) bool {
			return unicode.IsSpace(c)
		})
		line = strings.TrimRightFunc(line, func(c rune) bool {
			return c == '\r'
		})
		line = strings.TrimPrefix(line, marginPrefix)
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}
