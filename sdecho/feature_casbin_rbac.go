package sdecho

import (
	"context"
	"github.com/gaorx/stardust5/sdcasbin"
	"github.com/labstack/echo/v4"
)

type CasbinRbac struct {
	Rbac sdcasbin.Rbac
}

func (ac CasbinRbac) Apply(app *echo.Echo) error {
	checker := func(_ context.Context, ec echo.Context, token Token, object, action string) (bool, error) {
		if ac.Rbac == nil {
			return false, nil
		}
		ok := ac.Rbac.IsGranted(token.UID, object, action)
		return ok, nil
	}
	return AccessControl{Check: checker}.Apply(app)
}
