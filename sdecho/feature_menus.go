package sdecho

import (
	"github.com/gaorx/stardust5/sdslices"
	"github.com/labstack/echo/v4"
)

type Menus []*Menu

const (
	keyMenus = "sdecho.menus"
)

func (menus Menus) Apply(app *echo.Echo) error {
	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			ec.Set(keyMenus, sdslices.Ensure(menus))
			return next(ec)
		}
	}
	app.Use(middleware)
	return nil
}
