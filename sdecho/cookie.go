package sdecho

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdparse"
)

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
