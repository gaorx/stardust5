package sdcache

import (
	"github.com/gaorx/stardust5/sderr"
)

var (
	ErrNotFound = sderr.Sentinel("cache key not found")
)
