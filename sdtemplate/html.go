package sdtemplate

import (
	"bytes"
	"html/template"

	"github.com/gaorx/stardust5/sderr"
)

type htmlExecutor struct {
}

func (te htmlExecutor) Exec(tmpl string, data any) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", sderr.Wrap(err, "parse html template error")
	}
	buff := bytes.NewBufferString("")
	err = t.Execute(buff, data)
	if err != nil {
		return "", sderr.Wrap(err, "execute html template error")
	}
	return buff.String(), nil
}

func (te htmlExecutor) ExecDef(tmpl string, data any, def string) string {
	r, err := te.Exec(tmpl, data)
	if err != nil {
		return def
	}
	return r
}
