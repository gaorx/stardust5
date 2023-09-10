package sdblueprint

type Method interface {
	Id() string
	Code() MethodCode
}

type MethodCode string

var _ Method = &method{}

type method struct {
	id   string
	code MethodCode
}

func (m *method) Id() string {
	return m.id
}

func (m *method) Code() MethodCode {
	return m.code
}
