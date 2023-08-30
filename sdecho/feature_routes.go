package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdreflect"
	"github.com/gaorx/stardust5/sdurl"
	"github.com/labstack/echo/v4"
	"slices"
)

type Routes struct {
	BasePath  string
	Endpoints []any
	*ResultOptions
}

const (
	keyRoutes = "sdecho.routes"
)

func (routes Routes) Apply(app *echo.Echo) error {
	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			ec.Set(keyRoutes, &routes)
			return next(ec)
		}
	}
	app.Use(middleware)

	// expand endpoints
	endpoints, err := routes.ExpandEndpoints()
	if err != nil {
		return sderr.WithStack(err)
	}

	// prepare
	for _, endpoint := range endpoints {
		if err := endpoint.prepare(); err != nil {
			return sderr.WrapWith(err, "prepare endpoint error", endpoint.Path)
		}
	}

	pathOf := func(p string) string {
		if routes.BasePath != "" {
			return sdurl.JoinPath(routes.BasePath, p)
		} else {
			return p
		}
	}

	// add routes
	for _, endpoint := range endpoints {
		endpoint1 := *endpoint
		h := func(ec echo.Context) error {
			return endpoint1.render(ec)
		}
		if slices.Contains(endpoint.Methods, "ANY") {
			app.Any(pathOf(endpoint.Path), h, endpoint.Middlewares...)
		} else {
			for _, method := range endpoint.Methods {
				app.Add(method, pathOf(endpoint.Path), h, endpoint.Middlewares...)
			}
		}
	}
	return nil
}

func (routes Routes) ExpandEndpoints() ([]*Endpoint, error) {
	var endpoints []*Endpoint

	var appendEndpoint func(anyEndpoint any) bool
	appendEndpoint = func(anyEndpoint any) bool {
		if anyEndpoint == nil {
			return true
		}
		switch x := anyEndpoint.(type) {
		case Endpoint:
			endpoints = append(endpoints, &x)
			return true
		case *Endpoint:
			if x != nil {
				x1 := *x
				endpoints = append(endpoints, &x1)
			}
			return true
		case []Endpoint:
			for _, elem := range x {
				appendEndpoint(elem)
			}
			return true
		case []*Endpoint:
			for _, elem := range x {
				appendEndpoint(elem)
			}
			return true
		case interface{ ToEndpoint() *Endpoint }:
			if x != nil {
				appendEndpoint(x.ToEndpoint())
			}
			return true
		case interface{ ToEndpoint() Endpoint }:
			if x != nil {
				appendEndpoint(x.ToEndpoint())
			}
			return true
		case interface{ ToEndpoints() []*Endpoint }:
			if x != nil {
				appendEndpoint(x.ToEndpoints())
			}
			return true
		case interface{ ToEndpoints() []Endpoint }:
			if x != nil {
				appendEndpoint(x.ToEndpoints())
			}
			return true
		default:
			return false
		}
	}

	// get endpoints
	for _, anyEndpoint := range routes.Endpoints {
		if ok := appendEndpoint(anyEndpoint); !ok {
			return nil, sderr.NewWith("illegal endpoint", sdreflect.TypeOf(anyEndpoint).String())
		}
	}
	return endpoints, nil
}
