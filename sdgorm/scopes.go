package sdgorm

import (
	"gorm.io/gorm"
)

func ScopePaging0(pageNum, pageSize int) func(*gorm.DB) *gorm.DB {
	if pageNum <= 0 {
		pageNum = 0
	}
	if pageSize > 0 {
		return func(tx *gorm.DB) *gorm.DB {
			return tx.Offset(pageNum * pageSize).Limit(pageSize)
		}
	} else {
		return func(tx *gorm.DB) *gorm.DB {
			return tx
		}
	}
}

func ScopePaging1(pageNum, pageSize int) func(*gorm.DB) *gorm.DB {
	if pageNum <= 1 {
		pageNum = 1
	}
	if pageSize > 0 {
		return func(tx *gorm.DB) *gorm.DB {
			return tx.Offset((pageNum - 1) * pageSize).Limit(pageSize)
		}
	} else {
		return func(tx *gorm.DB) *gorm.DB {
			return tx
		}
	}
}
