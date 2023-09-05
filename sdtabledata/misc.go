package sdtabledata

import (
	"github.com/gaorx/stardust5/sderr"
)

func ListRows(src Source) ([]string, error) {
	if src.IsNil() {
		return nil, sderr.New("nil source")
	}
	t, err := newTable(src)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return t.listRowIds(), nil
}
