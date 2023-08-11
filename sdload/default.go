package sdload

func TextDef(loc, def string) string {
	s, err := Text(loc)
	if err != nil {
		return def
	}
	return s
}
