package sdecho

import (
	"context"
	"fmt"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdreflect"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/gaorx/stardust5/sdurl"
	"github.com/gaorx/stardust5/sdvalidator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"log/slog"
	"net/http"
	"reflect"
	"slices"
)

// endpoint

type Endpoint struct {
	Methods     []string
	Path        string
	Object      string
	Func        any
	Middlewares []echo.MiddlewareFunc
	funcVal     reflect.Value
	inTypes     []reflect.Type
}

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

type RecordID interface{ ~string | ~int | ~int64 }

type Record[ID RecordID] interface {
	RecordID() ID
}

type Pagination struct {
	Size  int
	No    int
	Total int
}

type FindResult[T Record[ID], ID RecordID] struct {
	Data       []T
	Filter     any
	Pagination Pagination
}

type CrudAPI[T Record[ID], ID RecordID, F any] struct {
	Path        string
	Create      func(echo.Context, T) (T, error)
	Update      func(echo.Context, T, []string) (T, error)
	Delete      func(echo.Context, ID) error
	Get         func(echo.Context, ID) (T, error)
	Find        func(echo.Context, F, Pagination) (*FindResult[T, ID], error)
	Object      string
	ObjectR     string
	ObjectW     string
	Middlewares []echo.MiddlewareFunc
	PageSize    int
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

func (api CrudAPI[T, ID, F]) ToEndpoints() []Endpoint {
	selectObject := func(first, second string) string {
		if first != "" {
			return first
		}
		return second
	}

	var endpoints []Endpoint

	// create
	if api.Create != nil {
		endpoints = append(endpoints, API{
			Path:   sdurl.JoinPath(api.Path, "create"),
			Object: selectObject(api.ObjectW, api.Object),
			Func: func(ec echo.Context, req T) *Result {
				created, err := api.Create(ec, req)
				return Of(created, err)
			},
			Middlewares: api.Middlewares,
		}.ToEndpoint())
	}

	// update
	if api.Update != nil {
		endpoints = append(endpoints, API{
			Path:   sdurl.JoinPath(api.Path, "create"),
			Object: selectObject(api.ObjectW, api.Object),
			Func: func(ec echo.Context, req T) *Result {
				fields := sdstrings.SplitNonempty(ec.QueryParam("fields"), ",", true)
				updated, err := api.Update(ec, req, fields)
				return Of(updated, err)
			},
			Middlewares: api.Middlewares,
		}.ToEndpoint())
	}

	// delete
	if api.Delete != nil {
		endpoints = append(endpoints, API{
			Path:   sdurl.JoinPath(api.Path, "delete"),
			Object: selectObject(api.ObjectW, api.Object),
			Func: func(ec echo.Context, id ID) *Result {
				err := api.Delete(ec, id)
				return Of("deleted", err)
			},
			Middlewares: api.Middlewares,
		}.ToEndpoint())
	}

	// get
	if api.Get != nil {
		endpoints = append(endpoints, API{
			Path:   sdurl.JoinPath(api.Path, "get"),
			Object: selectObject(api.ObjectR, api.Object),
			Func: func(ec echo.Context, req struct {
				Id ID `json:"id"`
			}) *Result {
				record, err := api.Get(ec, req.Id)
				return Of(record, err)
			},
			Middlewares: api.Middlewares,
		}.ToEndpoint())
	}

	// find
	if api.Find != nil {
		pageSize := api.PageSize
		if pageSize <= 0 {
			pageSize = 20
		}
		endpoints = append(endpoints, API{
			Path:   sdurl.JoinPath(api.Path, "find"),
			Object: selectObject(api.ObjectR, api.Object),
			Func: func(ec Context, filter F) *Result {
				pg := Pagination{
					Size: ec.ArgInt("pageSize", pageSize),
					No:   ec.ArgInt("page", 1),
				}
				fr, err := api.Find(ec, filter, pg)
				return Of(fr.Data, err).WithFields(map[string]any{
					"page":      pg.No,
					"pageSize":  pg.Size,
					"pageTotal": pg.Total,
					"filter":    filter,
				})
			},
			Middlewares: api.Middlewares,
		}.ToEndpoint())
	}

	return endpoints
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
		return Err(err).Write(ec, routes.ResultOptions)
	}
	err = AccessControlCheck(context.Background(), ec, token, endpoint.expandObject(ec), ActionCall)
	if err != nil {
		return Err(err).Write(ec, routes.ResultOptions)
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
				return Err(sderr.Wrap(ErrBadRequest, err.Error())).Write(ec, routes.ResultOptions)
			}
			if err := sdvalidator.Struct(reqPtr); err != nil {
				if _, ok := sderr.AsT[validator.ValidationErrors](err); ok {
					return Err(sderr.Wrap(ErrBadRequest, "validate error")).Write(ec, routes.ResultOptions)
				} else {
					return Err(sderr.Wrap(ErrBadRequest, err.Error())).Write(ec, routes.ResultOptions)
				}
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
		return Err(sderr.WithStack(ErrInternalServerError)).Write(ec, routes.ResultOptions)
	}
	res := outVals[0].Interface().(*Result)
	if res == nil {
		res = Ok(nil)
	}
	return res.Write(ec, routes.ResultOptions)
}

func (endpoint *Endpoint) expandObject(ec echo.Context) string {
	return sdstrings.ExpandShellLike(endpoint.Object, func(k string) string {
		v := ec.QueryParam(k)
		if v == "" {
			v = ec.Param(k)
		}
		if v == "" {
			v0 := ec.Get(k)
			if v0 != nil {
				v = fmt.Sprintf("%v", v0)
			}
		}
		return v
	})
}
