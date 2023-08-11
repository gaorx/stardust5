package sdcall

import (
	"github.com/gaorx/stardust5/sderr"
)

func Safe(action func()) (err error) {
	if action == nil {
		err = nil
		return
	}

	defer func() {
		if err0 := recover(); err0 != nil {
			err = sderr.AsErr(err0)
		}
	}()
	action()
	return
}
