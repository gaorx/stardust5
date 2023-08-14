package sdecho

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdtemplate"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"io"
	"io/fs"
)

type htmlRenderer struct {
	loader *sdtemplate.HtmlLoader
}

func MustHtmlRenderer(fsys fs.FS, eager bool) echo.Renderer {
	return lo.Must(NewHtmlRenderer(fsys, eager))
}

func NewHtmlRenderer(fsys fs.FS, eager bool) (echo.Renderer, error) {
	loader, err := sdtemplate.NewHtmlLoader(fsys, &sdtemplate.HtmlLoaderOptions{
		Eager: eager,
	})
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return &htmlRenderer{loader: loader}, nil
}

func (renderer htmlRenderer) Render(wr io.Writer, name string, data any, ec echo.Context) error {
	t, err := renderer.loader.Load(name)
	if err != nil {
		return sderr.WithStack(err)
	}
	err = t.Execute(wr, data)
	if err != nil {
		return sderr.WrapWith(err, "execute template error", name)
	}
	return nil
}
