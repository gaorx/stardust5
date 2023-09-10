package sdblueprint

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/samber/lo"
	"reflect"
	"slices"
	"strings"
)

type Table interface {
	Id() string
	Comment() string
	Group() string
	Attributes
	Jsonable
	DummyData() []DummyRecord
	NameForTable
	Columns() []Column
	Column(id string) Column
	Members() []Field
	Member(id string) Field
	Indexes() []Index
	PrimaryKey() Index
	Methods() []Method
	Method(id string) Method
}

type Column interface {
	Field
	Jsonable
	NameForTable
	NameForJson() string
	IsPrimaryKey() bool
	IsAutoIncrement() bool
	IsAllowNull() bool
	Default() any
}

type Index interface {
	Name() string
	Comment() string
	Jsonable
	Kind() IndexKind
	Columns() []string
	ReferenceTable() string
	ReferenceColumns() []string
	Order() IndexOrder
}

type NameForTable interface {
	NameForGo() string
	NameForDB() string
}

type (
	IndexKind  string
	IndexOrder string
)

const (
	IndexSimple = IndexKind("Simple")
	IndexPK     = IndexKind("PK")
	IndexFK     = IndexKind("FK")
	IndexUnique = IndexKind("UNIQUE")
)

const (
	IndexASC  = IndexOrder("ASC")
	IndexDESC = IndexOrder("DESC")
)

var (
	_ Table  = &table{}
	_ Column = column{}
	_ Index  = &index{}
)

type table struct {
	id      string
	comment string
	group   string
	attributes
	dummyData []DummyRecord
	columns   []column
	members   []*field
	indexes   []*index
	methods   []*method
}

type column struct {
	*field
}

type index struct {
	name             string
	comment          string
	kind             IndexKind
	columns          []string
	referenceTable   string
	referenceColumns []string
	order            IndexOrder
}

// table

func (t *table) Id() string {
	return t.id
}

func (t *table) Comment() string {
	return t.comment
}

func (t *table) Group() string {
	return t.group
}

func (t *table) DummyData() []DummyRecord {
	return t.dummyData
}

func (t *table) NameForGo() string {
	name := t.Get("go").AsStr()
	if name != "" {
		return name
	}
	return sdstrings.ToCamelU(t.id)
}

func (t *table) NameForDB() string {
	name := t.Get("db").AsStr()
	if name != "" {
		return name
	}
	return "t_" + sdstrings.ToSnakeL(t.id)
}

func (t *table) Columns() []Column {
	return lo.Map(t.columns, func(c column, _ int) Column { return c })
}

func (t *table) Column(id string) Column {
	for _, c := range t.columns {
		if c.id == id {
			return c
		}
	}
	return nil
}

func (t *table) Members() []Field {
	return lo.Map(t.members, func(member *field, _ int) Field { return member })
}

func (t *table) Member(id string) Field {
	for _, member := range t.members {
		if member.id == id {
			return member
		}
	}
	return nil
}

func (t *table) Indexes() []Index {
	return lo.Map(t.indexes, func(idx *index, _ int) Index { return idx })
}

func (t *table) PrimaryKey() Index {
	for _, idx := range t.indexes {
		if idx.kind == IndexPK {
			return idx
		}
	}
	return nil
}

func (t *table) Methods() []Method {
	return lo.Map(t.methods, func(m *method, _ int) Method { return m })
}

func (t *table) Method(id string) Method {
	for _, m := range t.methods {
		if m.id == id {
			return m
		}
	}
	return nil
}

func (t *table) addFieldAsColumn(f *field) column {
	c := column{f}
	t.columns = append(t.columns, c)
	return c
}

func (t *table) addMember(f *field) *field {
	t.members = append(t.members, f)
	return f
}

func (t *table) addIndex(idx *index) *index {
	t.indexes = append(t.indexes, idx)
	return idx
}

func (t *table) addPrimaryKey(cols []string) *index {
	return t.addIndex(&index{
		kind:    IndexPK,
		columns: cols,
	})
}

func (t *table) addUniqueIndex(cols []string) *index {
	return t.addIndex(&index{
		kind:    IndexUnique,
		columns: cols,
	})
}

func (t *table) addSimpleIndex(cols []string) *index {
	return t.addIndex(&index{
		kind:    IndexSimple,
		columns: cols,
	})
}

func (t *table) addForeignKey(cols []string, refTable string, refCols []string) *index {
	return t.addIndex(&index{
		kind:             IndexFK,
		columns:          cols,
		referenceTable:   refTable,
		referenceColumns: refCols,
	})
}

func (t *table) ensurePrimaryKey() *index {
	for _, idx := range t.indexes {
		if idx.kind == IndexPK {
			return idx
		}
	}
	pk := &index{kind: IndexPK}
	newIndexes := []*index{pk}
	newIndexes = append(newIndexes, t.indexes...)
	t.indexes = newIndexes
	return pk
}

