package sdecho

import (
	"fmt"
	"github.com/gaorx/stardust5/sdcall"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdparse"
	"github.com/gaorx/stardust5/sdslog"
	"github.com/gaorx/stardust5/sdtime"
	"github.com/samber/lo"
	"log/slog"
	"net/http"
	"strings"
	"time"

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
	app.Use(loggingRecover(opts1.LogSkipper))
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

func loggingRecover(logSkipper middleware.Skipper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			if logSkipper != nil && logSkipper(ec) {
				return next(ec)
			}
			req := ec.Request()
			res := ec.Response()
			startAt := time.Now()
			var nextErr, panicErr, finalErr error
			panicErr = sdcall.Safe(func() {
				nextErr = next(ec)
			})
			if panicErr != nil {
				finalErr = panicErr
			} else {
				finalErr = nextErr
			}
			if finalErr != nil {
				ec.Error(finalErr)
			}
			elapsedHuman := time.Since(startAt)
			elapsedMs := sdtime.ToMillisF(elapsedHuman)
			statusCode := res.Status
			method := req.Method
			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			bytesIn, err := sdparse.Int64(req.Header.Get(echo.HeaderContentLength))
			if err != nil {
				bytesIn = 0
			}

			logAttrs := []any{
				slog.Float64("latency", elapsedMs),
				slog.Duration("latency_h", elapsedHuman),
				slog.String("remote_ip", ec.RealIP()),
				slog.Int64("bytes_in", bytesIn),
				slog.Int64("bytes_out", res.Size),
			}
			if finalErr == nil {
				slog.With(logAttrs...).Info(fmt.Sprintf("%d %s %s", statusCode, method, path))
			} else {
				logAttrs = append(logAttrs, slog.String("error", fmt.Sprintf("%+v", finalErr)))
				slog.With(logAttrs...).Info(fmt.Sprintf("%d %s %s", statusCode, method, path))
			}
			return sderr.Wrap(finalErr, "logging recover middleware error")
		}
	}
}
