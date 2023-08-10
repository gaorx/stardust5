package sdparse

import (
	"time"

	"github.com/samber/lo"
)

func MustInt64(s string) int64 {
	return lo.Must(Int64(s))
}

func MustInt(s string) int {
	return lo.Must(Int(s))
}

func MustUint64(s string) uint64 {
	return lo.Must(Uint64(s))
}

func MustUint(s string) uint {
	return lo.Must(Uint(s))
}

func MustFloat64(s string) float64 {
	return lo.Must(Float64(s))
}

func MustBool(s string) bool {
	return lo.Must(Bool(s))
}

func MustTime(s string) time.Time {
	return lo.Must(Time(s))
}
