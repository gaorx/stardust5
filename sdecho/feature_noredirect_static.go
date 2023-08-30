package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/labstack/echo/v4"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type NoRedirectStatic struct {
	PathPrefix       string
	Fsys             fs.FS
	Root             string
	TrimPathPrefixes []string
}

func (d NoRedirectStatic) Apply(app *echo.Echo) error {
	fsys, err := getSubFsys(d.Fsys, d.Root)
	if err != nil {
		return sderr.WithStack(err)
	}
	app.Add(
		http.MethodGet,
		d.PathPrefix+"*",
		noRedirectStaticDirectoryHandler(fsys, d.TrimPathPrefixes, false),
	)
	return nil
}

const (
	defaultIndexPage = "index.html"
)

func noRedirectStaticDirectoryHandler(fsys fs.FS, trimPathPrefixes []string, disablePathUnescaping bool) echo.HandlerFunc {
	return func(ec echo.Context) error {
		return noRedirectStaticDirectory(ec, fsys, trimPathPrefixes, disablePathUnescaping)
	}
}

func noRedirectStaticDirectory(
	ec echo.Context,
	fsys fs.FS,
	trimPathPrefixes []string,
	disablePathUnescaping bool,
) error {
	p := ec.Param("*")
	if !disablePathUnescaping {
		tmpPath, err := url.PathUnescape(p)
		if err != nil {
			return sderr.NewWith("failed to unescape path variable", err)
		}
		p = tmpPath
	}

	toName := func(p string, trimPrefix string) string {
		if trimPrefix != "" {
			p = strings.TrimPrefix(p, strings.TrimPrefix(trimPrefix, "/"))
		}
		return filepath.ToSlash(filepath.Clean(strings.TrimPrefix(p, "/")))
	}

	candidateNames := []string{toName(p, "")}
	for _, trimPathPrefix := range trimPathPrefixes {
		candidateNames = append(candidateNames, toName(p, trimPathPrefix))
	}
	fi, name, err := fsStatFirst(fsys, candidateNames)
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

func fsStatFirst(fsys fs.FS, names []string) (fs.FileInfo, string, error) {
	var firstErr error = nil
	for _, name := range names {
		fi, err := fs.Stat(fsys, name)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
		} else {
			return fi, name, nil
		}
	}
	return nil, "", firstErr
}
