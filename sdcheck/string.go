package sdcheck

import (
	"regexp"
	"strings"
)

func MatchRegexp(s, pattern string, message any) CheckerFunc {
	return func() error {
		if matched, err := regexp.MatchString(pattern, s); err != nil {
			return errorOf(message)
		} else {
			if matched {
				return nil
			} else {
				return errorOf(message)
			}
		}
	}
}

func MatchRegexpPattern(s string, pattern *regexp.Regexp, message any) CheckerFunc {
	return func() error {
		if matched := pattern.MatchString(s); matched {
			return nil
		} else {
			return errorOf(message)
		}
	}
}

func HasSub(s string, substr string, message any) CheckerFunc {
	return func() error {
		if !strings.Contains(s, substr) {
			return errorOf(message)
		}
		return nil
	}
}

func HasPrefix(s string, prefix string, message any) CheckerFunc {
	return func() error {
		if !strings.HasPrefix(s, prefix) {
			return errorOf(message)
		}
		return nil
	}
}

func HasSuffix(s string, suffix string, message any) CheckerFunc {
	return func() error {
		if !strings.HasSuffix(s, suffix) {
			return errorOf(message)
		}
		return nil
	}
}
