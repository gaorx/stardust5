package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdfile/sdfiletype"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdslices"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	rkRaw  = "RAW"
	rkJson = "JSON"
	rkHtml = "HTML"
)

type Result struct {
	// kind
	kind string

	// http
	HttpStatus  int
	ContentType string
	Headers     map[string]string

	// data
	Facade any

	// API & HTML
	Code   any
	Fields map[string]any
	Data   any
	Error  error
	View   string
}

type ResultOptions struct {
	CodeOk           any
	CodeBadRequest   any
	CodeTokenExpired any
	CodeUnauthorized any
	CodeForbidden    any
	CodeLogin        any
	CodeUnknown      any
}

var defaultResultOptions = &ResultOptions{
	CodeOk:           "OK",
	CodeBadRequest:   "BadRequest",
	CodeTokenExpired: "TokenExpired",
	CodeUnauthorized: "Unauthorized",
	CodeForbidden:    "Forbidden",
	CodeLogin:        "Login",
	CodeUnknown:      "Unknown",
}

func (r *Result) Write(ec echo.Context, opts *ResultOptions) error {
	var r1 Result
	var opts1 ResultOptions
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
			r1.Code = selectCode(opts1.CodeOk, defaultResultOptions.CodeOk)
		} else {
			if sderr.Is(r1.Error, ErrBadRequest) {
				r1.Code = selectCode(opts1.CodeBadRequest, defaultResultOptions.CodeBadRequest)
			} else if sderr.Is(r1.Error, ErrTokenExpired) || sderr.Is(r1.Error, ErrDecodeToken) {
				r1.Code = selectCode(opts1.CodeTokenExpired, defaultResultOptions.CodeTokenExpired)
			} else if sderr.Is(r1.Error, ErrUnauthorized) {
				r1.Code = selectCode(opts1.CodeUnauthorized, defaultResultOptions.CodeUnauthorized)
			} else if sderr.Is(r1.Error, ErrForbidden) {
				r1.Code = selectCode(opts1.CodeForbidden, defaultResultOptions.CodeForbidden)
			} else if sderr.Is(r1.Error, ErrLogin) {
				r1.Code = selectCode(opts1.CodeLogin, defaultResultOptions.CodeLogin)
			} else {
				r1.Code = selectCode(opts1.CodeUnknown, defaultResultOptions.CodeUnknown)
			}
		}
	}

	// write
	switch r1.kind {
	case rkRaw:
		return r1.writeRaw(ec)
	case rkJson:
		return r1.writeJson(ec)
	case rkHtml:
		return r1.writeHtml(ec)
	default:
		panic(sderr.NewWith("illegal result format", r1.kind))
	}
}

func Raw(httpStatus int, contentType string, data []byte) *Result {
	return &Result{
		kind:        rkRaw,
		HttpStatus:  httpStatus,
		ContentType: contentType,
		Data:        sdslices.Ensure(data),
		Error:       nil,
	}
}

func Ok(data any) *Result {
	return &Result{kind: rkJson, Data: data, Error: nil}
}

func Err(err any) *Result {
	return &Result{kind: rkJson, Data: nil, Error: sderr.AsErr(err)}
}

func Of(data any, err any) *Result {
	if err != nil {
		return Err(err)
	} else {
		return Ok(data)
	}
}

func PageOf(data any, view string) *Result {
	return &Result{kind: rkHtml, Data: data, Error: nil, View: view}
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

func (r *Result) SetFacade(facade any) *Result {
	r.Facade = facade
	return r
}

func (r *Result) writeRaw(ec echo.Context) error {
	raw := r.Data.([]byte)
	contentType := r.ContentType
	if contentType == "" {
		contentType = sdfiletype.MatchMime(raw, echo.MIMEOctetStream)
	}
	for k, v := range r.Headers {
		ec.Response().Header().Set(k, v)
	}
	return ec.Blob(r.HttpStatus, contentType, raw)
}

func (r *Result) writeJson(ec echo.Context) error {
	for k, v := range r.Headers {
		ec.Response().Header().Set(k, v)
	}
	return ec.JSON(r.HttpStatus, r.makeResponse())
}

func (r *Result) writeHtml(ec echo.Context) error {
	for k, v := range r.Headers {
		ec.Response().Header().Set(k, v)
	}
	return ec.Render(r.HttpStatus, r.View, r.makeResponse())
}

func (r *Result) makeResponse() sdjson.Object {
	o := sdjson.Object{"code": r.Code}
	if r.Error == nil {
		o["data"] = r.Data
	} else {
		o["error"] = r.Error.Error()
	}
	for k, v := range r.Fields {
		o[k] = v
	}
	return o
}
