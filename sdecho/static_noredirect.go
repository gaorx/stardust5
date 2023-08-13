package sdecho

import (
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gaorx/stardust5/sderr"
	"github.com/labstack/echo/v4"
)

const (
	defaultIndexPage = "index.html"
)

func NoRedirectStaticFS(app *echo.Echo, pathPrefix string, fsys fs.FS) *echo.Route {
	return app.Add(
		http.MethodGet,
		pathPrefix+"*",
		NoRedirectStaticDirectoryHandler(fsys, false),
	)
}

func NoRedirectStaticDirectoryHandler(fsys fs.FS, disablePathUnescaping bool) echo.HandlerFunc {
	return func(ec echo.Context) error {
		return noRedirectStaticDirectory(ec, fsys, disablePathUnescaping)
	}
}

func noRedirectStaticDirectory(ec echo.Context, fsys fs.FS, disablePathUnescaping bool) error {
	p := ec.Param("*")
	if !disablePathUnescaping {
		tmpPath, err := url.PathUnescape(p)
		if err != nil {
			return sderr.NewWith("failed to unescape path variable", err)
		}
		p = tmpPath
	}

	name := filepath.ToSlash(filepath.Clean(strings.TrimPrefix(p, "/")))
	fi, err := fs.Stat(fsys, name)
	if err != nil {
		return echo.ErrNotFound
	}

	p = ec.Request().URL.Path
	if fi.IsDir() {
		name = defaultIndexPage
	}
	return fsFile2(ec, name, fsys)
}

func fsFile2(ec echo.Context, file string, fsys fs.FS) error {
	f, err := fsys.Open(file)
	if err != nil {
		return echo.ErrNotFound
	}
	defer func() { _ = f.Close() }()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.ToSlash(filepath.Join(file, defaultIndexPage))
		f, err = fsys.Open(file)
		if err != nil {
			return echo.ErrNotFound
		}
		defer func() { _ = f.Close() }()
		if fi, err = f.Stat(); err != nil {
			return err
		}
	}
	ff, ok := f.(io.ReadSeeker)
	if !ok {
		return sderr.New("file does not implement io.ReadSeeker")
	}
	http.ServeContent(ec.Response(), ec.Request(), fi.Name(), fi.ModTime(), ff)
	return nil
}
