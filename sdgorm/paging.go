package sdgorm

import (
	"database/sql"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdsql"
	"github.com/gaorx/stardust5/sdsqlparser"
	"gorm.io/gorm"
)

func PageScope(p sdsql.Page) func(*gorm.DB) *gorm.DB {
	limit, offset := p.LimitOffset()
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Limit(limit).Offset(offset)
	}
}

func FindPaging[T any](builder func() *gorm.DB, p sdsql.Page) (*sdsql.PagingResult[T], error) {
	var rows []T
	dbr := builder().Scopes(PageScope(p)).Find(&rows)
	if dbr.Error != nil {
		return nil, dbr.Error
	}
	var numRows int
	dbr = builder().Select("COUNT(*)").Scan(&numRows)
	if dbr.Error != nil {
		return nil, dbr.Error
	}

	limit, _ := p.LimitOffset()
	var pageTotal int
	if numRows%limit == 0 {
		pageTotal = numRows / limit
	} else {
		pageTotal = numRows/limit + 1
	}
	return &sdsql.PagingResult[T]{
		Rows:      rows,
		NumRows:   numRows,
		PageSize:  limit,
		PageNum:   p.TrimNum(),
		PageTotal: pageTotal,
	}, nil
}

func RawPaging[T any](tx *gorm.DB, selectSql string, args map[string]any, p sdsql.Page) (*sdsql.PagingResult[T], error) {
	q1, ok := sdsqlparser.SqlWithLimit(selectSql, "@__limit", "@__offset")
	if !ok {
		return nil, sderr.NewWith("parse select sql error", "1")
	}
	q2, ok := sdsqlparser.SqlForCount(selectSql)
	if !ok {
		return nil, sderr.NewWith("parse select sql error", "2")
	}
	limit, offset := p.LimitOffset()
	var args1 []any
	for k, v := range args {
		args1 = append(args1, sql.Named(k, v))
	}
	args1 = append(args1, sql.Named("__limit", limit), sql.Named("__offset", offset))
	var rows []T
	dbr := tx.Raw(q1, args1...).Find(&rows)
	if dbr.Error != nil {
		return nil, dbr.Error
	}
	var numRows int
	dbr = tx.Raw(q2, args1...).Scan(&numRows)
	if dbr.Error != nil {
		return nil, dbr.Error
	}

	var pageTotal int
	if numRows%limit == 0 {
		pageTotal = numRows / limit
	} else {
		pageTotal = numRows/limit + 1
	}

	return &sdsql.PagingResult[T]{
		Rows:      rows,
		NumRows:   numRows,
		PageSize:  limit,
		PageNum:   p.TrimNum(),
		PageTotal: pageTotal,
	}, nil
}
