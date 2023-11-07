package sdcheck

import (
	"github.com/gaorx/stardust5/sdreflect"
	"reflect"
)

// true/false

func True(b bool, message any) CheckerFunc {
	return func() error {
		if !b {
			return errorOf(message)
		}
		return nil
	}
}

func False(b bool, message any) CheckerFunc {
	return func() error {
		if b {
			return errorOf(message)
		}
		return nil
	}
}

// All/And/Or

func Not(c Checker, message any) CheckerFunc {
	return func() error {
		if c == nil {
			return errorOf(message)
		}
		if err := c.Check(); err == nil {
			return errorOf(message)
		} else {
			return nil
		}
	}
}

func All(checkers ...Checker) CheckerFunc {
	if len(checkers) == 0 {
		return CheckerFunc(nil)
	}
	if len(checkers) == 1 {
		return funcOf(checkers[0])
	}
	return func() error {
		for _, c := range checkers {
			if c != nil {
				if err := c.Check(); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func And(checkers []Checker, message any) CheckerFunc {
	if len(checkers) == 0 {
		return CheckerFunc(nil)
	}
	if len(checkers) == 1 {
		return funcOf(checkers[0])
	}
	return func() error {
		for _, c := range checkers {
			if c != nil {
				if err := c.Check(); err != nil {
					return errorOf(message)
				}
			}
		}
		return nil
	}
}

func Or(checkers []Checker, message any) CheckerFunc {
	if len(checkers) == 0 {
		return CheckerFunc(nil)
	}
	if len(checkers) == 1 {
		return funcOf(checkers[0])
	}
	return func() error {
		for _, c := range checkers {
			if c != nil {
				if err := c.Check(); err == nil {
					return nil
				}
			}
		}
		return errorOf(message)
	}
}

// if

func If(enabled bool, checker Checker) CheckerFunc {
	if !enabled {
		return CheckerFunc(nil)
	}
	return funcOf(checker)
}

// for

type ForFunc[T any] func() (T, error)

func For[T any](f ForFunc[T], ptr *T) CheckerFunc {
	return func() error {
		r, err := f()
		if err != nil {
			return err
		}
		if ptr != nil {
			*ptr = r
		}
		return nil
	}
}

// other

func Required(v any, message any) CheckerFunc {
	return func() error {
		v := sdreflect.ValueOf(v)
		k := v.Kind()
		if !v.IsValid() {
			return errorOf(message)
		}
		if (k == reflect.Pointer || k == reflect.Func) && v.IsNil() {
			return errorOf(message)
		}
		if (k == reflect.Slice || k == reflect.Array || k == reflect.Map) && (v.IsNil() || v.Len() <= 0) {
			return errorOf(message)
		}
		if v.IsZero() {
			return errorOf(message)
		}
		return nil
	}
}

func Len(v any, minLen, maxLen int, message any) CheckerFunc {
	if maxLen < minLen {
		minLen, maxLen = maxLen, minLen
	}
	return func() error {
		if n := sdreflect.ValueOf(v).Len(); n < minLen || n > maxLen {
			return errorOf(message)
		}
		return nil
	}
}
