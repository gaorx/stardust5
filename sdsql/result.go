package sdsql

import (
	"database/sql"
)

type Result struct {
	sql.Result
}

func ResultOf(sr sql.Result) Result {
	return Result{sr}
}

func (sr Result) RowsAffectedAndLastInsertId() (int64, int64, error) {
	affected, err := sr.Result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	lastInsertId, err := sr.Result.LastInsertId()
	if err != nil {
		return affected, 0, err
	}
	return affected, lastInsertId, nil
}

func (sr Result) RowsAffectedDef(def int64) int64 {
	n, err := sr.Result.RowsAffected()
	if err != nil {
		return def
	}
	return n
}

func (sr Result) LastInsertId(def int64) int64 {
	n, err := sr.Result.LastInsertId()
	if err != nil {
		return def
	}
	return n
}

func (sr Result) RowsAffectedAndLastInsertIdDef(affectedDef, lastInsertIdDef int64) (int64, int64) {
	affected, lastInsertId, err := sr.RowsAffectedAndLastInsertId()
	if err != nil {
		return affectedDef, lastInsertIdDef
	}
	return affected, lastInsertId
}