func (t *table) addMethod(m *method) *method {
	t.methods = append(t.methods, m)
	return m
}

func (t *table) setComment(comment string) *table {
	t.comment = comment
	return t
}

func (t *table) setGroup(group string) *table {
	t.group = group
	return t
}

func (t *table) ToJsonObject() sdjson.Object {
	if t == nil {
		return nil
	}
	o := sdjson.Object{
		"id":         t.id,
		"comment":    t.comment,
		"group":      t.group,
		"attributes": t.attributes.ensure(),
		"columns":    lo.Map(t.columns, func(col column, _ int) sdjson.Object { return col.ToJsonObject() }),
		"members":    lo.Map(t.members, func(member *field, _ int) sdjson.Object { return member.ToJsonObject() }),
		"indexes":    lo.Map(t.indexes, func(idx *index, _ int) sdjson.Object { return idx.ToJsonObject() }),
	}
	if len(t.dummyData) > 0 {
		o["dummy_data_count"] = len(t.dummyData)
		o["dummy_data_last"] = t.dummyData[len(t.dummyData)-1]
	}
	return o
}

// column

func (c column) IsPrimaryKey() bool {
	return c.First([]string{
		"db.primary_key",
		"db.pk",
		"primary_key",
		"pk",
	}).AsBool(false)
}

func (c column) IsAutoIncrement() bool {
	return c.First([]string{
		"db.auto_increment",
		"db.auto_incr",
		"auto_increment",
		"auto_incr",
	}).AsBool(false)
}

func (c column) IsAllowNull() bool {
	return c.First([]string{
		"db.allow_null",
		"allow_null",
	}).AsBool(false)
}

func (c column) Zero() any {
	return reflect.Zero(c.typ).Interface()
}

func (c column) Default() any {
	v, ok := c.Lookup("default")
	if !ok {
		return nil
	}
	switch c.typ.Name() {
	case "bool":
		return v.AsBool(false)
	case "string":
		return v.AsStr()
	case "byte":
		return byte(v.AsUint(0))
	case "int":
		return v.AsInt(0)
	case "int8":
		return int8(v.AsInt(0))
	case "int16":
		return int16(v.AsInt(0))
	case "int32":
		return int32(v.AsInt(0))
	case "int64":
		return v.AsInt64(0)
	case "uint":
		return v.AsUint(0)
	case "uint8":
		return uint8(v.AsUint(0))
	case "uint16":
		return uint16(v.AsUint(0))
	case "uint32":
		return uint32(v.AsUint(0))
	case "uint64":
		return v.AsUint64(0)
	case "float32":
		return float32(v.AsFloat64(0.0))
	case "float64":
		return v.AsFloat64(0.0)
	default:
		panic(sderr.NewWith("get default value error", c.typ))
	}
}

func (c column) NameForGo() string {
	name := c.Get("go").AsStr()
	if name != "" {
		return name
	}
	return c.id
}

func (c column) NameForDB() string {
	name := c.Get("db").AsStr()
	if name != "" {
		return name
	}
	return sdstrings.ToSnakeL(c.id)
}

func (c column) NameForJson() string {
	name := c.Get("json").AsStr()
	if name != "" {
		l := strings.SplitN(name, ",", 2)
		return strings.TrimSpace(l[0])
	} else {
		return sdstrings.ToSnakeL(c.id)
	}
}

func (c column) ToJsonObject() sdjson.Object {
	return c.field.ToJsonObject()
}

// index

func (idx *index) Name() string {
	return idx.name
}

func (idx *index) Comment() string {
	return idx.comment
}

func (idx *index) Kind() IndexKind {
	return idx.kind
}

func (idx *index) Columns() []string {
	return slices.Clone(idx.columns)
}

func (idx *index) ReferenceTable() string {
	return idx.referenceTable
}

func (idx *index) ReferenceColumns() []string {
	return slices.Clone(idx.referenceColumns)
}

func (idx *index) Order() IndexOrder {
	return idx.order
}

func (idx *index) addColumns(cols []string) *index {
	idx.columns = append(idx.columns, cols...)
	return idx
}

func (idx *index) setComment(comment string) *index {
	idx.comment = comment
	return idx
}

func (idx *index) setOrder(order IndexOrder) *index {
	idx.order = order
	return idx
}

func (idx *index) ToJsonObject() sdjson.Object {
	if idx == nil {
		return nil
	}
	return sdjson.Object{
		"comment":           idx.comment,
		"kind":              idx.kind,
		"columns":           idx.columns,
		"reference_table":   idx.referenceTable,
		"reference_columns": idx.referenceColumns,
		"order":             idx.order,
	}
}
