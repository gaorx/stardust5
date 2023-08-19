package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/labstack/echo/v4"
)

type Echo struct {
	*echo.Echo
}

func E(e *echo.Echo) Echo {
	return Echo{e}
}

func (e Echo) Install(features ...Feature) error {
	for _, feature := range features {
		if feature != nil {
			if err := feature.Apply(e.Echo); err != nil {
				return sderr.WithStack(err)
			}
		}
	}
	return nil
}
