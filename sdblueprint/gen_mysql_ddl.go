package sdblueprint

import (
	"fmt"
	"github.com/gaorx/stardust5/sdcodegen"
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
	"slices"
	"strings"
	"unicode"
)

type MysqlDDL struct {
	TableIds      []string
	FileForCreate string
	FileForDrop   string

	// callbacks
	OnHeader func(w sdcodegen.Writer, g *MysqlDDL, bp *Blueprint)
	OnCreate func(w sdcodegen.Writer, g *MysqlDDL, bp *Blueprint, t Table)
	OnDrop   func(w sdcodegen.Writer, g *MysqlDDL, bp *Blueprint, t Table)

	// options
	Charset       string
	Collation     string
	Engine        string
	DisableFK     bool
	WithDrop      bool
	WithoutCreate bool
}

var _ Generator = MysqlDDL{}

func (g MysqlDDL) GenerateTo(buffs *sdcodegen.Buffers, bp *Blueprint) error {
	tableIds := matchIds(bp.TableIds(), g.TableIds)
	if len(tableIds) <= 0 {
		return nil
	}

	// filename
	if !g.WithoutCreate {
		if g.FileForCreate == "" {
			return sderr.New("no filename on generate DDL for create tables")
		}
	}
	if g.WithDrop {
		if g.FileForDrop == "" {
			g.FileForDrop = g.FileForCreate
		}
		if g.FileForDrop == "" {
			return sderr.New("no filename on generate DDL for drop tables")
		}
	}

	// callbacks
	if g.OnHeader == nil {
		g.OnHeader = onMysqlHeader
	}
	if g.OnCreate == nil {
		g.OnCreate = onMysqlCreate
	}

	if g.OnDrop == nil {
		g.OnDrop = onMysqlDrop
	}

	// options
	if g.Charset == "" {
		g.Charset = "utf8mb4"
	}
	if g.Collation == "" {
		g.Collation = "utf8mb4_general_ci"
	}
	if g.Engine == "" {
		g.Engine = "InnoDB"
	}

	// drop table
	if g.WithDrop {
		for i, tableId := range lo.Reverse(slices.Clone(tableIds)) {
			t := bp.Table(tableId)
			if t == nil {
				panic(sderr.NewWith("not found table", t.Id()))
			}
			filename, err := executeTemplate(g.FileForDrop, map[string]any{"Id": t.Id()})
			if err != nil {
				return sderr.WithStack(err)
			}
			buff := buffs.Append(filename)
			if buff.IsEmpty() {
				if ok := lo.Try0(func() {
					g.OnHeader(buff, &g, bp)
				}); !ok {
					return sderr.NewWith("blueprint generate mysql DDL error", "on_header")
				}
			}
			if ok := lo.Try0(func() {
				g.OnDrop(buff, &g, bp, t)
			}); !ok {
				return sderr.NewWith("blueprint generate mysql DDL error", "on_drop")
			}
			if i >= len(tableIds)-1 {
				buff.NL()
			}
		}
	}

	// create table
	if !g.WithoutCreate {
		for _, tableId := range tableIds {
			t := bp.Table(tableId)
			if t == nil {
				panic(sderr.NewWith("not found table", t.Id()))
			}
			filename, err := executeTemplate(g.FileForCreate, map[string]any{"Id": t.Id()})
			if err != nil {
				return sderr.WithStack(err)
			}
			buff := buffs.Append(filename)
			if buff.IsEmpty() {
				if ok := lo.Try0(func() {
					g.OnHeader(buff, &g, bp)
				}); !ok {
					return sderr.NewWith("blueprint generate mysql DDL error", "on_header")
				}
			}
			if ok := lo.Try0(func() {
				g.OnCreate(buff, &g, bp, t)
			}); !ok {
				return sderr.NewWith("blueprint generate mysql DDL error", "on_create")
			}
		}
	}
	return nil
}

func onMysqlHeader(w sdcodegen.Writer, _ *MysqlDDL, _ *Blueprint) {
	w.NL()
	w.L("-- AUTO GENERATED, DO NOT EDIT")
	w.L("-- AUTO GENERATED, DO NOT EDIT")
	w.L("-- AUTO GENERATED, DO NOT EDIT")
	w.NL()
}

