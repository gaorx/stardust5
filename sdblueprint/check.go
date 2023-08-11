package sdblueprint

import (
	"github.com/gaorx/stardust5/sderr"
	"reflect"
)

func (bp *Blueprint) check() error {
	// check tables
	for _, t := range bp.tables {
		if err := t.checkSelf(); err != nil {
			return err
		}
	}

	// check queries
	for _, q := range bp.queries {
		if err := q.checkSelf(); err != nil {
			return err
		}
	}

	// check modules
	for _, m := range bp.modules {
		if err := m.checkSelf(); err != nil {
			return err
		}
	}

	// check same id
	ids := newIdSet()
	for _, t := range bp.tables {
		if !ids.add(t.id) {
			return sderr.NewWith("same id", t.id)
		}
	}
	for _, q := range bp.queries {
		if !ids.add(q.id) {
			return sderr.NewWith("same id", q.id)
		}
	}
	for _, m := range bp.modules {
		if !ids.add(m.id) {
			return sderr.NewWith("same id", m.id)
		}
	}

	// check reference
	for _, t := range bp.tables {
		if err := t.checkRef(bp); err != nil {
			return err
		}
	}
	for _, q := range bp.queries {
		if err := q.checkRef(bp); err != nil {
			return err
		}
	}
	return nil
}

func (bp *Blueprint) hasTable(id string) bool {
	for _, t := range bp.tables {
		if t.id == id {
			return true
		}
	}
	return false
}

func (bp *Blueprint) hasQuery(id string) bool {
	for _, q := range bp.queries {
		if q.id == id {
			return true
		}
	}
	return false
}

func (bp *Blueprint) hasTableColumn(tableId, columnId string) bool {
	for _, t := range bp.tables {
		if t.id == tableId {
			for _, c := range t.columns {
				if c.id == columnId {
					return true
				}
			}
		}
	}
	return false
}

func (t *table) checkSelf() error {
	// check id
	if t.id == "" {
		return sderr.New("no table id")
	}

	// check columns
	for _, c := range t.columns {
		if err := c.checkSelf(t.id); err != nil {
			return err
		}
	}

	// check members
	for _, member := range t.members {
		if err := member.checkSelfForMember(t.id); err != nil {
			return err
		}
	}

	// check indexes
	for _, idx := range t.indexes {
		if err := idx.checkSelf(t.id); err != nil {
			return err
		}
	}
	if t.PrimaryKey() == nil {
		return sderr.NewWith("no primary key in table", t.id)
	}

	// check same id
	ids := newIdSet()
	for _, c := range t.columns {
		if !ids.add(c.id) {
			return sderr.NewWith("same id", c.id)
		}
	}
	for _, member := range t.members {
		if !ids.add(member.id) {
			return sderr.NewWith("same id", member.id)
		}
	}

	return nil
}

func (t *table) checkRef(bp *Blueprint) error {
	// index
	for _, idx := range t.indexes {
		for _, colId := range idx.columns {
			if !bp.hasTableColumn(t.id, colId) {
				return sderr.NewWith("not found column id", sderr.Attrs{"t": t.id, "c": colId})
			}
		}
		if idx.kind == IndexFK {
			if !bp.hasTable(idx.referenceTable) {
				return sderr.NewWith("not found reference table", t.id)
			}
			for _, refColId := range idx.referenceColumns {
				if !bp.hasTableColumn(idx.referenceTable, refColId) {
					return sderr.NewWith("not found reference column", sderr.Attrs{"t": t.id, "ref": refColId})
				}
			}
		}
	}

	// check dummy data
	if len(t.dummyData) > 0 {
		for _, record := range t.dummyData {
			for colId, _ := range record {
				if !bp.hasTableColumn(t.id, colId) {
					return sderr.NewWith("unknown dummy data column", sderr.Attrs{"c": colId, "t": t.id})
				}
			}
		}
	}
	return nil
}

func (q *query) checkSelf() error {
	if q.id == "" {
		return sderr.New("no query id")
	}
	if q.q == "" {
		return sderr.NewWith("no SQL in query", q.id)
	}
	return nil
}

func (m *module) checkSelf() error {
	if m.id == "" {
		return sderr.New("no module id")
	}
	for _, t := range m.tasks {
		switch t1 := t.(type) {
		case *ModuleTaskGenerateSkeleton:
			if t1.Template == "" {
				return sderr.New("no template in generate skeleton task")
			}
			if t1.Dirname == "" {
				return sderr.New("no dir in generate skeleton task")
			}
		case *ModuleTaskGenerateGormModel:
			if t1.Dirname == "" {
				return sderr.New("no dir in generate model task")
			}
		case *ModuleTaskGenerateMysqlDDL:
			if t1.Dirname == "" {
				return sderr.New("no dir in generate sql task")
			}
		default:
			return sderr.NewWith("illegal task", sderr.Attrs{"task": reflect.TypeOf(t1).String(), "model": m.id})
		}
	}
	return nil
}

func (q *query) checkRef(bp *Blueprint) error {
	if q.tableId != "" {
		if !bp.hasTable(q.tableId) {
			return sderr.NewWith("not found table", sderr.Attrs{"t": q.tableId, "q": q.id})
		}
	}

	for _, param := range q.params {
		if param.table != "" && param.table != selfTableId {
			if !bp.hasTable(param.table) {
				return sderr.NewWith("not found table(param)", sderr.Attrs{"t": param.table, "q": q.id})
			}
		}
	}
	if q.result != nil && q.result.table != "" && q.result.table != selfTableId {
		if !bp.hasTable(q.result.table) {
			return sderr.NewWith("not found table(result)", sderr.Attrs{"t": q.result.table, "q": q.id})
		}
	}
	return nil
}

func (c column) checkSelf(tableId string) error {
	if c.id == "" {
		return sderr.NewWith("no column id in table", tableId)
	}
	return nil
}

func (f *field) checkSelfForMember(tableId string) error {
	if f.id == "" {
		return sderr.NewWith("no member id in table", tableId)
	}
	return nil
}

func (idx *index) checkSelf(tableId string) error {
	if len(idx.columns) <= 0 {
		return sderr.NewWith("no column ids for index in table ", tableId)
	}
	if idx.kind == IndexFK {
		if idx.referenceTable == "" {
			return sderr.NewWith("no reference for index in table", tableId)
		}
		if len(idx.referenceColumns) != len(idx.columns) {
			return sderr.NewWith("illegal reference column for index in table", tableId)
		}
	}
	return nil
}
