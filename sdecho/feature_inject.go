package sdecho

import (
	"github.com/labstack/echo/v4"
	"maps"
)

type Inject map[string]any

func (inject Inject) Apply(app *echo.Echo) error {
	if len(inject) <= 0 {
		return nil
	}
	data := maps.Clone(inject)
	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			for k, v := range data {
				ec.Set(k, v)
			}
			return next(ec)
		}
	}
	app.Use(middleware)
	return nil
}
