package sdblueprint

import (
	"github.com/gaorx/stardust5/sdreflect"
	"reflect"
)

type (
	MarkAsTable  int
	MarkAsQuery  int
	MarkAsModule int
)

type (
	MarkAsPrimaryKey  int
	MarkAsUniqueIndex int
	MarkAsSimpleIndex int
	MarkAsForeignKey  int
)

type (
	MarkAsInlineQuery int
)

type (
	MarkAsGenerateSkeleton  int
	MarkAsGenerateGormModel int
	MarkAsGenerateMysqlDDL  int
)

type markSet []reflect.Type

var (
	markAsTable             = sdreflect.T[MarkAsTable]()
	markAsQuery             = sdreflect.T[MarkAsQuery]()
	markAsModule            = sdreflect.T[MarkAsModule]()
	markAsPrimaryKey        = sdreflect.T[MarkAsPrimaryKey]()
	markAsUniqueIndex       = sdreflect.T[MarkAsUniqueIndex]()
	markAsSimpleIndex       = sdreflect.T[MarkAsSimpleIndex]()
	markAsForeignKey        = sdreflect.T[MarkAsForeignKey]()
	markAsInlineQuery       = sdreflect.T[MarkAsInlineQuery]()
	markAsGenerateSkeleton  = sdreflect.T[MarkAsGenerateSkeleton]()
	markAsGenerateGormModel = sdreflect.T[MarkAsGenerateGormModel]()
	markAsGenerateMysqlDDL  = sdreflect.T[MarkAsGenerateMysqlDDL]()

	// struct mark
	structMarks = markSet{
		markAsTable,
		markAsQuery,
		markAsModule,
	}

	indexMarks = markSet{
		markAsPrimaryKey,
		markAsUniqueIndex,
		markAsSimpleIndex,
		markAsForeignKey,
	}

	inlineMarks = markSet{
		markAsInlineQuery,
	}

	moduleTaskMarks = markSet{
		markAsGenerateSkeleton,
		markAsGenerateGormModel,
		markAsGenerateMysqlDDL,
	}

	// all
	allMarks = func() markSet {
		all := markSet{}
		all = append(all, structMarks...)
		all = append(all, indexMarks...)
		all = append(all, inlineMarks...)
		all = append(all, moduleTaskMarks...)
		return all
	}()
)

func (set markSet) has(mark reflect.Type) bool {
	for _, mark0 := range set {
		if mark == mark0 {
			return true
		}
	}
	return false
}

func getFieldMark(ft reflect.Type, markSet markSet) (reflect.Type, bool) {
	if ft.Kind() == reflect.Ptr {
		return getFieldMark(ft.Elem(), markSet)
	} else {
		if markSet.has(ft) {
			return ft, true
		} else {
			return nil, false
		}
	}
}

func isFieldMark(ft reflect.Type, marSet markSet) bool {
	_, ok := getFieldMark(ft, marSet)
	return ok
}
