package sdecho

import (
	"errors"
)

var (
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrBadRequest          = errors.New("bad request")
	ErrInternalServerError = errors.New("internal server error")
	ErrDecodeToken         = errors.New("decode token error")
	ErrTokenExpired        = errors.New("token expired")
	ErrLogin               = errors.New("login error")
)
