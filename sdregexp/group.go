package sdregexp

import (
	"regexp"
)

func FindStringSubmatchGroup(exp *regexp.Regexp, s string) map[string]string {
	match := exp.FindStringSubmatch(s)
	if match == nil {
		return nil
	}
	r := make(map[string]string)
	for i, name := range exp.SubexpNames() {
		if i != 0 && name != "" {
			r[name] = match[i]
		}
	}
	return r
}
