package sdecho

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdreflect"
	"github.com/gaorx/stardust5/sdslog"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/gaorx/stardust5/sdurl"
	"github.com/gaorx/stardust5/sdvalidator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
	"reflect"
	"slices"
)

// endpoint

type Endpoint struct {
	Methods     []string
	Path        string
	Object      Object
	Func        any
	Middlewares []echo.MiddlewareFunc
	handler     echo.HandlerFunc
}

type Page struct {
	Method      string
	Path        string
	Object      Object
	Func        any
	Middlewares []echo.MiddlewareFunc
}

type API struct {
	Path        string
	Object      Object
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
	Object      Object
	ObjectR     Object
	ObjectW     Object
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
	selectObject := func(first, second Object) Object {
		if !first.IsEmpty() {
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
				return ResultOf(created, err)
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
				return ResultOf(updated, err)
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
				return ResultOf("deleted", err)
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
				return ResultOf(record, err)
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
				return ResultOf(fr.Data, err).WithFields(map[string]any{
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
	// path
	if endpoint.Path == "" {
		return sderr.New("no path in endpoint")
	}

	// method
	if len(endpoint.Methods) <= 0 {
		return sderr.NewWith("no methods in endpoint", endpoint.Path)
	}
	if slices.Contains(endpoint.Methods, "*") {
		endpoint.Methods = []string{"ANY"}
	}

	h, err := endpoint.ToHandler()
	if err != nil {
		return sderr.WithStack(err)
	}
	endpoint.handler = h
	return nil
}

func (endpoint *Endpoint) ToHandler() (echo.HandlerFunc, error) {
	if endpoint.Func == nil {
		return nil, sderr.New("nil func")
	} else if h, ok := endpoint.Func.(echo.HandlerFunc); ok {
		return h, nil
	} else if h, ok := endpoint.Func.(func(ec echo.Context) error); ok {
		return h, nil
	} else if h, ok := endpoint.Func.(func(ec Context) error); ok {
		return func(ec echo.Context) error {
			return h(C(ec))
		}, nil
	} else {
		funcVal := sdreflect.ValueOf(endpoint.Func)
		inTypes, outTypes := sdreflect.InOutTypes(funcVal.Type())
		if len(outTypes) != 1 || outTypes[0] != sdreflect.T[*Result]() {
			return nil, sderr.NewWith("illegal result type", endpoint.Path)
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
			return nil, sderr.NewWith("illegal argument type", endpoint.Path)
		}
		return func(ec echo.Context) error {
			return endpoint.renderDefault(ec, funcVal, inTypes)
		}, nil
	}
}

func (endpoint *Endpoint) renderDefault(ec echo.Context, funcVal reflect.Value, inTypes []reflect.Type) error {
	routes := MustGet[*Routes](ec, keyRoutes)

	token, err := TokenDecode(context.Background(), ec)
	if err != nil {
		return ResultErr(err).Write(ec, routes.ResultOptions)
	}
	err = AccessControlCheck(context.Background(), ec, token, endpoint.Object, ActionCall)
	if err != nil {
		return ResultErr(err).Write(ec, routes.ResultOptions)
	}
	var inVals, outVals []reflect.Value
	for _, inTyp := range inTypes {
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
				return ResultErr(sderr.Wrap(ErrBadRequest, err.Error())).Write(ec, routes.ResultOptions)
			}
			if isStructOrStructPtr(inTyp) {
				if err := sdvalidator.Struct(reqPtr); err != nil {
					if _, ok := sderr.AsT[validator.ValidationErrors](err); ok {
						return ResultErr(sderr.Wrap(ErrBadRequest, "validate error")).Write(ec, routes.ResultOptions)
					} else {
						return ResultErr(sderr.Wrap(ErrBadRequest, err.Error())).Write(ec, routes.ResultOptions)
					}
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
		outVals = funcVal.Call(inVals)
	}); !ok {
		sdslog.WithAttr("path", endpoint.Path).Error("call endpoint error")
		return ResultErr(sderr.WithStack(ErrInternalServerError)).Write(ec, routes.ResultOptions)
	}
	res := outVals[0].Interface().(*Result)
	if res == nil {
		res = ResultOk(nil)
	}
	return res.Write(ec, routes.ResultOptions)
}

func isStructOrStructPtr(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr:
		return isStructOrStructPtr(typ.Elem())
	default:
		return false
	}
}