func onMysqlCreate(b sdcodegen.Writer, g *MysqlDDL, bp *Blueprint, t Table) {
	q := func(name string) string {
		return "`" + name + "`"
	}

	ql := func(names []string) string {
		if len(names) <= 0 {
			return ""
		}
		return strings.Join(lo.Map(names, func(name string, _ int) string {
			return q(name)
		}), ",")
	}

	toDbColumnNames := func(t Table, colNames []string) []string {
		return lo.Map(colNames, func(colName string, _ int) string {
			return t.Column(colName).NameForDB()
		})
	}

	defToStr := func(v any) string {
		if v1, ok := v.(string); ok {
			return fmt.Sprintf("'%s'", v1)
		} else {
			return fmt.Sprint(v)
		}
	}

	b.FL("-- %s", t.Id())
	b.FL("CREATE TABLE IF NOT EXISTS %s (", q(t.NameForDB()))
	for _, c := range t.Columns() {
		b.I(1)
		b.F("%s %s", q(c.NameForDB()), mysqlDataTypeOf(c))
		b.If(c.IsAutoIncrement(), " AUTO_INCREMENT")
		b.If(!c.IsAllowNull(), " NOT NULL")
		if c.Default() != nil {
			b.F(" DEFAULT %s", defToStr(c.Default()))
		}
		b.P(",")
		b.If(c.Comment() != "", " -- "+c.Comment())
		b.NL()
	}
	b.NL()
	for _, idx := range t.Indexes() {
		b.I(1)
		switch idx.Kind() {
		case IndexPK:
			b.F("PRIMARY KEY (%s)", ql(toDbColumnNames(t, idx.Columns()))).P(",")
		case IndexFK:
			if !g.DisableFK {
				refTable := bp.Table(idx.ReferenceTable())
				if refTable == nil {
					panic(sderr.NewWith("not found foreign key reference table in query", idx.ReferenceTable))
				}
				if idx.Name() != "" {
					b.F("CONSTRAINT %s", q(idx.Name()))
				}
				b.F("FOREIGN KEY (%s) REFERENCES %s (%s)",
					ql(toDbColumnNames(t, idx.Columns())),
					q(refTable.NameForDB()),
					ql(toDbColumnNames(refTable, idx.ReferenceColumns())),
				)
				b.P(",")
			}
		case IndexSimple:
			if idx.Name() != "" {
				b.F("INDEX %s (%s)", q(idx.Name()), ql(toDbColumnNames(t, idx.Columns())))
			} else {
				b.F("INDEX (%s)", ql(toDbColumnNames(t, idx.Columns())))
			}
			b.P(",")
		case IndexUnique:
			if idx.Name() != "" {
				b.F("UNIQUE INDEX %s (%s)", q(idx.Name()), ql(toDbColumnNames(t, idx.Columns())))
			} else {
				b.F("UNIQUE INDEX (%s)", ql(toDbColumnNames(t, idx.Columns())))
			}
			b.P(",")
		}
		b.NL()
	}
	b.Modify(func(code string) string {
		code = strings.TrimRightFunc(code, func(c rune) bool {
			return unicode.IsSpace(c)
		})
		code = strings.TrimSuffix(code, ",")
		return code
	})
	b.NL()
	b.F(")")
	b.If(g.Engine != "", " ENGINE="+g.Engine)
	b.If(g.Charset != "", " CHARACTER SET="+g.Charset)
	b.If(g.Collation != "", " COLLATE="+g.Collation)
	b.L(";")
	b.NL()
}

func onMysqlDrop(w sdcodegen.Writer, g *MysqlDDL, bp *Blueprint, t Table) {
	q := func(name string) string {
		return "`" + name + "`"
	}
	w.FL("DROP TABLE IF EXISTS %s;", q(t.NameForDB()))
}

func mysqlDataTypeOf(c Column) string {
	dbTyp := c.First([]string{"db_type", "dbtype"}).AsStr()
	if dbTyp != "" {
		return dbTyp
	}
	switch c.Type().String() {
	case "string":
		return "VARCHAR(255)"
	case "bool":
		return "TINYINT"
	case "[]byte":
		return "MEDIUMBLOB"
	case "int", "int8", "int16", "int32":
		return "INT"
	case "int64":
		return "BIGINT"
	case "uint", "uint8", "uint16", "uint32":
		return "INT UNSIGNED"
	case "uint64":
		return "BIGINT UNSIGNED"
	case "float32", "float64":
		return "DOUBLE"
	default:
		panic("illegal type for convert to MYSQL data type")
	}
}
