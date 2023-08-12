package sdecho

import (
	"fmt"
	"github.com/gaorx/stardust5/sdcall"
	"log/slog"
	"time"

	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdparse"
	"github.com/gaorx/stardust5/sdtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LoggingRecover(logSkipper middleware.Skipper) echo.MiddlewareFunc {
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
			return sderr.Wrap(finalErr, "sdecho logging recover middleware error")
		}
	}
}
