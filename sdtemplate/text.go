package sdtemplate

import (
	"bytes"
	"text/template"

	"github.com/gaorx/stardust5/sderr"
)

type textExecutor struct {
}

func (te textExecutor) Exec(tmpl string, data any) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", sderr.Wrap(err, "parse text template error")
	}
	buff := bytes.NewBufferString("")
	err = t.Execute(buff, data)
	if err != nil {
		return "", sderr.Wrap(err, "execute text template error")
	}
	return buff.String(), nil
}

func (te textExecutor) ExecDef(tmpl string, data any, def string) string {
	r, err := te.Exec(tmpl, data)
	if err != nil {
		return def
	}
	return r
}
