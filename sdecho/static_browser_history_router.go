package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/labstack/echo/v4"
	"io"
	"io/fs"
	"net/http"
)

func BrowserHistoryRouterStaticFS(app *echo.Echo, pathPrefix string, fsys fs.FS) *echo.Route {
	return app.Add(
		http.MethodGet,
		pathPrefix+"*",
		BrowserHistoryRouterStaticDirectoryHandler(fsys, "index.html", false),
	)
}

func BrowserHistoryRouterStaticDirectoryHandler(fsys fs.FS, recoveryFilename string, disablePathUnescaping bool) echo.HandlerFunc {
	return func(ec echo.Context) error {
		err := noRedirectStaticDirectory(ec, fsys, disablePathUnescaping)
		if err != nil {
			if httpErr, ok := sderr.AsT[*echo.HTTPError](err); ok && httpErr != nil && httpErr.Code == http.StatusNotFound {
				f, err1 := fsys.Open(recoveryFilename)
				if err1 == nil && f != nil {
					defer func() { _ = f.Close() }()
					ff, ok := f.(io.ReadSeeker)
					if !ok {
						return sderr.New("file does not implement io.ReadSeeker")
					}
					fi, _ := f.Stat()
					http.ServeContent(ec.Response(), ec.Request(), fi.Name(), fi.ModTime(), ff)
					return nil
				}
			}
		}
		return err
	}
}
