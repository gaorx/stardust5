package sdcasbin

import (
	"github.com/gaorx/stardust5/sderr"
)

var (
	ErrIllegalUserId   = sderr.Sentinel("illegal user id")
	ErrIllegalRoleId   = sderr.Sentinel("illegal role id")
	ErrIllegalObjectId = sderr.Sentinel("illegal object id")
	ErrIllegalAction   = sderr.Sentinel("illegal action")
)
