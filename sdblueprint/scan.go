package sdblueprint

import (
	"fmt"
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
	"reflect"
	"strings"
)

func scanProto(bp *Blueprint, sv reflect.Value, st structType, name string) error {
	marks := st.findStructMarkIn(structMarks)
	if len(marks) > 1 {
		return sderr.New("the struct is marked only once")
	} else if len(marks) == 1 {
		mark := marks[0]
		switch mark.mark {
		case markAsTable:
			return scanTableProto(bp, sv, st, mark)
		case markAsQuery:
			return scanQueryProto(bp, sv, st, mark)
		case markAsModule:
			return scanModuleProto(bp, sv, st, mark)
		default:
			return sderr.NewWith("unknown struct mark", mark.mark)
		}
	} else {
		return sderr.New("no mark in proto")
	}
}

func scanTableProto(bp *Blueprint, sv reflect.Value, st structType, mark markedField) error {
	id := mark.getId(&st)
	newTable := bp.addTable(id, func() attributes {
		attrs := attributes{}
		mark.tag.toAttrs(attrs, "db", "go")
		return attrs
	}()).setComment(mark.tag.comment()).setGroup(mark.tag.group())
	n := st.NumField()
	for i := 0; i < n; i++ {
		sf := st.Field(i)
		st := structTag(sf.Tag)
		mark1, ok := getFieldMark(sf.Type, allMarks)
		if ok {
			if mark1 == markAsTable {
				// pass
			} else if mark1 == markAsPrimaryKey {
				newTable.addPrimaryKey(st.split("db")).setComment(st.comment())
			} else if mark1 == markAsUniqueIndex {
				newTable.addUniqueIndex(st.split("db")).setComment(st.comment())
			} else if mark1 == markAsSimpleIndex {
				newTable.addSimpleIndex(st.split("db")).setComment(st.comment())
			} else if mark1 == markAsForeignKey {
				cols, refTable, refCols, ok := st.toForeignKey("db")
				if !ok {
					return sderr.NewWith("parse foreign key error", sf.Tag)
				}
				newTable.addForeignKey(cols, refTable, refCols).setComment(st.comment())
			} else if mark1 == markAsInlineQuery {
				if err := addQueryWithField(bp, nil, newFieldMark(sf), newTable.id); err != nil {
					return sderr.WithStack(err)
				}
			} else {
				return sderr.NewWith("illegal mark in table", mark1)
			}
		} else {
			if sf.Name == dummyDataFieldName {
				continue
			}
			f := &field{
				id:      sf.Name,
				comment: sf.Tag.Get("comment"),
				typ:     sf.Type,
				attributes: func() attributes {
					attrs := attributes{}
					structTag(sf.Tag).toAttrs(attrs,
						"json", "xml", "validate", "go", "default", "db_type", "dbtype",
						"pk", "primary_key", "unique", "index",
					)
					structTag(sf.Tag).toAttrsForFlags(attrs, "db")
					return attrs
				}(),
			}
			if !isPublic(f.id) || f.Get("db") == "-" {
				newTable.addMember(f)
			} else {
				newTable.addFieldAsColumn(f)
			}
			if f.First([]string{"db.index", "index"}).AsBool(false) {
				newTable.addSimpleIndex([]string{f.id})
			} else if f.First([]string{"db.unique", "unique"}).AsBool(false) {
				newTable.addUniqueIndex([]string{f.id})
			}
		}
	}
	pkColsByFlag := lo.FilterMap(newTable.columns, func(col column, _ int) (string, bool) {
		if col.IsPrimaryKey() {
			return col.id, true
		}
		return "", false
	})
	newTable.ensurePrimaryKey().addColumns(pkColsByFlag)

	// dummy data
	dummyData, err := getDummyData(sv, st, newTable.columns)
	if err != nil {
		return sderr.WithStack(err)
	}
	newTable.dummyData = dummyData
	return nil
}

func scanQueryProto(bp *Blueprint, _ reflect.Value, st structType, mark markedField) error {
	return addQueryWithField(bp, &st, mark, "")
}

func addQueryWithField(bp *Blueprint, st *structType, mark markedField, inTable string) error {
	id := mark.getId(st)
	q := mark.tag.Get("db")
	tableId := selectNotEmpty(mark.tag.Get("table"), inTable, "")
	funcDecl := mark.tag.Get("go")
	if funcDecl == "" {
		sqlUpper := strings.ToUpper(q)
		if sqlUpper == "CREATE" || sqlUpper == "NEW" || sqlUpper == "UPDATE" || sqlUpper == "SAVE" {
			funcDecl = fmt.Sprintf("func(o *%s)", selfTableId)
		}
	}
	qft, err := parseQueryFuncType(funcDecl)
	if err != nil {
		return sderr.WithStack(err)
	}
	q1 := bp.addQuery(id, q, qft).setComment(mark.tag.comment()).setGroup(mark.tag.group()).setTable(tableId)
	q1.attributes = func() attributes {
		attrs := attributes{}
		mark.tag.toAttrs(attrs, "return_row_affected")
		return attrs
	}()
	return nil
}

func scanModuleProto(bp *Blueprint, _ reflect.Value, st structType, mark markedField) error {
	id := mark.getId(&st)
	newModule := bp.addModule(id).setComment(mark.tag.comment())
	n := st.NumField()
	for i := 0; i < n; i++ {
		sf := st.Field(i)
		st := structTag(sf.Tag)
		mark1, ok := getFieldMark(sf.Type, allMarks)
		attrs := func() attributes {
			attrs0 := attributes{}
			st.toAttrs(attrs0, "template", "dir", "group", "groups", "query_with_context")
			return attrs0
		}()
		if ok {
			if mark1 == markAsModule {
				// pass
			} else if mark1 == markAsGenerateSkeleton {
				newModule.addTask(&ModuleTaskGenerateSkeleton{
					Template: attrs.Get("template").AsStr(),
					Dirname:  attrs.Get("dir").AsStr(),
				})
			} else if mark1 == markAsGenerateGormModel {
				newModule.addTask(&ModuleTaskGenerateGormModel{
					Dirname:          attrs.Get("dir").AsStr(),
					Groups:           attrs.First([]string{"groups", "group"}).AsSlice(","),
					QueryWithContext: attrs.Get("query_with_context").AsBool(false),
				})
			} else if mark1 == markAsGenerateMysqlDDL {
				newModule.addTask(&ModuleTaskGenerateMysqlDDL{
					Dirname: attrs.Get("dir").AsStr(),
					Groups:  attrs.First([]string{"groups", "group"}).AsSlice(","),
				})
			} else {
				return sderr.NewWith("illegal mark in module", mark1)
			}
		} else {
			continue // ignore
		}
	}
	return nil
}
