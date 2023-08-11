package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdslog"
	"github.com/samber/lo"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Options struct {
	DebugMode    bool
	LogSkipper   middleware.Skipper
	ErrorHandler echo.HTTPErrorHandler
}

func New(opts *Options) *echo.Echo {
	opts1 := lo.FromPtr(opts)
	if opts1.ErrorHandler == nil {
		opts1.ErrorHandler = defaultHttpErrorHandler
	}
	app := echo.New()
	app.Debug = opts1.DebugMode
	app.HideBanner = true
	app.HidePort = true
	app.Use(LoggingRecover(opts1.LogSkipper))
	app.HTTPErrorHandler = opts1.ErrorHandler
	return app
}

func defaultHttpErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	he, ok := sderr.AsT[*echo.HTTPError](err)
	if ok {
		if he.Internal != nil {
			if herr, ok := sderr.AsT[*echo.HTTPError](he.Internal); ok {
				he = herr
			}
		}
	} else {
		errMsg := http.StatusText(http.StatusInternalServerError)
		if c.QueryParam("_show_error") == "1" {
			errMsg = errMsg + strings.Repeat("\r\n", 2) + err.Error()
		}
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: errMsg,
		}
	}

	var errMsg string
	if m, ok := he.Message.(string); ok {
		errMsg = m
	} else {
		errMsg = "Unknown error"
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(he.Code)
	} else {
		err = c.String(he.Code, errMsg)
	}
	if err != nil {
		sdslog.Errorf("http error handler error: %s", err)
	}
}
