package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdfile/sdfiletype"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	ViewJson = "JSON"
)

type ResultRenderer func(echo.Context, *Result, *RendererOptions) error

type Result struct {
	// http
	HttpStatus  int
	ContentType string
	Headers     map[string]string
	Raw         []byte

	// data
	Code     any
	Fields   map[string]any
	Data     any
	Error    error
	View     string
	Facade   any
	Renderer ResultRenderer
}

type RendererOptions struct {
	GetRenderer      func(view string) ResultRenderer
	CodeOk           any
	CodeBadRequest   any
	CodeTokenExpired any
	CodeUnauthorized any
	CodeForbidden    any
	CodeLogin        any
	CodeUnknown      any
}

var defaultRenderOptions = &RendererOptions{
	CodeOk:           "OK",
	CodeBadRequest:   "BadRequest",
	CodeTokenExpired: "TokenExpired",
	CodeUnauthorized: "Unauthorized",
	CodeForbidden:    "Forbidden",
	CodeLogin:        "Login",
	CodeUnknown:      "Unknown",
}

func (r *Result) Render(ec echo.Context, opts *RendererOptions) error {
	var r1 Result
	var opts1 RendererOptions
	r1 = *r
	if opts != nil {
		opts1 = *opts
	}

	selectCode := func(first, second any) any {
		if first != nil {
			return first
		} else {
			return second
		}
	}

	// facade
	// TODO

	// http status code
	if r1.HttpStatus <= 0 {
		r1.HttpStatus = http.StatusOK
	}

	// code
	if r1.Code == nil {
		if r1.Error == nil {
			r1.Code = selectCode(opts1.CodeOk, defaultRenderOptions.CodeOk)
		} else {
			if sderr.Is(r1.Error, ErrBadRequest) {
				r1.Code = selectCode(opts1.CodeBadRequest, defaultRenderOptions.CodeBadRequest)
			} else if sderr.Is(r1.Error, ErrTokenExpired) || sderr.Is(r1.Error, ErrDecodeToken) {
				r1.Code = selectCode(opts1.CodeTokenExpired, defaultRenderOptions.CodeTokenExpired)
			} else if sderr.Is(r1.Error, ErrUnauthorized) {
				r1.Code = selectCode(opts1.CodeUnauthorized, defaultRenderOptions.CodeUnauthorized)
			} else if sderr.Is(r1.Error, ErrForbidden) {
				r1.Code = selectCode(opts1.CodeForbidden, defaultRenderOptions.CodeForbidden)
			} else if sderr.Is(r1.Error, ErrLogin) {
				r1.Code = selectCode(opts1.CodeLogin, defaultRenderOptions.CodeLogin)
			} else {
				r1.Code = selectCode(opts1.CodeUnknown, defaultRenderOptions.CodeUnknown)
			}
		}
	}

	// view
	if r1.View == "" {
		r1.View = ViewJson
	}

	// renderer
	if r1.Renderer == nil {
		if r1.Raw != nil {
			if r1.ContentType == "" {
				r1.ContentType = sdfiletype.MatchMime(r1.Raw, echo.MIMEOctetStream)
			}
			r1.Renderer = renderRaw
		} else {
			if r1.View == ViewJson {
				r1.Renderer = renderJson
			} else {
				if opts1.GetRenderer == nil {
					panic(sderr.New("get result renderer error"))
				}
				r1.Renderer = opts1.GetRenderer(r1.View)
			}
		}
	}

	// render
	return r1.Renderer(ec, &r1, &opts1)
}

func Ok(data any) *Result {
	return &Result{Data: data, Error: nil}
}

func Err(err any) *Result {
	return &Result{Data: nil, Error: sderr.AsErr(err)}
}

func Of(data any, err any) *Result {
	if err != nil {
		return Err(err)
	} else {
		return Ok(data)
	}
}

func Raw(httpStatus int, contentType string, data []byte) *Result {
	return &Result{HttpStatus: httpStatus, ContentType: contentType, Raw: data}
}

func (r *Result) SetHttpStatus(status int) *Result {
	r.HttpStatus = status
	return r
}

func (r *Result) SetContentType(contentType string) *Result {
	r.ContentType = contentType
	return r
}

func (r *Result) WithHeader(k, v string) *Result {
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}
	r.Headers[k] = v
	return r
}

func (r *Result) WithHeaders(headers map[string]string) *Result {
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}
	for k, v := range headers {
		r.Headers[k] = v
	}
	return r
}

func (r *Result) SetCode(code any) *Result {
	r.Code = code
	return r
}

func (r *Result) WithField(k string, v any) *Result {
	if r.Fields == nil {
		r.Fields = map[string]any{}
	}
	r.Fields[k] = v
	return r
}

func (r *Result) WithFields(fields map[string]any) *Result {
	if r.Fields == nil {
		r.Fields = map[string]any{}
	}
	for k, v := range fields {
		r.Fields[k] = v
	}
	return r
}

func (r *Result) SetView(view string) *Result {
	r.View = view
	return r
}

func (r *Result) SetFacade(facade any) *Result {
	r.Facade = facade
	return r
}

func (r *Result) SetRenderer(renderer ResultRenderer) *Result {
	r.Renderer = renderer
	return r
}

func renderRaw(ec echo.Context, r1 *Result, _ *RendererOptions) error {
	for k, v := range r1.Headers {
		ec.Response().Header().Set(k, v)
	}
	return ec.Blob(r1.HttpStatus, r1.ContentType, r1.Raw)
}

func renderJson(ec echo.Context, r1 *Result, _ *RendererOptions) error {
	for k, v := range r1.Headers {
		ec.Response().Header().Set(k, v)
	}
	o := sdjson.Object{"code": r1.Code}
	if r1.Error == nil {
		o["data"] = r1.Data
	} else {
		o["error"] = r1.Error.Error()
	}
	for k, v := range r1.Fields {
		o[k] = v
	}
	return ec.JSON(r1.HttpStatus, o)
}
