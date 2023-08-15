package sdecho

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"slices"
)

const (
	ObjectPublic = "public"
	ActionCall   = "call"
	ActionShow   = "show"
)

type AccessControl struct {
	Check func(ctx context.Context, ec echo.Context, token Token, object, action string) (bool, error)
}

const (
	keyAccessControl = "sdecho.access_control"
)

type accessControlChecker func(context.Context, echo.Context, Token, string, string) error

func (ac AccessControl) Apply(app *echo.Echo) error {
	checker := func(ctx context.Context, ec echo.Context, token Token, object, action string) error {
		if object == ObjectPublic {
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
			ec.Set(keyAccessControl, accessControlChecker(checker))
			return next(ec)
		}
	}
	app.Use(middleware)
	return nil
}

func AccessControlCheck(ctx context.Context, ec echo.Context, token Token, object, action string) error {
	checker := MustGet[accessControlChecker](ec, keyAccessControl)
	return checker(ctx, ec, token, object, action)
}

func ObjectExpand(object string, dataOrMapper any) string {
	if dataOrMapper == nil {
		return object
	}
	switch x := dataOrMapper.(type) {
	case func(string) string:
		if x == nil {
			return object
		}
		return sdstrings.ExpandShellLike(object, x)
	case map[string]string:
		return sdstrings.ExpandShellLikeMap(object, x)
	default:
		panic(sderr.New("illegal dataOrMapper"))
	}
}

func ObjectsExpand(objects []string, dataOrMapper any) []string {
	if dataOrMapper == nil {
		return slices.Clone(objects)
	}
	switch x := dataOrMapper.(type) {
	case func(string) string:
		if x == nil {
			return slices.Clone(objects)
		}
		return lo.Map(objects, func(object string, _ int) string {
			return sdstrings.ExpandShellLike(object, x)
		})
	case map[string]string:
		return lo.Map(objects, func(object string, _ int) string {
			return sdstrings.ExpandShellLikeMap(object, x)
		})
	default:
		panic(sderr.New("illegal dataOrMapper"))
	}

}
