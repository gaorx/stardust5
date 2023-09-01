package sdecho

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdmaps"
	"github.com/labstack/echo/v4"
)

const (
	ActionCall = "call"
	ActionShow = "show"
)

type AccessControl struct {
	Check             AccessControlChecker
	DefaultObjectVars map[string]string
}

type AccessControlChecker func(ctx context.Context, ec echo.Context, token Token, object Object, action string) (bool, error)

const (
	keyAccessControlChecker    = "sdecho.access_control_checker"
	keyAccessControlObjectVars = "sdecho.access_control_object_vars"
)

type accessControlChecker func(context.Context, echo.Context, Token, Object, string) error

func (ac AccessControl) Apply(app *echo.Echo) error {
	checker := func(ctx context.Context, ec echo.Context, token Token, object Object, action string) error {
		if object.IsPublic() {
			return nil
		}
		if token.UID == "" {
			return sderr.WithStack(ErrUnauthorized)
		}

		ok, err := ac.Check(ctx, ec, token, object, action)
		if err != nil {
			return sderr.WithStack(err)
		}
		if !ok {
			return sderr.WithStack(ErrForbidden)
		}
		return nil
	}

	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			ec.Set(keyAccessControlChecker, accessControlChecker(checker))
			ec.Set(keyAccessControlObjectVars, sdmaps.Ensure(ac.DefaultObjectVars))
			return next(ec)
		}
	}
	app.Use(middleware)
	return nil
}

func AccessControlCheck(ctx context.Context, ec echo.Context, token Token, object Object, action string) error {
	checker := MustGet[accessControlChecker](ec, keyAccessControlChecker)
	defaultObjectVars := MustGet[map[string]string](ec, keyAccessControlObjectVars)
	object1 := object.Expand(contextExpandMapper(ec), defaultObjectVars)
	return checker(ctx, ec, token, object1, action)
}
