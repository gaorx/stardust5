package sdecho

import (
	"github.com/labstack/echo/v4"
)

func CrossOrigin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			ec.Response().Header().Add("Access-Control-Allow-Origin", "*")
			ec.Response().Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			ec.Response().Header().Add("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
			return next(ec)
		}
	}
}
