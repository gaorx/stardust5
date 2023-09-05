package sdtabledata

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdobjectstore"
	"github.com/samber/lo"
)

type LoadOptions struct {
	Store            sdobjectstore.Store
	StoreObjectName  string
	StoreHttpUrl     bool
	IgnoreIllegalRow bool
}

func LoadAll(src Source, rowsPtr any, opts *LoadOptions) error {
	if src.IsNil() {
		return sderr.New("nil source")
	}
	opts1 := lo.FromPtr(opts)
	t, err := newTable(src)
	if err != nil {
		return sderr.WithStack(err)
	}
	err = t.absorb(&opts1)
	if err != nil {
		return sderr.WithStack(err)
	}
	return convertByJson(t.data(), rowsPtr)
}

func LoadSome(src Source, rowsPtr any, rows []string, opts *LoadOptions) error {
	if src.IsNil() {
		return sderr.New("nil source")
	}
	opts1 := lo.FromPtr(opts)
	t, err := newTable(src)
	if err != nil {
		return sderr.WithStack(err)
	}
	var rowsData []sdjson.Object
	for _, rowId := range rows {
		row1 := t.getRow(rowId)
		if row1 == nil {
			if opts.IgnoreIllegalRow {
				continue
			} else {
				return sderr.NewWith("not found row in table data", row1)
			}
		}
		err = row1.absorb(&opts1)
		if err != nil {
			return sderr.WithStack(err)
		}
		rowsData = append(rowsData, row1.data)
	}
	if rowsData == nil {
		rowsData = []sdjson.Object{}
	}
	return convertByJson(rowsData, rowsPtr)
}

func LoadOne(src Source, row string, rowPtr any, opts *LoadOptions) error {
	if src.IsNil() {
		return sderr.New("nil source")
	}
	opts1 := lo.FromPtr(opts)
	t, err := newTable(src)
	if err != nil {
		return sderr.WithStack(err)
	}
	row1 := t.getRow(row)
	if row1 == nil {
		return sderr.NewWith("not found row in table data", row)
	}
	err = row1.absorb(&opts1)
	if err != nil {
		return sderr.WithStack(err)
	}
	return convertByJson(row1.data, rowPtr)
}
