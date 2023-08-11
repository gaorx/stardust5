package sdtabledata

import (
	"github.com/gaorx/stardust5/sderr"
	"path/filepath"
)

func ListRows(dir string) ([]string, error) {
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return nil, sderr.WrapWith(err, "get data absolute directory error", dir)
	}
	t, err := newTable(dirAbs)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return t.listRowIds(), nil
}
