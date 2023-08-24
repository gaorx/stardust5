package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/labstack/echo/v4"
	"io"
	"io/fs"
	"net/http"
)

type BrowserHistoryRouterStatic struct {
	PathPrefix string
	Fsys       fs.FS
	Root       string
}

func (d BrowserHistoryRouterStatic) Apply(app *echo.Echo) error {
	fsys, err := getSubFsys(d.Fsys, d.Root)
	if err != nil {
		return sderr.WithStack(err)
	}
	app.Add(
		http.MethodGet,
		d.PathPrefix+"*",
		browserHistoryRouterStaticDirectoryHandler(fsys, "index.html", false),
	)
	return nil
}

func browserHistoryRouterStaticDirectoryHandler(fsys fs.FS, recoveryFilename string, disablePathUnescaping bool) echo.HandlerFunc {
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

func getSubFsys(fsys fs.FS, root string) (fs.FS, error) {
	var fsys1 fs.FS
	if root != "" {
		sub, err := fs.Sub(fsys, root)
		if err != nil {
			return nil, sderr.WrapWith(err, "get sub dir error", root)
		}
		fsys1 = sub
	} else {
		fsys1 = fsys
	}
	return fsys1, nil
}
