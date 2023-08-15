package sdgorm

import (
	"database/sql"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func Transaction[R any](db *gorm.DB, action func(tx *gorm.DB) (R, error), opts ...*sql.TxOptions) (R, error) {
	var r R
	err := db.Transaction(func(tx *gorm.DB) error {
		r0, err := action(tx)
		if err != nil {
			return err
		}
		r = r0
		return nil
	}, opts...)
	return r, err
}

func First[T any](tx *gorm.DB, conds ...any) (T, error) {
	var r T
	dbr := tx.First(&r, conds...)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
}

func Last[T any](tx *gorm.DB, conds ...any) (T, error) {
	var r T
	dbr := tx.Last(&r, conds...)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
}

func Take[T any](tx *gorm.DB, conds ...any) (T, error) {
	var r T
	dbr := tx.Take(&r, conds...)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
}

func Find[T any](tx *gorm.DB, conds ...any) ([]T, error) {
	var r []T
	dbr := tx.Find(&r, conds...)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
}

func Scan[T any](tx *gorm.DB) (T, error) {
	var r T
	dbr := tx.Scan(&r)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
}

func Create(tx *gorm.DB, row any) (int64, error) {
	dbr := tx.Create(row)
	if dbr.Error != nil {
		return 0, dbr.Error
	}
	return dbr.RowsAffected, nil
}

func CreateTake[T any](tx *gorm.DB, row T, q any, args ...any) (T, error) {
	var err error
	_, err = Create(tx, row)
	if err != nil {
		return lo.Empty[T](), err
	}
	created, err := Take[T](tx, append([]any{q}, args...)...)
	if err != nil {
		return lo.Empty[T](), err
	}
	return created, nil
}

func Modify[T any](tx *gorm.DB, modifier func(T) T, q any, args ...any) (int64, error) {
	if modifier == nil {
		return 0, nil
	}
	var row T
	dbr := tx.Where(q, args...).Take(&row)
	if dbr.Error != nil {
		return 0, dbr.Error
	}
	dbr = tx.Where(q, args...).Save(modifier(row))
	if dbr.Error != nil {
		return 0, dbr.Error
	}
	return dbr.RowsAffected, nil
}

func ModifyTake[T any](tx *gorm.DB, modifier func(T) T, q any, args ...any) (T, error) {
	_, err := Modify[T](tx, modifier, q, args...)
	if err != nil {
		return lo.Empty[T](), err
	}
	return Take[T](tx, append([]any{q}, args...)...)
}

func UpdateColumns(tx *gorm.DB, model any, colVals map[string]any, q any, args ...any) (int64, error) {
	if len(colVals) <= 0 {
		return 0, nil
	}
	if tableName, ok := model.(string); ok {
		tx = tx.Table(tableName)
	} else {
		tx = tx.Model(model)
	}
	dbr := tx.Where(q, args...).Updates(colVals)
	if dbr.Error != nil {
		return 0, dbr.Error
	}
	return dbr.RowsAffected, nil
}

func UpdateColumnsTake[T any](tx *gorm.DB, colVals map[string]any, q any, args ...any) (T, error) {
	_, err := UpdateColumns(tx, lo.Empty[T](), colVals, q, args...)
	if err != nil {
		return lo.Empty[T](), err
	}
	return Take[T](tx, append([]any{q}, args...)...)
}
