package sdblueprint

import (
	"fmt"
	"github.com/gaorx/stardust5/sdcodegen"
	"github.com/gaorx/stardust5/sdcodegen/sdgengo"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdslog"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/gaorx/stardust5/sdtemplate"
	"github.com/samber/lo"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type GormModel struct {
	// tables
	TableIds     []string
	FileForModel string

	// queries
	WithQuery        bool
	QueryIds         []string
	FileForQuery     string
	QueryWithContext bool

	// callback
	OnHeader func(w sdcodegen.Writer, g *GormModel, bp *Blueprint)
	OnModel  func(w sdcodegen.Writer, g *GormModel, bp *Blueprint, t Table)
	OnQuery  func(w sdcodegen.Writer, g *GormModel, bp *Blueprint, q Query)

	// options
	Package string
}

var _ Generator = GormModel{}

func (g GormModel) GenerateTo(buffs *sdcodegen.Buffers, bp *Blueprint) error {
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
		g.OnHeader = onGormHeader
	}
	if g.OnModel == nil {
		g.OnModel = onGormTable
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

	// queries
	if !g.WithQuery {
		return nil
	}

	queryIds := matchIds(bp.QueryIds(), g.QueryIds)
	if len(queryIds) <= 0 {
		return nil
	}

	// query filename
	if g.FileForQuery == "" {
		g.FileForQuery = g.FileForModel
	}

	// query callback
	if g.OnQuery == nil {
		g.OnQuery = onGormQuery
	}

	// generate query
	for _, queryId := range queryIds {
		q := bp.Query(queryId)
		if q == nil {
			panic(sderr.NewWith("not found query", queryId))
		}
		filename, err := executeTemplate(g.FileForQuery, map[string]any{"Id": q.Id()})
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
			g.OnQuery(buff, &g, bp, q)
		}); !ok {
			return sderr.NewWith("blueprint generate GORM error", "on_query")
		}
	}

	return nil
}

func onGormHeader(w sdcodegen.Writer, g *GormModel, bp *Blueprint) {
	if g.WithQuery && len(bp.queries) > 0 {
		sdgengo.Header(w, w.Filename(), g.Package, []string{
			"github.com/gaorx/stardust5/sderr",
			"gorm.io/gorm",
		}).NL()
	} else {
		sdgengo.Header(w, w.Filename(), g.Package, nil).NL()
	}
}

