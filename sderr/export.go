package sderr

import (
	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
)

// export types
type MultipleError = multierror.Error

// export
var (
	New    = eris.New
	Newf   = eris.Errorf
	Wrap   = eris.Wrap
	Wrapf  = eris.Wrapf
	Unwrap = eris.Unwrap
	Cause  = eris.Cause
	Is     = eris.Is
	Try    = eris.As

	// multiple
	Append = multierror.Append
)

// misc

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if len(eris.StackFrames(err)) > 0 {
		return err
	}
	return Wrap(err, "")
}

func TryT[E error](err error) (E, bool) {
	var e E
	if Try(err, &e) {
		return e, true
	} else {
		return e, false
	}
}

func AsErr(v any) error {
	switch err := v.(type) {
	case nil:
		return nil
	case error:
		return err
	case string:
		return New(err)
	default:
		return Newf("%v", err)
	}
}
