package sdecho

import (
	"time"

	"github.com/gaorx/stardust5/sdparse"
)

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
