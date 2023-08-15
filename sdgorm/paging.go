package sdgorm

import (
	"database/sql"
	"fmt"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdsqlparser"
	"gorm.io/gorm"
	"strings"
)

type Page struct {
	Num   int
	Size  int
	base0 bool
}

func Page0(num, size int) Page {
	return Page{Num: num, Size: size, base0: true}
}

func Page1(num, size int) Page {
	return Page{Num: num, Size: size, base0: false}
}

func (p Page) Scope() func(*gorm.DB) *gorm.DB {
	limit, offset := p.LimitOffset()
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Limit(limit).Offset(offset)
	}
}

func (p Page) Sql() string {
	limit, offset := p.LimitOffset()
	return fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)
}

func (p Page) Wrap(rawSql string) string {
	if strings.TrimSpace(rawSql) == "" {
		return ""
	}
	return rawSql + p.Sql()
}

func (p Page) LimitOffset() (int, int) {
	const maxLimit = 1000000
	limit, num := p.Size, p.TrimNum()
	if limit <= 0 {
		limit = maxLimit
	}
	if p.base0 {
		return limit, num * limit
	} else {
		return limit, (num - 1) * limit
	}
}

func (p Page) TrimNum() int {
	num := p.Num
	if p.base0 {
		if num <= 0 {
			num = 0
		}
	} else {
		if num <= 1 {
			num = 1
		}
	}
	return num
}

type PagingResult[T any] struct {
	Rows      []T
	NumRows   int
	PageSize  int
	PageNum   int
	PageTotal int
}

func FindPaging[T any](builder func() *gorm.DB, p Page) (*PagingResult[T], error) {
	var rows []T
	dbr := builder().Scopes(p.Scope()).Find(&rows)
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
	return &PagingResult[T]{
		Rows:      rows,
		NumRows:   numRows,
		PageSize:  limit,
		PageNum:   p.TrimNum(),
		PageTotal: pageTotal,
	}, nil
}

func RawPaging[T any](tx *gorm.DB, selectSql string, args map[string]any, p Page) (*PagingResult[T], error) {
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

	return &PagingResult[T]{
		Rows:      rows,
		NumRows:   numRows,
		PageSize:  limit,
		PageNum:   p.TrimNum(),
		PageTotal: pageTotal,
	}, nil
}
