package sdgorm

import (
	"database/sql"
	"fmt"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdreflect"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
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

func Exec(tx *gorm.DB, q string, values ...any) (int64, error) {
	dbr := tx.Exec(q, values...)
	if dbr.Error != nil {
		return 0, dbr.Error
	}
	return dbr.RowsAffected, nil
}

func Raw[T any](tx *gorm.DB, q string, values ...any) (T, error) {
	var r T
	dbr := tx.Raw(q, values...).Scan(&r)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
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
		return nil, dbr.Error
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

func Exists[T any](tx *gorm.DB, conds ...any) (bool, error) {
	q := tx.ToSQL(func(tx1 *gorm.DB) *gorm.DB {
		var r T
		return tx1.Find(&r, conds...)
	})
	q = fmt.Sprintf("SELECT EXISTS(%s)", q)
	return Raw[bool](tx, q)
}

func Create(tx *gorm.DB, row any) (int64, error) {
	dbr := tx.Create(row)
	if dbr.Error != nil {
		return 0, dbr.Error
	}
	return dbr.RowsAffected, nil
}

func CreateAndTake[T any](tx *gorm.DB, row T, q any, args ...any) (T, error) {
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

func ModifyAndTake[T any](tx *gorm.DB, modifier func(T) T, q any, args ...any) (T, error) {
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

func UpdateColumnsAndTake[T any](tx *gorm.DB, colVals map[string]any, q any, args ...any) (T, error) {
	_, err := UpdateColumns(tx, lo.Empty[T](), colVals, q, args...)
	if err != nil {
		return lo.Empty[T](), err
	}
	return Take[T](tx, append([]any{q}, args...)...)
}

type CreateInBatchesOptions struct {
	Table     string // 可以指定表名，如果为空，则使用model中表名
	Clear     bool   // 在插入数据前删除表中所有数据(危险操作)
	Overwrite bool   // 如果为true，如果表中存在相同ID的行，则覆盖掉原来的行；否则不修改任何数据
}

func CreateInBatches(tx *gorm.DB, rows any, opts *CreateInBatchesOptions) (int64, error) {
	if rows == nil {
		return 0, nil
	}
	rowsVal := sdreflect.ValueOf(rows)
	if rowsVal.Kind() != reflect.Slice && rowsVal.Kind() != reflect.Array {
		return 0, sderr.New("illegal rows type")
	}

	nRows := rowsVal.Len()
	if nRows <= 0 {
		return 0, nil
	}

	rowTyp := rowsVal.Type().Elem()
	opts1 := lo.FromPtr(opts)
	txByTableNameOrModel := func(tx *gorm.DB) *gorm.DB {
		if opts1.Table != "" {
			return tx.Table(opts1.Table)
		} else {
			return tx.Model(rowsVal.Index(0).Interface())
		}
	}

	if opts1.Clear {
		emptyRow := reflect.New(rowTyp).Elem()
		dbr := txByTableNameOrModel(tx).
			Session(&gorm.Session{AllowGlobalUpdate: true}).
			Clauses().Delete(emptyRow)
		if dbr.Error != nil {
			return 0, dbr.Error
		}
	}

	rowsAffected := int64(0)
	for i := 0; i < nRows; i++ {
		row := rowsVal.Index(i).Interface()
		var dbr *gorm.DB
		if opts1.Overwrite {
			dbr = txByTableNameOrModel(tx).
				Clauses(clause.OnConflict{UpdateAll: true}).
				Create(row)
		} else {
			dbr = txByTableNameOrModel(tx).
				Clauses(clause.OnConflict{DoNothing: true}).
				Create(row)
		}
		if dbr.Error != nil {
			return rowsAffected, dbr.Error
		}
		rowsAffected += dbr.RowsAffected
	}
	return rowsAffected, nil
}
