package sdhttpua

type Predicate func(*UA) bool

func (p Predicate) Not() Predicate {
	return func(ua *UA) bool {
		return !p(ua)
	}
}

func Or(first Predicate, others ...Predicate) Predicate {
	if len(others) <= 0 {
		return first
	}
	return func(ua *UA) bool {
		if first != nil && first(ua) {
			return true
		}
		for _, other := range others {
			if other != nil && other(ua) {
				return true
			}
		}
		return false
	}
}

func PlatformIs(platform string, others ...string) Predicate {
	return func(ua *UA) bool {
		if ua.Platform == platform {
			return true
		}
		for _, other := range others {
			if ua.Platform == other {
				return true
			}
		}
		return false
	}
}

func OSIs(os string, others ...string) Predicate {
	return func(ua *UA) bool {
		if ua.OS == os {
			return true
		}
		for _, other := range others {
			if ua.OS == other {
				return true
			}
		}
		return false
	}
}

func IsMobile() Predicate {
	return func(ua *UA) bool {
		return ua.Mobile
	}
}
