package sdblueprint

import (
	"github.com/gaorx/stardust5/sdcodegen"
	"github.com/gaorx/stardust5/sderr"
	"path/filepath"
	"reflect"
	"strings"
)

type ForModule struct {
	Ids []string
}

var _ Generator = ForModule{}

func (g ForModule) GenerateTo(buffs *sdcodegen.Buffers, bp *Blueprint) error {
	moduleIds := matchIds(bp.ModuleIds(), g.Ids)
	if len(moduleIds) <= 0 {
		return nil
	}
	for _, moduleId := range moduleIds {
		m := bp.Module(moduleId)
		if m == nil {
			panic(sderr.NewWith("not found module", moduleId))
		}
		for _, t := range m.Tasks() {
			switch t1 := t.(type) {
			case *ModuleTaskGenerateSkeleton:
				if t1.Dirname == "" {
					return sderr.NewWith("no dir in generate skeleton task", m.Id())
				}
				if t1.Template == "" {
					return sderr.NewWith("no template in generate skeleton task", m.Id())
				}
				if err := generateSkeleton(buffs, t1.Dirname, t1.Template); err != nil {
					return sderr.WithStack(err)
				}
			case *ModuleTaskGenerateGormModel:
				dirname := t1.Dirname
				if dirname == "" {
					return sderr.NewWith("no dir in generate GORM model task", m.Id())
				}
				if err := (GormModel{
					FileForModel:     filepath.Join(dirname, "models.gen.go"),
					WithQuery:        true,
					QueryWithContext: t1.QueryWithContext,
				}).GenerateTo(buffs, getSub(bp, t1.Groups)); err != nil {
					return sderr.WithStack(err)
				}
			case *ModuleTaskGenerateMysqlDDL:
				dirname := t1.Dirname
				if dirname == "" {
					return sderr.NewWith("no dir in generate MYSQL ddl task", m.Id())
				}
				if err := (MysqlDDL{
					FileForCreate: filepath.Join(dirname, "create_table.gen.sql"),
					FileForDrop:   filepath.Join(dirname, "drop_table.gen.sql"),
					WithDrop:      true,
				}).GenerateTo(buffs, getSub(bp, t1.Groups)); err != nil {
					return sderr.WithStack(err)
				}
			default:
				panic(sderr.NewWith("illegal task", sderr.Attrs{"task": reflect.TypeOf(t1).String(), "model": m.Id()}))
			}
		}
	}
	return nil
}

func generateSkeleton(buffs *sdcodegen.Buffers, dirname string, template string) error {
	switch strings.ToLower(template) {
	case "cmd", "simple":
		return generateSimpleCmd(buffs, dirname)
	case "web_admin_antd":
		panic("TODO: generate_admin_antd")
	default:
		return sderr.NewWith("illegal template in generate skeleton task", template)
	}
}

func generateSimpleCmd(buffs *sdcodegen.Buffers, dirname string) error {
	buffs.Add(filepath.Join(dirname, "main.go"), func(w sdcodegen.Writer) {
		w.SetOverwrite(false)
		w.L("package main")
		w.NL()
		w.L("import (")
		w.I(1).L("\"fmt\"")
		w.L(")")
		w.NL()
		w.L("func main() {")
		w.I(1).L("fmt.Println(\"TODO: add code\")")
		w.L("}")
		w.NL()
	})
	return nil
}

func getSub(bp *Blueprint, groups []string) *Blueprint {
	if len(groups) <= 0 {
		return bp
	}
	return bp.Sub(groups...)
}
