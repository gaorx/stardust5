package sdblueprint

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/samber/lo"
	"strings"
)

type Query interface {
	Id() string
	Comment() string
	Group() string
	Jsonable
	Attributes
	Kind() QueryKind
	SQL() string
	Table() string
	Params() []QueryParam
	ParamByName(name string) QueryParam
	ParamByIndex(index int) QueryParam
	Result() QueryResult
	ReturnRowAffected() bool
}

type QueryParam interface {
	Name() string
	Type() string
	Table() string
	IsList() bool
	Jsonable
}

type QueryResult interface {
	Type() string
	Table() string
	IsList() bool
	Jsonable
}

type QueryKind string

const (
	QueryForCreate  = QueryKind("CREATE")
	QueryForUpdate  = QueryKind("UPDATE")
	QueryForExec    = QueryKind("EXEC")
	QueryForRecord  = QueryKind("RECORD")
	QueryForRecords = QueryKind("RECORDS")
	QueryForScalar  = QueryKind("SCALAR")
)

type query struct {
	id      string
	comment string
	group   string
	attributes
	q       string
	tableId string
	params  []*queryParam
	result  *queryResult
}

type queryParam struct {
	name   string
	typ    string
	table  string
	isList bool
}

type queryResult struct {
	typ    string
	table  string
	isList bool
}

var (
	_ Query       = &query{}
	_ QueryParam  = &queryParam{}
	_ QueryResult = &queryResult{}
)

// query

func (q *query) Id() string {
	return q.id
}

func (q *query) Comment() string {
	return q.comment
}

func (q *query) Group() string {
	return q.group
}

func (q *query) Kind() QueryKind {
	sqlUpper := strings.ToUpper(q.q)
	if sqlUpper == "CREATE" || sqlUpper == "NEW" {
		return QueryForCreate
	} else if sqlUpper == "UPDATE" || sqlUpper == "SAVE" {
		return QueryForUpdate
	} else if q.result == nil {
		return QueryForExec
	} else {
		if q.result.table != "" {
			if q.result.isList {
				return QueryForRecords
			} else {
				return QueryForRecord
			}
		} else {
			return QueryForScalar
		}
	}
}

func (q *query) SQL() string {
	return q.q
}

func (q *query) Table() string {
	if q.tableId != "" {
		return q.tableId
	} else {
		if q.result != nil {
			return q.result.table
		} else {
			return ""
		}
	}
}

func (q *query) Params() []QueryParam {
	return lo.Map(q.params, func(param *queryParam, _ int) QueryParam { return param })
}

func (q *query) ParamByName(name string) QueryParam {
	for _, param := range q.params {
		if param.name == name {
			return param
		}
	}
	return nil
}

func (q *query) ParamByIndex(index int) QueryParam {
	if index < 0 || index >= len(q.params) {
		return nil
	}
	return q.params[index]
}

func (q *query) Result() QueryResult {
	return q.result
}

func (q *query) ReturnRowAffected() bool {
	return q.Get("return_row_affected").AsBool(false)
}

func (q *query) ToJsonObject() sdjson.Object {
	if q == nil {
		return nil
	}
	return sdjson.Object{
		"id":      q.id,
		"comment": q.comment,
		"group":   q.group,
		"sql":     q.q,
		"table":   q.tableId,
		"params":  lo.Map(q.params, func(param *queryParam, _ int) sdjson.Object { return param.ToJsonObject() }),
		"result":  q.result.ToJsonObject(),
	}
}

func (q *query) setComment(comment string) *query {
	q.comment = comment
	return q
}

func (q *query) setGroup(group string) *query {
	q.group = group
	return q
}

func (q *query) setTable(tableId string) *query {
	q.tableId = tableId
	return q
}

// query param

func (param *queryParam) Name() string {
	if param == nil {
		return ""
	}
	return param.name
}

func (param *queryParam) Type() string {
	if param == nil {
		return ""
	}
	return param.typ
}

func (param *queryParam) Table() string {
	if param == nil {
		return ""
	}
	return param.table
}

func (param *queryParam) IsList() bool {
	if param == nil {
		return false
	}
	return param.isList
}

func (param *queryParam) ToJsonObject() sdjson.Object {
	if param == nil {
		return nil
	}
	return sdjson.Object{
		"name":    param.name,
		"type":    param.typ,
		"table":   param.table,
		"is_list": param.isList,
	}
}

// query result

func (res *queryResult) Type() string {
	if res == nil {
		return ""
	}
	return res.typ
}

func (res *queryResult) Table() string {
	if res == nil {
		return ""
	}
	return res.table
}

func (res *queryResult) IsList() bool {
	if res == nil {
		return false
	}
	return res.isList
}

func (res *queryResult) ToJsonObject() sdjson.Object {
	if res == nil {
		return nil
	}
	return sdjson.Object{
		"type":    res.Type(),
		"table":   res.table,
		"is_list": res.isList,
	}
}

// expand

func expandQuery(bp *Blueprint) error {
	expandTable := func(elemTableId string, tableId string) (string, bool) {
		if elemTableId == selfTableId {
			elemTableId = tableId
		}
		t := bp.Table(elemTableId)
		if t == nil {
			return "", false
		}
		return t.NameForGo(), true
	}

	expandType := func(typ string, tableId string) (string, bool) {
		if tableId == "" {
			return typ, true
		} else {
			if strings.HasPrefix(typ, "[]*") {
				model, ok := expandTable(typ[len("[]*"):], tableId)
				if !ok {
					return "", false
				}
				return "[]*" + model, true
			} else if strings.HasPrefix(typ, "[]") {
				model, ok := expandTable(typ[len("[]"):], tableId)
				if !ok {
					return "", false
				}
				return "[]" + model, true
			} else if strings.HasPrefix(typ, "*") {
				model, ok := expandTable(typ[len("*"):], tableId)
				if !ok {
					return "", false
				}
				return "*" + model, true
			} else {
				model, ok := expandTable(typ, tableId)
				if !ok {
					return "", false
				}
				return model, true
			}
		}
	}

	for _, q := range bp.queries {
		tableId := q.Table()
		if tableId == selfTableId {
			return sderr.NewWith("no table in query", q.id)
		}

		// params
		for _, param := range q.params {
			if param.table == selfTableId {
				param.table = tableId
			}
			if typ1, ok := expandType(param.typ, param.table); !ok {
				return sderr.NewWith("expand param type", sderr.Attrs{"t": param.typ, "q": q.id})
			} else {
				param.typ = typ1
			}
		}

		// result
		if q.result != nil {
			if q.result.table == selfTableId {
				q.result.table = tableId
			}
			if typ1, ok := expandType(q.result.typ, q.result.table); !ok {
				return sderr.NewWith("expand result type", sderr.Attrs{"type": q.result.typ, "q": q.id})
			} else {
				q.result.typ = typ1
			}
		}

		// sql
		expandedSql, err := newSqlTemplate(q.q).expand(bp, tableId)
		if err != nil {
			return sderr.NewWith("expand SQL error in query", q.q)
		}
		q.q = expandedSql
	}
	return nil
}
