package sdgorm

import (
	"github.com/gaorx/stardust5/sderr/sdnotfounderr"
)

func RowExists[T any](_ T, err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	if sdnotfounderr.Is(err) {
		return false, nil
	}
	return false, err
}
