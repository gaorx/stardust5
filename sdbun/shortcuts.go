package sdbun

import (
	"context"
	"database/sql"
	"github.com/gaorx/stardust5/sdreflect"
	"github.com/gaorx/stardust5/sdsql"
	"github.com/uptrace/bun"
	"reflect"
)

func Transaction(ctx context.Context, db bun.IDB, action func(context.Context, bun.Tx) error, opts *sql.TxOptions) error {
	return db.RunInTx(ctx, opts, action)
}

func TransactionFor[R any](ctx context.Context, db bun.IDB, action func(context.Context, bun.Tx) (R, error), opts *sql.TxOptions) (R, error) {
	var r R
	err := db.RunInTx(ctx, opts, func(ctx context.Context, tx bun.Tx) error {
		r0, err := action(ctx, tx)
		if err != nil {
			return err
		}
		r = r0
		return nil
	})
	return r, err
}

func Insert(ctx context.Context, db bun.IDB, v any, qfn func(query *bun.InsertQuery) *bun.InsertQuery) (sdsql.Result, error) {
	sr, err := db.NewInsert().Model(ptrOfStruct(v)).Apply(qfn).Exec(ctx)
	if err != nil {
		return sdsql.Result{}, err
	}
	return sdsql.ResultOf(sr), nil
}

func Update(ctx context.Context, db bun.IDB, v any, qfn func(query *bun.UpdateQuery) *bun.UpdateQuery) (sdsql.Result, error) {
	sr, err := db.NewUpdate().Model(ptrOfStruct(v)).Apply(qfn).Exec(ctx)
	if err != nil {
		return sdsql.Result{}, err
	}
	return sdsql.ResultOf(sr), nil
}

func Delete[ROW any](ctx context.Context, db bun.IDB, qfn func(query *bun.DeleteQuery) *bun.DeleteQuery) (sdsql.Result, error) {
	sr, err := db.NewDelete().Apply(qfn).Apply(modelApplier[*bun.DeleteQuery, ROW]()).Exec(ctx)
	if err != nil {
		return sdsql.Result{}, err
	}
	return sdsql.ResultOf(sr), nil
}

func SelectMany[ROW any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery) ([]ROW, error) {
	var r []ROW
	err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).Scan(ctx, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func SelectFirst[ROW any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery) (ROW, error) {
	t := sdreflect.T[ROW]()
	if isPtrToStruct(t) {
		dest := reflect.New(t.Elem()).Interface()
		err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).Scan(ctx, dest)
		if err != nil {
			var zero ROW
			return zero, err
		}
		return dest.(ROW), nil
	} else {
		var r ROW
		err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).Scan(ctx, &r)
		if err != nil {
			var zero ROW
			return zero, err
		}
		return r, nil
	}
}

func SelectManyRaw[ROW any](ctx context.Context, db bun.IDB, q string, args []any) ([]ROW, error) {
	var r []ROW
	err := db.NewRaw(q, args...).Scan(ctx, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func SelectFirstRaw[ROW any](ctx context.Context, db bun.IDB, q string, args []any) (ROW, error) {
	t := sdreflect.T[ROW]()
	if isPtrToStruct(t) {
		dest := reflect.New(t.Elem()).Interface()
		err := db.NewRaw(q, args...).Scan(ctx, dest)
		if err != nil {
			var zero ROW
			return zero, err
		}
		return dest.(ROW), nil
	} else {
		var r ROW
		err := db.NewRaw(q, args...).Scan(ctx, &r)
		if err != nil {
			var zero ROW
			return zero, err
		}
		return r, nil
	}
}

func SelectOne[T any](ctx context.Context, db bun.IDB, q string, args []any) (T, error) {
	return SelectFirstRaw[T](ctx, db, q, args)
}

func Count[ROW any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery) (int64, error) {
	n, err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).Count(ctx)
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}

func Exists[ROW any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery) (bool, error) {
	exists, err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func SelectManyAndCount[ROW any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery) ([]ROW, int64, error) {
	var r []ROW
	n, err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).ScanAndCount(ctx, &r)
	if err != nil {
		var zero []ROW
		return zero, 0, err
	}
	return r, int64(n), nil
}

type modelQuery[Q any] interface {
	GetTableName() string
	Model(model any) Q
}

func modelApplier[Q modelQuery[Q], ROW any]() func(Q) Q {
	return func(q Q) Q {
		if q.GetTableName() == "" {
			if model := modelOfTyped[ROW](); model != nil {
				q = q.Model(model)
			}
		}
		return q
	}
}

func modelOf(model any) any {
	if model == nil {
		return nil
	}
	t := reflect.TypeOf(model)
	k := t.Kind()
	if k == reflect.Struct {
		return reflect.Zero(reflect.PtrTo(t)).Interface()
	} else if k == reflect.Pointer {
		if t.Elem().Kind() == reflect.Struct {
			return reflect.Zero(t).Interface()
		} else {
			var getElem func(reflect.Type) reflect.Type
			getElem = func(t1 reflect.Type) reflect.Type {
				if t1.Kind() == reflect.Pointer {
					return getElem(t1.Elem())
				} else {
					return t1
				}
			}
			base := getElem(t)
			if base.Kind() == reflect.Struct {
				return reflect.Zero(reflect.PtrTo(base)).Interface()
			} else {
				return nil
			}
		}
	} else {
		return nil
	}
}

func modelOfTyped[T any]() any {
	var model T
	return modelOf(model)
}

func isPtrToStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct
}

func ptrOfStruct(v any) any {
	if v == nil {
		return nil
	}
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Struct {
		if vv.CanAddr() {
			return vv.Addr().Interface()
		} else {
			p := reflect.New(vv.Type())
			p.Elem().Set(vv)
			return p.Interface()
		}
	} else {
		return v
	}
}
