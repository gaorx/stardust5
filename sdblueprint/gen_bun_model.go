package sdblueprint

import (
	"fmt"
	"github.com/gaorx/stardust5/sdcodegen"
	"github.com/gaorx/stardust5/sdcodegen/sdgengo"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdslog"
	"github.com/gaorx/stardust5/sdtemplate"
	"github.com/samber/lo"
)

type BunModel struct {
	// tables
	TableIds     []string
	FileForModel string

	// callback
	OnHeader func(w sdcodegen.Writer, g *BunModel, bp *Blueprint)
	OnModel  func(w sdcodegen.Writer, g *BunModel, bp *Blueprint, t Table)

	// options
	Package string
}

var _ Generator = BunModel{}

func (g BunModel) GenerateTo(buffs *sdcodegen.Buffers, bp *Blueprint) error {
	tableIds := matchIds(bp.TableIds(), g.TableIds)
	if len(tableIds) <= 0 {
		return nil
	}

	// filename
	if g.FileForModel == "" {
		return sderr.New("no filename on generate go code for model")
	}

	// callback
	if g.OnHeader == nil {
		g.OnHeader = onBunHeader
	}
	if g.OnModel == nil {
		g.OnModel = onBunTable
	}

	// generate table
	for _, tableId := range tableIds {
		t := bp.Table(tableId)
		if t == nil {
			panic(sderr.NewWith("not found table", tableId))
		}
		filename, err := executeTemplate(g.FileForModel, map[string]any{"Id": t.Id()})
		if err != nil {
			return sderr.WithStack(err)
		}
		buff := buffs.Append(filename)
		if buff.IsEmpty() {
			if ok := lo.Try0(func() {
				g.OnHeader(buff, &g, bp)
			}); !ok {
				return sderr.NewWith("blueprint generate GORM error", "on_header")
			}
		}
		if ok := lo.Try0(func() {
			g.OnModel(buff, &g, bp, t)
		}); !ok {
			return sderr.NewWith("blueprint generate GORM error", "on_model")
		}
	}

	return nil
}

func onBunHeader(w sdcodegen.Writer, g *BunModel, bp *Blueprint) {
	if len(bp.tables) > 0 {
		sdgengo.Header(w, w.Filename(), g.Package, []string{
			"github.com/uptrace/bun",
		}).NL()
	} else {
		sdgengo.Header(w, w.Filename(), g.Package, nil).NL()
	}
}

func onBunTable(w sdcodegen.Writer, g *BunModel, bp *Blueprint, t Table) {
	bunTag := func(c Column) string {
		v := c.NameForDB()
		if lo.Contains(t.PrimaryKey().Columns(), c.Id()) {
			v += ",pk"
		}
		if c.IsAutoIncrement() {
			v += ",autoincrement"
		}
		if !c.IsAllowNull() {
			v += ",notnull"
		}
		return v
	}

	if t.Comment() != "" {
		w.FL("// %s %s", t.NameForGo(), t.Comment())
	}

	col2sf := func(c Column) sdgengo.Field {
		goTyp := getMemberGoType(
			c.Type(),
			c.Get("go_type").AsStr(),
			c.Get("go_import").AsStr(),
		)
		sdgengo.AddImportPackages(w, goTyp.pkgPaths)
		var tags1 []sdgengo.FieldTag
		tags1 = append(tags1, sdgengo.FieldTag{K: "bun", V: bunTag(c)})
		tags1 = appendStructFieldTagsByAttrs(tags1, c, "json", "xml", "validate")
		return sdgengo.Field{
			Name:    c.Id(),
			Type:    c.Type().String(),
			Tags:    tags1,
			Comment: c.Comment(),
		}
	}

	member2sf := func(member Field) sdgengo.Field {
		goTyp := getMemberGoType(
			member.Type(),
			member.Get("go_type").AsStr(),
			member.Get("go_import").AsStr(),
		)
		sdgengo.AddImportPackages(w, goTyp.pkgPaths)
		var tags1 []sdgengo.FieldTag
		tags1 = append(tags1, sdgengo.FieldTag{K: "bun", V: "-"})
		tags1 = appendStructFieldTagsByAttrs(tags1, member, "json", "xml", "validate")
		return sdgengo.Field{
			Name:    member.Id(),
			Type:    goTyp.typ,
			Tags:    tags1,
			Comment: member.Comment(),
		}
	}

	var genFields []sdgengo.Field
	genFields = append(genFields, sdgengo.Field{
		Name: "",
		Type: "bun.BaseModel",
		Tags: []sdgengo.FieldTag{{K: "bun", V: fmt.Sprintf("table:%s", t.NameForDB())}},
	})
	for _, c := range t.Columns() {
		genFields = append(genFields, col2sf(c))
	}
	for _, member := range t.Members() {
		genFields = append(genFields, member2sf(member))
	}

	// model
	sdgengo.Struct(w, t.NameForGo(), genFields)
	w.NL()

	// methods
	for _, method := range t.Methods() {
		code := string(method.Code())
		data := map[string]any{
			"T": t.NameForGo(),
		}
		rendered, err := sdtemplate.Text.Exec(code, data)
		if err != nil {
			sdslog.WithError(err).With("method", t.Id()+"."+method.Id()).Info("render method error")
			w.FL("// render method error %s.%s", t.Id(), method.Id())
		} else {
			w.FL(rendered)
		}
	}

	w.NL()
}
