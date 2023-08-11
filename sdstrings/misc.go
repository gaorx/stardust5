package sdstrings

func EmptyAs(s string, def string) string {
	if s == "" {
		return def
	}
	return s
}
