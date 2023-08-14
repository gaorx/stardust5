package sdecho

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Page struct {
	Method      string
	Path        string
	Object      string
	Func        any
	Middlewares []echo.MiddlewareFunc
}

type API struct {
	Path        string
	Object      string
	Func        any
	Middlewares []echo.MiddlewareFunc
}

type FindResult[T any, F any] struct {
	Data      []*T
	Filter    *F
	PageSize  int
	PageNo    int
	PageTotal int
}

func (p Page) ToEndpoint() Endpoint {
	if p.Method == "" {
		p.Method = http.MethodGet
	}
	return Endpoint{
		Methods:     []string{p.Method},
		Path:        p.Path,
		Object:      p.Object,
		Func:        p.Func,
		Middlewares: p.Middlewares,
	}
}

func (api API) ToEndpoint() Endpoint {
	return Endpoint{
		Methods:     []string{http.MethodPost},
		Path:        api.Path,
		Object:      api.Object,
		Func:        api.Func,
		Middlewares: api.Middlewares,
	}
}
