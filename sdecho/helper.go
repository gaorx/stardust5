package sdecho

import (
	"fmt"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/labstack/echo/v4"
)

func contextExpandMapper(ec echo.Context) sdstrings.ExpandMapper {
	return func(k string) string {
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
	}
}
