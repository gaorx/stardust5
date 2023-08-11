package sdtemplate

type TemplateExecutor interface {
	Exec(template string, data any) (string, error)
	ExecDef(template string, data any, def string) string
}
