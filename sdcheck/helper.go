package sdcheck

func funcOf(c Checker) Func {
	if f, ok := c.(Func); ok {
		return f
	} else {
		return func() error {
			return c.Check()
		}
	}
}
