package sdtabledata

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/samber/lo"
	"path"
)

type table struct {
	Source
	rows []*row
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

func newTable(src Source) (*table, error) {
	src = src.Trim()
	items, err := loadItems(src)
	if err != nil {
		return nil, err
	}
	t := &table{Source: src}
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
	if err := readJsonFile(row.t.Root, path.Join(row.t.Dir, row.dataFn), &rowData); err != nil {
		return sderr.WithStack(err)
	}
	if rowData == nil {
		rowData = sdjson.Object{}
	}

	// store column files
	if len(row.columns) > 0 {
		storeFile := opts.StoreFile

		storeFileForUrl := func(fn string) (string, error) {
			if storeFile == nil {
				return "", nil
			}
			target, err := storeFile(row.t.Root, path.Join(row.t.Dir, fn))
			if err != nil {
				return "", sderr.WrapWith(err, "store column file error", fn)
			}
			return target, nil
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

	if opts.Modifier != nil {
		rowData = opts.Modifier.ModifyRow(rowData)
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
