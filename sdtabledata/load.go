package sdtabledata

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdslices"
	"github.com/samber/lo"
)

type LoadOptions struct {
	StoreFile        StoreFile
	IgnoreIllegalRow bool
}

func LoadAll(src Source, opts *LoadOptions) ([]sdjson.Object, error) {
	if src.IsNil() {
		return nil, sderr.New("nil source")
	}
	opts1 := lo.FromPtr(opts)
	t, err := newTable(src)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	err = t.absorb(&opts1)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return t.data(), nil
}

func LoadSome(src Source, rows []string, opts *LoadOptions) ([]sdjson.Object, error) {
	if src.IsNil() {
		return nil, sderr.New("nil source")
	}
	opts1 := lo.FromPtr(opts)
	t, err := newTable(src)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	var rowsData []sdjson.Object
	for _, rowId := range rows {
		row1 := t.getRow(rowId)
		if row1 == nil {
			if opts.IgnoreIllegalRow {
				continue
			} else {
				return nil, sderr.NewWith("not found row in table data", row1)
			}
		}
		err = row1.absorb(&opts1)
		if err != nil {
			return nil, sderr.WithStack(err)
		}
		rowsData = append(rowsData, row1.data)
	}
	return sdslices.Ensure(rowsData), nil
}

func LoadOne(src Source, row string, opts *LoadOptions) (sdjson.Object, error) {
	if src.IsNil() {
		return nil, sderr.New("nil source")
	}
	opts1 := lo.FromPtr(opts)
	t, err := newTable(src)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	row1 := t.getRow(row)
	if row1 == nil {
		return nil, sderr.NewWith("not found row in table data", row)
	}
	err = row1.absorb(&opts1)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return row1.data, nil
}
