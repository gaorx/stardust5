package sdsql

import (
	"fmt"
	"slices"
	"strings"
)

type Page struct {
	Num   int
	Size  int
	base0 bool
}

type PagingResult[T any] struct {
	Rows      []T
	NumRows   int
	PageSize  int
	PageNum   int
	PageTotal int
}

func Page0(num, size int) Page {
	return Page{Num: num, Size: size, base0: true}
}

func Page1(num, size int) Page {
	return Page{Num: num, Size: size, base0: false}
}

func (p Page) WithDefaultSize(defaultSize int) Page {
	if p.Size <= 0 {
		p.Size = defaultSize
	}
	return p
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

func PageRows[T any](rows []T, page Page) *PagingResult[T] {
	limit, offset := page.LimitOffset()
	start, end := offset, offset+limit
	numRows := len(rows)
	if start > end {
		start, end = end, start
	}
	if start > numRows {
		start = numRows
	}
	if start < 0 {
		start = 0
	}

	if end > numRows {
		end = numRows
	}
	if end < 0 {
		end = 0
	}
	return &PagingResult[T]{
		Rows:      slices.Clone(rows[start:end]),
		NumRows:   numRows,
		PageSize:  limit,
		PageNum:   page.Num,
		PageTotal: (numRows + limit - 1) / limit,
	}
}

func NewPagingResultTo[T, R any](pr *PagingResult[T], rows []R) *PagingResult[R] {
	return &PagingResult[R]{
		Rows:      rows,
		NumRows:   pr.NumRows,
		PageSize:  pr.PageSize,
		PageNum:   pr.PageNum,
		PageTotal: pr.PageTotal,
	}
}
