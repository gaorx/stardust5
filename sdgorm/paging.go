package sdgorm

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type Page struct {
	Num   int
	Size  int
	base0 bool
}

func PageOf(num, size int) Page {
	return Page{Num: num, Size: size, base0: true}
}

func Page1Of(num, size int) Page {
	return Page{Num: num, Size: size, base0: false}
}

func (p Page) Scope() func(*gorm.DB) *gorm.DB {
	if p.base0 {
		num, size := p.Num, p.Size
		if num <= 0 {
			num = 0
		}
		if size > 0 {
			return func(tx *gorm.DB) *gorm.DB {
				return tx.Offset(num * size).Limit(size)
			}
		} else {
			return func(tx *gorm.DB) *gorm.DB {
				return tx
			}
		}
	} else {
		num, size := p.Num, p.Size
		if num <= 1 {
			num = 1
		}
		if size > 0 {
			return func(tx *gorm.DB) *gorm.DB {
				return tx.Offset((num - 1) * size).Limit(size)
			}
		} else {
			return func(tx *gorm.DB) *gorm.DB {
				return tx
			}
		}
	}
}

func (p Page) Sql() string {
	if p.base0 {
		num, size := p.Num, p.Size
		if num <= 0 {
			num = 0
		}
		if size > 0 {
			return fmt.Sprintf(" LIMIT %d OFFSET %d ", size, num*size)
		} else {
			return ""
		}
	} else {
		num, size := p.Num, p.Size
		if num <= 0 {
			num = 0
		}
		if size > 0 {
			return fmt.Sprintf(" LIMIT %d OFFSET %d ", size, (num-1)*size)
		} else {
			return ""
		}
	}
}

func (p Page) Wrap(rawSql string) string {
	if strings.TrimSpace(rawSql) == "" {
		return ""
	}
	return rawSql + p.Sql()
}
