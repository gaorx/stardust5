package sdecho

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdparse"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Context struct {
	echo.Context
}

func C(c echo.Context) Context {
	if c1, ok := c.(Context); ok {
		return c1
	} else {
		return Context{c}
	}
}

// context data

func Get[T any](ec echo.Context, k string) (T, bool) {
	v := ec.Get(k)
	if v == nil {
		return lo.Empty[T](), false
	}
	typed, ok := v.(T)
	if !ok {
		return lo.Empty[T](), false
	}
	return typed, true
}

func MustGet[T any](ec echo.Context, k string) T {
	v, ok := Get[T](ec, k)
	if !ok {
		panic(sderr.NewWith("not found in context", k))
	}
	return v
}

// request

func (c Context) RequestBodyBytes() ([]byte, error) {
	reader := c.Request().Body
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, sderr.Wrap(err, "sdecho read request body error")
	}
	return r, nil
}

func (c Context) RequestBodyString() (string, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c Context) RequestBodyAsJsonValue() (sdjson.Value, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return sdjson.Value{}, err
	}
	v, err := sdjson.UnmarshalValue(b)
	if err != nil {
		return sdjson.Value{}, err
	}
	return v, nil
}

func (c Context) RequestBodyAsJsonObject() (sdjson.Object, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c Context) RequestBodyAsJsonArray() (sdjson.Array, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return nil, err
	}
	var a []any
	err = json.Unmarshal(b, &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (c Context) RequestBodyAs(v any) error {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// request arguments

func (c Context) ArgString(name, def string) string {
	v := c.QueryParam(name)
	if v != "" {
		return v
	}
	v = c.Param(name)
	if v != "" {
		return v
	}
	return def
}

func (c Context) ArgInt(name string, def int) int {
	return sdparse.IntDef(c.ArgString(name, ""), def)
}

func (c Context) ArgInt64(name string, def int64) int64 {
	return sdparse.Int64Def(c.ArgString(name, ""), def)
}

func (c Context) ArgFloat64(name string, def float64) float64 {
	return sdparse.Float64Def(c.ArgString(name, ""), def)
}

func (c Context) ArgBool(name string, def bool) bool {
	return sdparse.BoolDef(c.ArgString(name, ""), def)
}

func (c Context) ArgTime(name string, def time.Time) time.Time {
	return sdparse.TimeDef(c.ArgString(name, ""), def)
}

func (c Context) ArgStringFirst(names []string, def string) string {
	for _, name := range names {
		if name == "" {
			continue
		}
		if arg := c.ArgString(name, ""); arg != "" {
			return arg
		}
	}
	return def
}

func (c Context) ArgIntFirst(names []string, def int) int {
	return sdparse.IntDef(c.ArgStringFirst(names, ""), def)
}

func (c Context) ArgInt64First(names []string, def int64) int64 {
	return sdparse.Int64Def(c.ArgStringFirst(names, ""), def)
}

func (c Context) ArgFloat64First(names []string, def float64) float64 {
	return sdparse.Float64Def(c.ArgStringFirst(names, ""), def)
}

func (c Context) ArgBoolFirst(names []string, def bool) bool {
	return sdparse.BoolDef(c.ArgStringFirst(names, ""), def)
}

func (c Context) ArgTimeFirst(names []string, def time.Time) time.Time {
	return sdparse.TimeDef(c.ArgStringFirst(names, ""), def)
}

// cookie

func (c Context) CookieString(name, def string) string {
	v, err := c.Cookie(name)
	if err != nil {
		return def
	}
	return v.Value
}

func (c Context) SetCookieString(name, val string, path string, maxAge int) {
	c.SetCookie(&http.Cookie{
		Name:   name,
		Value:  val,
		Path:   path,
		MaxAge: maxAge,
	})
}

func (c Context) DeleteCookie(name string, path string) {
	c.SetCookieString(name, "", path, -1)
}

func (c Context) CookieInt(name string, def int) int {
	return sdparse.IntDef(c.CookieString(name, ""), def)
}

func (c Context) SetCookieInt(name string, val int, path string, maxAge int) {
	c.SetCookieString(name, strconv.Itoa(val), path, maxAge)
}

func (c Context) CookieInt64(name string, def int64) int64 {
	return sdparse.Int64Def(c.CookieString(name, ""), def)
}

func (c Context) SetCookieInt64(name string, val int64, path string, maxAge int) {
	c.SetCookieString(name, strconv.FormatInt(val, 10), path, maxAge)
}

func (c Context) CookieJson(name string, v any) error {
	base64Str := c.CookieString(name, "")
	jsonBytes, err := base64.URLEncoding.DecodeString(base64Str)
	if err != nil {
		return err
	}
	err = sdjson.Unmarshal(jsonBytes, v)
	if err != nil {
		return sderr.Wrap(err, "sdecho unmarshal cookie json error")
	}
	return nil
}

func (c Context) SetCookieJson(name string, v any, path string, maxAge int) error {
	jsonBytes, err := sdjson.Marshal(v)
	if err != nil {
		return sderr.Wrap(err, "sdecho marshal json cookie error")
	}
	c.SetCookieString(name, base64.URLEncoding.EncodeToString(jsonBytes), path, maxAge)
	return nil
}

func (c Context) CookieJsonObject(name string, def sdjson.Object) sdjson.Object {
	base64Str := c.CookieString(name, "")
	jsonBytes, err := base64.URLEncoding.DecodeString(base64Str)
	if err != nil {
		return def
	}
	v, err := sdjson.UnmarshalValue(jsonBytes)
	if err != nil {
		return def
	}
	return v.AsObjectDef(def)
}
