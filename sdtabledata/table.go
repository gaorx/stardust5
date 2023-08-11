package sdtabledata

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdobjectstore"
	"github.com/samber/lo"
	"path/filepath"
)

type table struct {
	dirAbs string
	rows   []*row
}

type row struct {
	t       *table
	id      string
	dataFn  string
	columns []*column
	data    sdjson.Object
}

type column struct {
	id   string
	fn   string
	subs map[string]string
}

func newTable(dirAbs string) (*table, error) {
	items, err := loadItems(dirAbs)
	if err != nil {
		return nil, err
	}
	t := &table{dirAbs: dirAbs}
	for _, item := range items {
		switch item.kind {
		case fileMeta:
			// 目前不支持meta文件，所以先忽略
		case fileRow:
			t.ensureRow(item.row).dataFn = item.filename
		case fileColumn:
			t.ensureRow(item.row).ensureColumn(item.column).fn = item.filename
		case fileColumnSub:
			t.ensureRow(item.row).ensureColumn(item.column).addSub(item.sub, item.filename)
		}
	}
	return t, nil
}

func (t *table) getRow(rowId string) *row {
	for _, row := range t.rows {
		if row.id == rowId {
			return row
		}
	}
	return nil
}

func (t *table) data() []sdjson.Object {
	var data []sdjson.Object
	for _, row := range t.rows {
		data = append(data, row.data)
	}
	return data
}

func (t *table) listRowIds() []string {
	return lo.Map(t.rows, func(row1 *row, _ int) string {
		return row1.id
	})
}

func (t *table) absFn(fn string) string {
	if filepath.IsAbs(fn) {
		return fn
	} else {
		return filepath.Join(t.dirAbs, fn)
	}
}

func (t *table) ensureRow(rowId string) *row {
	for _, row := range t.rows {
		if row.id == rowId {
			return row
		}
	}
	row := &row{t: t, id: rowId}
	t.rows = append(t.rows, row)
	return row
}

func (t *table) absorb(opts *LoadOptions) error {
	for _, row := range t.rows {
		if err := row.absorb(opts); err != nil {
			return err
		}
	}
	return nil
}

func (row *row) ensureColumn(colId string) *column {
	for _, col := range row.columns {
		if col.id == colId {
			return col
		}
	}
	col := &column{id: colId}
	row.columns = append(row.columns, col)
	return col
}

func (row *row) absorb(opts *LoadOptions) error {
	var rowData sdjson.Object
	if err := readJsonFile(row.t.absFn(row.dataFn), &rowData); err != nil {
		return sderr.WithStack(err)
	}
	if rowData == nil {
		rowData = sdjson.Object{}
	}

	// store column files
	if len(row.columns) > 0 {
		store, objectName, httpUrl := opts.Store, opts.StoreObjectName, opts.StoreHttpUrl
		if store == nil {
			return sderr.New("nil object store")
		}

		storeFileForUrl := func(fn string) (string, error) {
			fnAbs := row.t.absFn(fn)
			target, err := store.Store(sdobjectstore.File(fnAbs, ""), objectName)
			if err != nil {
				return "", sderr.WrapWith(err, "store column file error", fnAbs)
			}
			if httpUrl {
				return target.Url(), nil
			} else {
				return target.HttpsUrl(), nil
			}
		}

		for _, col := range row.columns {
			if len(col.subs) > 0 {
				colObj := sdjson.Object{}
				for sub, subFn := range col.subs {
					url, err := storeFileForUrl(subFn)
					if err != nil {
						return err
					}
					colObj[sub] = url
				}
				rowData[col.id] = sdjson.MarshalStringDef(colObj, "{}")
			} else {
				url, err := storeFileForUrl(col.fn)
				if err != nil {
					return err
				}
				rowData[col.id] = url
			}
		}
	}

	row.data = rowData
	return nil
}

func (col *column) addSub(sub, fn string) {
	if col.subs == nil {
		col.subs = map[string]string{}
	}
	col.subs[sub] = fn
}
