package sdblueprint

import (
	"fmt"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdreflect"
	"github.com/samber/lo"
	"io"
	"reflect"
)

type Blueprint struct {
	session          any
	protos           []any
	finalized        bool
	tables           []*table
	queries          []*query
	modules          []*module
	disableMethod    bool
	disableDummyData bool
}

func New(session any) *Blueprint {
	return &Blueprint{
		session:   session,
		finalized: false,
	}
}

func (bp *Blueprint) Add(protos ...any) *Blueprint {
	var add func(proto any)
	add = func(proto any) {
		if proto == nil {
			return
		}
		protoVal := sdreflect.ValueOf(proto)
		if protoVal.Kind() == reflect.Slice || protoVal.Kind() == reflect.Array {
			n := protoVal.Len()
			for i := 0; i < n; i++ {
				proto0Val := protoVal.Index(i)
				add(proto0Val.Interface())
			}
		} else {
			bp.protos = append(bp.protos, proto)
		}
	}

	for _, proto := range protos {
		add(proto)
	}
	return bp
}

func (bp *Blueprint) IsFinalized() bool {
	return bp.finalized
}

func (bp *Blueprint) MustFinalize() *Blueprint {
	lo.Must0(bp.Finalize())
	return bp
}

func (bp *Blueprint) Finalize() error {
	if bp.finalized {
		return nil
	}
	bpCopy := &Blueprint{
		session:   bp.session,
		protos:    bp.protos,
		finalized: false,
	}

	// scan
	for _, proto := range bp.protos {
		st, ok := structTypeOf(proto)
		if !ok {
			return sderr.New("prototype is not struct")
		}
		sv := sdreflect.ValueOf(proto)
		if err := scanProto(bpCopy, sv, st, ""); err != nil {
			return sderr.WithStack(err)
		}
	}

	// check
	if err := bpCopy.check(); err != nil {
		return sderr.WithStack(err)
	}

	// expand SQL in query
	if err := expandQuery(bpCopy); err != nil {
		return sderr.WithStack(err)
	}

	// ok
	bp.tables = bpCopy.tables
	bp.queries = bpCopy.queries
	bp.modules = bpCopy.modules
	bp.finalized = true
	return nil
}

func (bp *Blueprint) ToJsonObject() sdjson.Object {
	if bp == nil {
		return nil
	}
	return sdjson.Object{
		"finalized": bp.finalized,
		"tables":    lo.Map(bp.tables, func(t *table, _ int) sdjson.Object { return t.ToJsonObject() }),
		"queries":   lo.Map(bp.queries, func(q *query, _ int) sdjson.Object { return q.ToJsonObject() }),
	}
}

func (bp *Blueprint) Dump(w io.Writer) {
	if w == nil {
		return
	}
	_, _ = fmt.Fprintln(w, sdjson.MarshalPretty(bp.ToJsonObject()))
}

func (bp *Blueprint) Sub(groups ...string) *Blueprint {
	if !bp.finalized {
		panic(sderr.New("blueprint is not finalized"))
	}
	return &Blueprint{
		protos:    nil,
		finalized: true,
		tables: lo.Filter(bp.tables, func(t *table, _ int) bool {
			return lo.Contains(groups, t.group)
		}),
		queries: lo.Filter(bp.queries, func(q *query, _ int) bool {
			return lo.Contains(groups, q.group)
		}),
	}
}

func (bp *Blueprint) Tables() []Table {
	return lo.Map(bp.tables, func(t *table, _ int) Table { return t })
}

func (bp *Blueprint) Table(id string) Table {
	for _, t := range bp.tables {
		if t.id == id {
			return t
		}
	}
	return nil
}

func (bp *Blueprint) TableIds() []string {
	return lo.Map(bp.tables, func(t *table, _ int) string { return t.id })
}

func (bp *Blueprint) Queries() []Query {
	return lo.Map(bp.queries, func(q *query, _ int) Query { return q })
}

func (bp *Blueprint) Query(id string) Query {
	for _, q := range bp.queries {
		if q.id == id {
			return q
		}
	}
	return nil
}

func (bp *Blueprint) QueryIds() []string {
	return lo.Map(bp.queries, func(q *query, _ int) string { return q.id })
}

func (bp *Blueprint) Modules() []Module {
	return lo.Map(bp.modules, func(m *module, _ int) Module { return m })
}

func (bp *Blueprint) Module(id string) Module {
	for _, m := range bp.modules {
		if m.id == id {
			return m
		}
	}
	return nil
}

func (bp *Blueprint) ModuleIds() []string {
	return lo.Map(bp.modules, func(m *module, _ int) string { return m.id })
}

func (bp *Blueprint) has(id string) bool {
	if bp.tableById(id) != nil {
		return true
	}
	return false
}

func (bp *Blueprint) tableById(id string) *table {
	return lo.FindOrElse[*table](bp.tables, nil, func(t *table) bool { return t.id == id })
}

func (bp *Blueprint) addTable(id string, attrs attributes) *table {
	t := &table{id: id, attributes: attrs}
	bp.tables = append(bp.tables, t)
	return t
}

func (bp *Blueprint) addQuery(id string, q string, qft *queryFuncType) *query {
	q1 := &query{
		id:     id,
		q:      q,
		params: qft.params,
		result: qft.result,
	}
	bp.queries = append(bp.queries, q1)
	return q1
}

func (bp *Blueprint) addModule(id string) *module {
	m1 := &module{
		id: id,
	}
	bp.modules = append(bp.modules, m1)
	return m1
}
