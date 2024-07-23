package sdbun

import (
	"context"
	"github.com/gaorx/stardust5/sdsql"
	"github.com/uptrace/bun"
)

func PageApplier(p sdsql.Page) func(*bun.SelectQuery) *bun.SelectQuery {
	limit, offset := p.LimitOffset()
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Limit(limit).Offset(offset)
	}
}

func SelectPage[ROW any](ctx context.Context, db bun.IDB, p sdsql.Page, qfn func(*bun.SelectQuery) *bun.SelectQuery, postProcs ...sdsql.RowsProc[ROW]) (*sdsql.PagingResult[ROW], error) {
	var rows []ROW
	numRows, err := db.NewSelect().Apply(qfn).Apply(PageApplier(p)).Apply(modelApplier[*bun.SelectQuery, ROW]()).ScanAndCount(ctx, &rows)
	if err != nil {
		return nil, err
	}
	rows, err = sdsql.ProcRows(rows, postProcs...)
	if err != nil {
		return nil, err
	}
	limit, _ := p.LimitOffset()
	var pageTotal int
	if numRows%limit == 0 {
		pageTotal = numRows / limit
	} else {
		pageTotal = numRows/limit + 1
	}
	return &sdsql.PagingResult[ROW]{
		Rows:      rows,
		NumRows:   numRows,
		PageSize:  limit,
		PageNum:   p.TrimNum(),
		PageTotal: pageTotal,
	}, nil
}
