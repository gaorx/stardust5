package sdecho

import (
	"github.com/labstack/echo/v4"
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
