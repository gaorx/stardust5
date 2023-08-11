package sdhttpua

var All = make([]*UA, 0)

func init() {
	for _, raw := range rawUserAgents {
		ua, err := Parse(raw, nil)
		if err != nil {
			continue
		}
		All = append(All, ua)
	}
}

func Find(predicates ...Predicate) []*UA {
	var r []*UA
	for _, ua := range All {
		if isOK(ua, predicates) {
			r = append(r, ua)
		}
	}
	return r
}

func FindRaw(predicates ...Predicate) []string {
	var r []string
	for _, ua := range All {
		if isOK(ua, predicates) {
			r = append(r, ua.UA)
		}
	}
	return r
}

func isOK(ua *UA, predicates []Predicate) bool {
	for _, p := range predicates {
		if p != nil && !p(ua) {
			return false
		}
	}
	return true
}
