package sdecho

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdreflect"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"log/slog"
	"reflect"
	"slices"
)

type Endpoint struct {
	Methods     []string
	Path        string
	Object      string
	Func        any
	Middlewares []echo.MiddlewareFunc
	funcVal     reflect.Value
	inTypes     []reflect.Type
}

type Routes struct {
	Endpoints []any
	*RendererOptions
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
		case interface{ AsEndpoints() []Endpoint }:
			if x != nil {
				appendEndpoint(x.AsEndpoints())
			}
			return true
		default:
			return false
		}
	}

	// get endpoints
	for _, anyEndpoint := range routes.Endpoints {
		if ok := appendEndpoint(anyEndpoint); !ok {
			return sderr.NewWith("illegal endpoint", sdreflect.TypeOf(anyEndpoint).String())
		}
	}

	// prepare
	for _, endpoint := range endpoints {
		if err := endpoint.prepare(); err != nil {
			return sderr.WrapWith(err, "prepare endpoint error", endpoint.Path)
		}
	}

	// add routes
	for _, endpoint := range endpoints {
		endpoint1 := *endpoint
		h := func(ec echo.Context) error {
			return endpoint1.render(ec)
		}
		if slices.Contains(endpoint.Methods, "ANY") {
			app.Any(endpoint.Path, h, endpoint.Middlewares...)
		} else {
			for _, method := range endpoint.Methods {
				app.Add(method, endpoint.Path, h, endpoint.Middlewares...)
			}
		}
	}
	return nil
}

func (endpoint *Endpoint) prepare() error {
	if endpoint.Path == "" {
		return sderr.New("no path in endpoint")
	}
	if len(endpoint.Methods) <= 0 {
		return sderr.NewWith("no methods in endpoint", endpoint.Path)
	}
	if slices.Contains(endpoint.Methods, "*") {
		endpoint.Methods = []string{"ANY"}
	}
	if endpoint.Func == nil {
		return sderr.NewWith("no func in endpoint", endpoint.Path)
	}
	funcVal := sdreflect.ValueOf(endpoint.Func)
	inTypes, outTypes := sdreflect.InOutTypes(funcVal.Type())
	if len(outTypes) != 1 || outTypes[0] != sdreflect.T[*Result]() {
		return sderr.NewWith("illegal result type", endpoint.Path)
	}
	numFreeParam := 0
	for _, inType := range inTypes {
		if !slices.Contains([]reflect.Type{
			sdreflect.T[echo.Context](),
			sdreflect.T[Context](),
			sdreflect.T[Token](),
			sdreflect.T[*Token](),
		}, inType) {
			numFreeParam += 1
		}
	}
	if numFreeParam > 1 {
		return sderr.NewWith("illegal argument type", endpoint.Path)
	}
	endpoint.funcVal, endpoint.inTypes = funcVal, inTypes
	return nil
}

func (endpoint *Endpoint) render(ec echo.Context) error {
	routes := MustGet[*Routes](ec, keyRoutes)
	token, err := TokenDecode(context.Background(), ec)
	if err != nil {
		return Err(err).Render(ec, routes.RendererOptions)
	}
	err = AccessControlCheck(context.Background(), ec, token, endpoint.Object, ActionCall)
	if err != nil {
		return Err(err).Render(ec, routes.RendererOptions)
	}
	var inVals, outVals []reflect.Value
	for _, inTyp := range endpoint.inTypes {
		switch inTyp {
		case sdreflect.T[echo.Context]():
			inVals = append(inVals, reflect.ValueOf(ec))
		case sdreflect.T[Context]():
			inVals = append(inVals, reflect.ValueOf(C(ec)))
		case sdreflect.T[Token]():
			inVals = append(inVals, reflect.ValueOf(token))
		case sdreflect.T[*Token]():
			inVals = append(inVals, reflect.ValueOf(&token))
		default:
			reqIsPtr := inTyp.Kind() == reflect.Ptr
			var reqPtr any
			if reqIsPtr {
				reqPtr = reflect.New(inTyp.Elem()).Interface()
			} else {
				reqPtr = reflect.New(inTyp).Interface()
			}
			if err := ec.Bind(reqPtr); err != nil {
				return Err(sderr.Wrap(ErrBadRequest, err.Error())).Render(ec, routes.RendererOptions)
			}
			if reqIsPtr {
				inVals = append(inVals, reflect.ValueOf(reqPtr))
			} else {
				inVals = append(inVals, reflect.ValueOf(reqPtr).Elem())
			}
		}
	}
	if ok := lo.Try0(func() {
		outVals = endpoint.funcVal.Call(inVals)
	}); !ok {
		slog.With("path", endpoint.Path).Error("call endpoint error")
		return Err(sderr.WithStack(ErrInternalServerError)).Render(ec, routes.RendererOptions)
	}
	res := outVals[0].Interface().(*Result)
	if res == nil {
		res = Ok(nil)
	}
	return res.Render(ec, routes.RendererOptions)
}
