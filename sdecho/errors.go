package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
)

var (
	ErrUnauthorized        = sderr.Sentinel("unauthorized")
	ErrForbidden           = sderr.Sentinel("forbidden")
	ErrBadRequest          = sderr.Sentinel("bad request")
	ErrInternalServerError = sderr.Sentinel("internal server error")
	ErrDecodeToken         = sderr.Sentinel("decode token error")
	ErrTokenExpired        = sderr.Sentinel("token expired")
	ErrLogin               = sderr.Sentinel("login error")
)
