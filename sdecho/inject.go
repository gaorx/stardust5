package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/labstack/echo/v4"
)

func Set(ec echo.Context, k string, v any) {
	ec.Set(k, v)
}

func Inject(k string, v any) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			Set(ec, k, v)
			return h(ec)
		}
	}
}

func Lookup[T any](ec echo.Context, k string) (T, bool) {
	v := ec.Get(k)
	if v == nil {
		var empty T
		return empty, false
	}
	typed, ok := v.(T)
	if !ok {
		var empty T
		return empty, false
	}
	return typed, true
}

func Get[S any](ec echo.Context, k string) S {
	state, ok := Lookup[S](ec, k)
	if !ok {
		panic(sderr.New("get in echo context error"))
	}
	return state
}
