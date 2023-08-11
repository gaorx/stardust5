package sdnotfounderr

import (
	"database/sql"
	"github.com/gaorx/stardust5/sdcache"
	"github.com/gaorx/stardust5/sderr"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var notFoundErrs = []error{
	gorm.ErrRecordNotFound,
	sdcache.ErrNotFound,
	sql.ErrNoRows,
	redis.Nil,
}

func Register(err error) {
	if err == nil {
		return
	}
	for _, err0 := range notFoundErrs {
		if sderr.Is(err, err0) {
			return
		}
	}
	notFoundErrs = append(notFoundErrs, err)
}

func Is(err error) bool {
	for _, notFoundErr := range notFoundErrs {
		if sderr.Is(err, notFoundErr) {
			return true
		}
	}
	return false
}
