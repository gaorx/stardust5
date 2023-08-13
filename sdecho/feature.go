package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/labstack/echo/v4"
)

type Feature interface {
	Apply(app *echo.Echo) error
}

func Install(app *echo.Echo, features ...Feature) error {
	for _, feature := range features {
		if feature != nil {
			if err := feature.Apply(app); err != nil {
				return sderr.WithStack(err)
			}
		}
	}
	return nil
}

type FeatureFunc func(*echo.Echo) error

func (f FeatureFunc) Apply(app *echo.Echo) error {
	return f(app)
}
