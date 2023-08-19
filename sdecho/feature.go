package sdecho

import (
	"github.com/labstack/echo/v4"
)

type Feature interface {
	Apply(app *echo.Echo) error
}

type FeatureFunc func(*echo.Echo) error

func (f FeatureFunc) Apply(app *echo.Echo) error {
	return f(app)
}
