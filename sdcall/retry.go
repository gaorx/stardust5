package sdcall

import (
	"github.com/gaorx/stardust5/sderr"
)

func Retry(maxRetries int, action func() error) error {
	err0 := action()
	if err0 == nil {
		return nil
	}
	err := sderr.Append(err0)
	for i := 1; i <= maxRetries; i++ {
		err0 := action()
		if err0 != nil {
			err = sderr.Append(err, err0)
		} else {
			return nil
		}
	}
	return err
}
