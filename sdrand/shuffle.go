package sdrand

import (
	"github.com/samber/lo"
	"slices"
)

func Shuffle[T any](collection []T) {
	lo.Shuffle(collection)
}

func ShuffleClone[T any](collection []T) []T {
	return lo.Shuffle(slices.Clone(collection))
}
