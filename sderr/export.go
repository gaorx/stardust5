package sderr

import (
	stderr "errors"
	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
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
	As     = eris.As
)

var (
	Append = multierror.Append
)

// misc

func Sentinel(text string) error {
	return stderr.New(text)
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if len(eris.StackFrames(err)) > 0 {
		return err
	}
	return Wrap(err, "")
}

func AsT[E error](err error) (E, bool) {
	return lo.ErrorsAs[E](err)
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
