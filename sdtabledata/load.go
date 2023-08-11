package sdtabledata

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdobjectstore"
	"path/filepath"
)

type LoadOptions struct {
	Store            sdobjectstore.Store
	StoreObjectName  string
	StoreHttpUrl     bool
	IgnoreIllegalRow bool
}

func LoadAll(dir string, rowsPtr any, opts LoadOptions) error {
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return sderr.WrapWith(err, "get data absolute directory error", dir)
	}
	t, err := newTable(dirAbs)
	if err != nil {
		return sderr.WithStack(err)
	}
	err = t.absorb(&opts)
	if err != nil {
		return sderr.WithStack(err)
	}
	return convertByJson(t.data(), rowsPtr)
}

func LoadSome(dir string, rowsPtr any, rows []string, opts LoadOptions) error {
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return sderr.WrapWith(err, "get data absolute directory error", dir)
	}
	t, err := newTable(dirAbs)
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
		err = row1.absorb(&opts)
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

func LoadOne(dir string, row string, rowPtr any, opts LoadOptions) error {
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return sderr.WrapWith(err, "get data absolute directory error", dir)
	}
	t, err := newTable(dirAbs)
	if err != nil {
		return sderr.WithStack(err)
	}
	row1 := t.getRow(row)
	if row1 == nil {
		return sderr.NewWith("not found row in table data", row)
	}
	err = row1.absorb(&opts)
	if err != nil {
		return sderr.WithStack(err)
	}
	return convertByJson(row1.data, rowPtr)
}
