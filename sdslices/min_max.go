package sdslices

import (
	"cmp"
	"slices"
)

func MinV[T cmp.Ordered](some ...T) T {
	return slices.Min(some)
}

func MaxV[T cmp.Ordered](some ...T) T {
	return slices.Max(some)
}
