package sdcheck

func funcOf(c Checker) CheckerFunc {
	if f, ok := c.(CheckerFunc); ok {
		return f
	} else {
		return func() error {
			return c.Check()
		}
	}
}
