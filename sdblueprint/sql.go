package sdblueprint

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdtemplate"
)

const selfTableId = "T"

type sqlTemplate struct {
	q string
}

func newSqlTemplate(q string) sqlTemplate {
	return sqlTemplate{q}
}

func (qt sqlTemplate) expand(bp *Blueprint, tableId string) (string, error) {
	data := qt.getData(bp, tableId)
	r, err := sdtemplate.Text.Exec(qt.q, data)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	return r, nil
}

func (qt sqlTemplate) getData(bp *Blueprint, tableId string) map[string]string {
	data := map[string]string{}

	addTableName := func(t, name string) {
		data[t] = name
	}

	addSelfColumnName := func(c, name string) {
		data[c] = name
	}

	addRefColumnName := func(t, c, name string) {
		data[t+"_"+c] = name
	}

	for _, t := range bp.tables {
		addTableName(t.id, t.NameForDB())
		for _, c := range t.columns {
			addRefColumnName(t.id, c.id, c.NameForDB())
		}
		if tableId != "" && tableId == t.id {
			addTableName(selfTableId, t.NameForDB())
			for _, c := range t.columns {
				addSelfColumnName(c.id, c.NameForDB())
				addRefColumnName(selfTableId, c.id, c.NameForDB())
			}
		}
	}
	return data
}
