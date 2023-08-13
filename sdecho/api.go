package sdecho

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

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

type CrudAPIs[ID ~string | ~int | ~int64, T any, F any] struct {
	Path    string
	Create  func(ec echo.Context, o *T) (*T, error)
	Update  func(ec echo.Context, o *T, fields []string) (*T, error)
	Delete  func(ec echo.Context, id ID) error
	GetById func(ec echo.Context, id ID) (*T, error)
	FindBy  func(ec echo.Context, filter *F) (FindResult[T, F], error)
	Object  string
	ObjectR string
	ObjectW string
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

func (crud CrudAPIs[ID, T, F]) ToEndpoints() []Endpoint {
	panic("TODO")
}
