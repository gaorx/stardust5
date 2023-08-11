package sdcasbin

import (
	"errors"
)

var (
	ErrIllegalUserId   = errors.New("illegal user id")
	ErrIllegalRoleId   = errors.New("illegal role id")
	ErrIllegalObjectId = errors.New("illegal object id")
	ErrIllegalAction   = errors.New("illegal action")
)
