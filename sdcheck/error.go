package sdcheck

import (
	"fmt"
	"github.com/gaorx/stardust5/sderr"
)

const defaultFailedMessage = "check failed"

func errorOf(message any) error {
	switch v := message.(type) {
	case nil:
		return sderr.New(defaultFailedMessage)
	case string:
		return sderr.New(v)
	case error:
		return v
	case func() error:
		return v()
	case func() string:
		return sderr.New(v())
	case fmt.Stringer:
		return sderr.New(v.String())
	default:
		return sderr.New(defaultFailedMessage)
	}
}
