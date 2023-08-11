package sdgorm

import (
	"database/sql"
	"fmt"
	"github.com/gaorx/stardust5/sdreflect"
	"github.com/samber/lo"

	"github.com/gaorx/stardust5/sderr"
	"gorm.io/gorm"
)

type TableOptions struct {
	IdColumn string
}

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
	if err != nil {
		return r, sderr.Wrap(err, "sdgorm transaction error")
	}
	return r, nil
}

func IdFromRow(tx *gorm.DB, t any) (any, error) {
	ids, err := GetPrimaryKeys(t, tx.NamingStrategy)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	if len(ids) <= 0 {
		return "", sderr.New("no primary key column")
	}
	if len(ids) > 1 {
		return "", sderr.New("multiple primary key column")
	}
	return ids[0], nil
}

func IdFromRowT[ID any](tx *gorm.DB, row any) (ID, error) {
	id, err := IdFromRow(tx, row)
	if err != nil {
		return lo.Empty[ID](), sderr.WithStack(err)
	}
	id1, ok := sdreflect.To[ID](id)
	if !ok {
		return lo.Empty[ID](), sderr.New("invalid ID type")
	}
	return id1, nil
}

func Create[T any](tx *gorm.DB, row T) error {
	dbr := tx.Create(row)
	if dbr.Error != nil {
		return dbr.Error
	}
	return nil
}

func GetById[T any, ID any](tx *gorm.DB, id ID, opts TableOptions) (T, error) {
	idCol, err := idColumnWithOptions[T](tx, opts)
	if err != nil {
		return lo.Empty[T](), sderr.WithStack(err)
	}

	var r T
	dbr := tx.First(&r, fmt.Sprintf("%s = ?", idCol), id)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
}

func GetByIds[T any, ID any](tx *gorm.DB, ids []ID, opts TableOptions) ([]T, error) {
	idCol, err := idColumnWithOptions[T](tx, opts)
	if err != nil {
		return lo.Empty[[]T](), sderr.WithStack(err)
	}
	var r []T
	dbr := tx.Where(fmt.Sprintf("%s IN ?", idCol), ids).Find(&r)
	if dbr.Error != nil {
		return nil, dbr.Error
	}
	return r, nil
}

func DeleteById[T any, ID any](tx *gorm.DB, nullRow T, id ID, opts TableOptions) error {
	idCol, err := idColumnWithOptions[T](tx, opts)
	if err != nil {
		return sderr.WithStack(err)
	}
	dbr := tx.Where(fmt.Sprintf("%s=?", idCol), id).Delete(nullRow)
	if dbr.Error != nil {
		return sderr.WithStack(dbr.Error)
	}
	return nil
}

func DeleteByIds[T any, ID any](tx *gorm.DB, nullRow T, ids []ID, opts TableOptions) error {
	if len(ids) <= 0 {
		return nil
	}

	idCol, err := idColumnWithOptions[T](tx, opts)
	if err != nil {
		return sderr.WithStack(err)
	}
	dbr := tx.Where(fmt.Sprintf("%s IN ?", idCol), ids).Delete(nullRow)
	if dbr.Error != nil {
		return sderr.WithStack(dbr.Error)
	}
	return nil
}

func First[T any](tx *gorm.DB) (T, error) {
	var r T
	dbr := tx.First(&r)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
}

func Last[T any](tx *gorm.DB) (T, error) {
	var r T
	dbr := tx.Last(&r)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
}

func Take[T any](tx *gorm.DB) (T, error) {
	var r T
	dbr := tx.Take(&r)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
}

func Find[T any](tx *gorm.DB) (T, error) {
	var r T
	dbr := tx.Find(&r)
	if dbr.Error != nil {
		return r, dbr.Error
	}
	return r, nil
}

func getPrimaryColumnName(tx *gorm.DB, t any) (string, error) {
	cols, err := ParsePrimaryColumnNames(t, tx.NamingStrategy)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	if len(cols) <= 0 {
		return "", sderr.New("no primary key column")
	}
	if len(cols) > 1 {
		return "", sderr.New("multiple primary key column")
	}
	return cols[0], nil
}

func idColumnWithOptions[T any](tx *gorm.DB, opts TableOptions) (string, error) {
	if opts.IdColumn != "" {
		return opts.IdColumn, nil
	} else {
		var row T
		idCol, err := getPrimaryColumnName(tx, row)
		if err != nil {
			return "", sderr.WithStack(err)
		}
		return idCol, nil
	}
}
