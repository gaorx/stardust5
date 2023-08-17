package sdecho

import (
	"context"
	"github.com/gaorx/stardust5/sdcasbin"
	"github.com/labstack/echo/v4"
)

type CasbinRbac struct {
	Rbac              sdcasbin.Rbac
	DefaultObjectVars map[string]string
}

func (ac CasbinRbac) Apply(app *echo.Echo) error {
	checker := func(_ context.Context, ec echo.Context, token Token, object Object, action string) (bool, error) {
		if ac.Rbac == nil {
			return false, nil
		}
		ok := ac.Rbac.IsGranted(token.UID, object.String(), action)
		return ok, nil
	}
	return AccessControl{
		Check:             checker,
		DefaultObjectVars: ac.DefaultObjectVars,
	}.Apply(app)
}
