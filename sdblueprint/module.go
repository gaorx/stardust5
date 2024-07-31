package sdblueprint

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/samber/lo"
	"slices"
)

type Module interface {
	Id() string
	Comment() string
	Tasks() []any
	Jsonable
}

type ModuleTaskGenerateSkeleton struct {
	Template string `json:"template"`
	Dirname  string `json:"dir"`
}

type ModuleTaskGenerateGormModel struct {
	Dirname          string   `json:"dir,omitempty"`
	Groups           []string `json:"groups,omitempty"`
	QueryWithContext bool     `json:"query_with_context,omitempty"`
}

type ModuleTaskGenerateBunModel struct {
	Dirname string   `json:"dir,omitempty"`
	Groups  []string `json:"groups,omitempty"`
}

type ModuleTaskGenerateMysqlDDL struct {
	Dirname string   `json:"dir,omitempty"`
	Groups  []string `json:"groups,omitempty"`
}

var (
	_ Module     = &module{}
	_ moduleTask = &ModuleTaskGenerateSkeleton{}
	_ moduleTask = &ModuleTaskGenerateGormModel{}
	_ moduleTask = &ModuleTaskGenerateBunModel{}
	_ moduleTask = &ModuleTaskGenerateMysqlDDL{}
)

type module struct {
	id      string
	comment string
	tasks   []moduleTask
}

type moduleTask interface {
	Jsonable
	clone() any
}

// module

func (m *module) Id() string {
	return m.id
}

func (m *module) Comment() string {
	return m.comment
}

func (m *module) Tasks() []any {
	return lo.Map(m.tasks, func(t moduleTask, _ int) any { return t.clone() })
}

func (m *module) ToJsonObject() sdjson.Object {
	if m == nil {
		return nil
	}
	return sdjson.Object{
		"id":      m.id,
		"comment": m.comment,
		"tasks":   lo.Map(m.tasks, func(t moduleTask, _ int) sdjson.Object { return t.ToJsonObject() }),
	}
}

func (m *module) setComment(comment string) *module {
	m.comment = comment
	return m
}

func (m *module) addTask(t moduleTask) {
	m.tasks = append(m.tasks, t)
}

// tasks

func (t *ModuleTaskGenerateSkeleton) clone() any {
	if t == nil {
		return nil
	}
	t1 := *t
	return &t1
}

func (t *ModuleTaskGenerateGormModel) clone() any {
	if t == nil {
		return nil
	}
	t1 := *t
	t1.Groups = slices.Clone(t.Groups)
	return &t1
}

func (t *ModuleTaskGenerateBunModel) clone() any {
	if t == nil {
		return nil
	}
	t1 := *t
	t1.Groups = slices.Clone(t.Groups)
	return &t1
}

func (t *ModuleTaskGenerateMysqlDDL) clone() any {
	if t == nil {
		return nil
	}
	t1 := *t
	t1.Groups = slices.Clone(t.Groups)
	return &t1
}

func (t *ModuleTaskGenerateSkeleton) ToJsonObject() sdjson.Object {
	if t == nil {
		return nil
	}
	return moduleTaskToJson(t, "generate_skeleton")
}

func (t *ModuleTaskGenerateGormModel) ToJsonObject() sdjson.Object {
	if t == nil {
		return nil
	}
	return moduleTaskToJson(t, "generate_gorm_model")
}

func (t *ModuleTaskGenerateBunModel) ToJsonObject() sdjson.Object {
	if t == nil {
		return nil
	}
	return moduleTaskToJson(t, "generate_bun_model")
}

func (t *ModuleTaskGenerateMysqlDDL) ToJsonObject() sdjson.Object {
	if t == nil {
		return nil
	}
	return moduleTaskToJson(t, "generate_mysql_ddl")
}

func moduleTaskToJson(v any, typ string) sdjson.Object {
	o, err := sdjson.StructToObject(v)
	if err != nil {
		panic(sderr.WithStack(err))
	}
	o["type"] = typ
	return o
}
