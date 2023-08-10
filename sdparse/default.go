package sdparse

import (
	"time"

	"github.com/samber/mo"
)

func Int64Def(s string, def int64) int64 {
	return mo.TupleToResult(Int64(s)).OrElse(def)
}

func IntDef(s string, def int) int {
	return mo.TupleToResult(Int(s)).OrElse(def)
}

func Uint64Def(s string, def uint64) uint64 {
	return mo.TupleToResult(Uint64(s)).OrElse(def)
}

func UintDef(s string, def uint) uint {
	return mo.TupleToResult(Uint(s)).OrElse(def)
}

func Float64Def(s string, def float64) float64 {
	return mo.TupleToResult(Float64(s)).OrElse(def)
}

func BoolDef(s string, def bool) bool {
	return mo.TupleToResult(Bool(s)).OrElse(def)
}

func TimeDef(s string, def time.Time) time.Time {
	return mo.TupleToResult(Time(s)).OrElse(def)
}
