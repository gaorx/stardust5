package sdcheck

// checker

type Checker interface {
	Check() error
}

// checker func

type Func func() error

func (f Func) Check() error {
	if f == nil {
		return nil
	}
	return f()
}
