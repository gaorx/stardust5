package sdcheck

// checker

type Checker interface {
	Check() error
}

// checker func

type CheckerFunc func() error

func (f CheckerFunc) Check() error {
	if f == nil {
		return nil
	}
	return f()
}