func onGormTable(w sdcodegen.Writer, g *GormModel, bp *Blueprint, t Table) {
	gormTag := func(c Column) string {
		v := "column:" + c.NameForDB()
		if lo.Contains(t.PrimaryKey().Columns(), c.Id()) {
			v += ";primaryKey"
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
		tags1 = append(tags1, sdgengo.FieldTag{K: "gorm", V: gormTag(c)})
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
		tags1 = append(tags1, sdgengo.FieldTag{K: "gorm", V: "-:all"})
		tags1 = appendStructFieldTagsByAttrs(tags1, member, "json", "xml", "validate")
		return sdgengo.Field{
			Name:    member.Id(),
			Type:    goTyp.typ,
			Tags:    tags1,
			Comment: member.Comment(),
		}
	}

	var genFields []sdgengo.Field
	for _, c := range t.Columns() {
		genFields = append(genFields, col2sf(c))
	}
	for _, member := range t.Members() {
		genFields = append(genFields, member2sf(member))
	}

	// model
	sdgengo.Struct(w, t.NameForGo(), genFields)
	w.NL()

	// table name
	sdgengo.Method(
		w,
		"TableName",
		sdgengo.NamedType{Type: t.NameForGo()},
		nil,
		[]sdgengo.NamedType{{Type: "string"}},
		func(w sdcodegen.Writer) {
			w.I(1).FL("return \"%s\"", t.NameForDB())
		},
	)
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

func onGormQuery(w sdcodegen.Writer, g *GormModel, _ *Blueprint, q Query) {
	getGenParams := func() []sdgengo.NamedType {
		var pairs []sdgengo.NamedType
		if g.QueryWithContext {
			sdgengo.AddImportPackages(w, []string{"context"})
			pairs = append(pairs, sdgengo.NamedType{Name: "ctx", Type: "context.Context"})
		}
		pairs = append(pairs, sdgengo.NamedType{Name: "tx", Type: "*gorm.DB"})
		for _, param := range q.Params() {
			pairs = append(pairs, sdgengo.NamedType{Name: param.Name(), Type: param.Type()})
		}
		return pairs
	}

	genDbrErr := func(w sdcodegen.Writer, dbr string, returnRowAffected bool) {
		if !returnRowAffected {
			w.I(1).FL("if %s.Error != nil {", dbr)
			w.I(2).FL("return sderr.WithStack(%s.Error)", dbr)
			w.I(1).FL("}")
			w.I(1).FL("return nil")
		} else {
			w.I(1).FL("if %s.Error != nil {", dbr)
			w.I(2).FL("return 0, sderr.WithStack(%s.Error)", dbr)
			w.I(1).FL("}")
			w.I(1).FL("return %s.RowsAffected, nil", dbr)
		}
	}

	joinNamedParams := func() string {
		if len(q.Params()) <= 0 {
			return ""
		}
		sdgengo.AddImportPackages(w, []string{"database/sql"})
		return ", " + strings.Join(lo.Map(q.Params(), func(param QueryParam, _ int) string {
			return fmt.Sprintf("sql.Named(\"%s\", %s)", param.Name(), param.Name())
		}), ", ")
	}

	withContext := func(q string) string {
		if g.QueryWithContext {
			return strings.Replace(q, "tx.", "tx.WithContext(ctx).", -1)
		} else {
			return q
		}
	}

	k := q.Kind()
	switch k {
	case QueryForCreate:
		if !q.ReturnRowAffected() {
			sdgengo.Func(w, q.Id(), getGenParams(), sdgengo.Return("error"), func(w sdcodegen.Writer) {
				w.I(1).FL(withContext("dbr := tx.Create(%s)"), q.ParamByIndex(0).Name())
				genDbrErr(w, "dbr", false)
			}).NL()
		} else {
			sdgengo.Func(w, q.Id(), getGenParams(), sdgengo.Return("int64", "error"), func(w sdcodegen.Writer) {
				w.I(1).FL(withContext("dbr := tx.Create(%s)"), q.ParamByIndex(0).Name())
				genDbrErr(w, "dbr", true)
			}).NL()
		}
	case QueryForUpdate:
		if !q.ReturnRowAffected() {
			sdgengo.Func(w, q.Id(), getGenParams(), sdgengo.Return("error"), func(w sdcodegen.Writer) {
				w.I(1).FL(withContext("dbr := tx.Save(%s)"), q.ParamByIndex(0).Name())
				genDbrErr(w, "dbr", false)
			}).NL()
		} else {
			sdgengo.Func(w, q.Id(), getGenParams(), sdgengo.Return("int64", "error"), func(w sdcodegen.Writer) {
				w.I(1).FL(withContext("dbr := tx.Save(%s)"), q.ParamByIndex(0).Name())
				genDbrErr(w, "dbr", true)
			}).NL()
		}
	case QueryForExec:
		if !q.ReturnRowAffected() {
			sdgengo.Func(w, q.Id(), getGenParams(), sdgengo.Return("error"), func(w sdcodegen.Writer) {
				w.I(1).FL(withContext("dbr := tx.Exec(%s%s)"), strconv.Quote(q.SQL()), joinNamedParams())
				genDbrErr(w, "dbr", false)
			}).NL()
		} else {
			sdgengo.Func(w, q.Id(), getGenParams(), sdgengo.Return("int64", "error"), func(w sdcodegen.Writer) {
				w.I(1).FL(withContext("dbr := tx.Exec(%s%s)"), strconv.Quote(q.SQL()), joinNamedParams())
				genDbrErr(w, "dbr", true)
			}).NL()
		}
	case QueryForRecord:
		sdgengo.Func(w, q.Id(), getGenParams(), sdgengo.Return(q.Result().Type(), "error"), func(w sdcodegen.Writer) {
			w.I(1).FL("var row %s", q.Result().Type())
			w.I(1).FL(withContext("dbr := tx.Raw(%s%s).Take(&row)"), strconv.Quote(q.SQL()), joinNamedParams())
			w.I(1).FL("if dbr.Error != nil {")
			w.I(2).FL("return nil, sderr.WithStack(dbr.Error)")
			w.I(1).FL("}")
			w.I(1).FL("return row, nil")
		}).NL()
	case QueryForRecords:
		sdgengo.Func(w, q.Id(), getGenParams(), sdgengo.Return(q.Result().Type(), "error"), func(w sdcodegen.Writer) {
			w.I(1).FL("var rows %s", q.Result().Type())
			w.I(1).FL(withContext("dbr := tx.Raw(%s%s).Find(&rows)"), strconv.Quote(q.SQL()), joinNamedParams())
			w.I(1).FL("if dbr.Error != nil {")
			w.I(2).FL("return nil, sderr.WithStack(dbr.Error)")
			w.I(1).FL("}")
			w.I(1).FL("return rows, nil")
		}).NL()
	case QueryForScalar:
		sdgengo.Func(w, q.Id(), getGenParams(), sdgengo.Return(q.Result().Type(), "error"), func(w sdcodegen.Writer) {
			w.I(1).FL("var r %s", q.Result().Type())
			w.I(1).FL(withContext("dbr := tx.Raw(%s%s).Scan(&r)"), strconv.Quote(q.SQL()), joinNamedParams())
			w.I(1).FL("if dbr.Error != nil {")
			w.I(2).FL("return r, sderr.WithStack(dbr.Error)")
			w.I(1).FL("}")
			w.I(1).FL("return r, nil")
		}).NL()
	default:
		panic(sderr.NewWith("illegal kind in query for generate code", sderr.Attrs{"kind": k, "q": q.Id()}))
	}
}

type goModelFieldType struct {
	pkgPaths []string
	typ      string
}

func getMemberGoType(typ reflect.Type, goTyp, goImport string) goModelFieldType {
	trimPkgPaths := func(l []string) []string {
		l1 := lo.Filter(l, func(x string, _ int) bool {
			return x != ""
		})
		slices.Sort(l1)
		return l1
	}

	if goTyp != "" {
		return goModelFieldType{
			pkgPaths: trimPkgPaths(sdstrings.SplitNonempty(goImport, ";", true)),
			typ:      strings.TrimSpace(goTyp),
		}
	} else {
		var getPkgPaths func(reflect.Type) []string
		getPkgPaths = func(typ0 reflect.Type) []string {
			if typ0.PkgPath() != "" {
				return []string{typ0.PkgPath()}
			} else {
				switch typ0.Kind() {
				case reflect.Pointer, reflect.Slice, reflect.Array, reflect.Chan:
					return getPkgPaths(typ0.Elem())
				case reflect.Map:
					keyImportPkgs := getPkgPaths(typ0.Key())
					valImportPkgs := getPkgPaths(typ0.Elem())
					return append(keyImportPkgs, valImportPkgs...)
				case reflect.Func:
					var importPkgs []string
					for i := 0; i < typ0.NumIn(); i++ {
						importPkgs = append(importPkgs, getPkgPaths(typ0.In(i))...)
					}
					for i := 0; i < typ0.NumOut(); i++ {
						importPkgs = append(importPkgs, getPkgPaths(typ0.Out(i))...)
					}
					return importPkgs
				default:
					return []string{}
				}
			}
		}
		return goModelFieldType{
			pkgPaths: trimPkgPaths(getPkgPaths(typ)),
			typ:      typ.String(),
		}
	}
}
