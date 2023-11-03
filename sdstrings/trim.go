package sdstrings

import (
	"strings"
)

func TrimPrefixes(s string, prefixes ...string) string {
	if len(prefixes) <= 0 {
		return s
	}
	for _, prefix := range prefixes {
		s = strings.TrimPrefix(s, prefix)
	}
	return s
}

func TrimSuffixes(s string, suffixes ...string) string {
	if len(suffixes) <= 0 {
		return s
	}
	for _, suffix := range suffixes {
		s = strings.TrimSuffix(s, suffix)
	}
	return s
}
