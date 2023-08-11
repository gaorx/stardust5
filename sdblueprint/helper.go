package sdblueprint

import (
	"bytes"
	"github.com/gaorx/stardust5/sdcodegen/sdgengo"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/samber/lo"
	"go/ast"
	"go/parser"
	"path/filepath"
	"slices"
	"text/template"
	"unicode"
)

func selectNotEmpty[T comparable](first, second, last T) T {
	var empty T
	if first != empty {
		return first
	}
	if second != empty {
		return second
	}
	return last
}

type Jsonable interface {
	ToJsonObject() sdjson.Object
}

type columnReference struct {
	col      string
	refTable string
	refCol   string
}

type columnReferences []columnReference

func (refs columnReferences) sameTable() (string, bool) {
	if len(refs) <= 0 {
		return "", false
	} else if len(refs) == 1 {
		return refs[0].refTable, true
	} else {
		uniqTables := lo.Uniq(lo.Map(refs, func(ref columnReference, _ int) string { return ref.refTable }))
		if len(uniqTables) > 1 {
			return "", false
		} else {
			return uniqTables[0], true
		}
	}
}

func (refs columnReferences) forFK() ([]string, string, []string, bool) {
	refTable, ok := refs.sameTable()
	if !ok {
		return nil, "", nil, false
	}
	cols := lo.Map(refs, func(ref columnReference, _ int) string { return ref.col })
	refCols := lo.Map(refs, func(ref columnReference, _ int) string { return ref.refCol })
	return cols, refTable, refCols, true
}

type queryFuncType struct {
	params []*queryParam
	result *queryResult
}

func parseQueryFuncType(s string) (*queryFuncType, error) {
	if s == "" {
		return nil, sderr.New("no GO func type in query")
	}

	primaryTypes := []string{
		"bool",
		"string",
		"byte",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
	}

	getText := func(node ast.Node) string {
		return s[node.Pos()-1 : node.End()-1]
	}

	var getTableId func(node ast.Node) string
	var getTypeText func(node ast.Node) string

	getTableId = func(node ast.Node) string {
		if ident, ok := node.(*ast.Ident); ok {
			if lo.Contains(primaryTypes, ident.Name) {
				return ""
			} else {
				return ident.Name
			}
		} else if a, ok := node.(*ast.ArrayType); ok {
			return getTableId(a.Elt)
		} else if s, ok := node.(*ast.StarExpr); ok {
			return getTableId(s.X)
		} else {
			return ""
		}
	}

	getTypeText = func(node ast.Node) string {
		if ident, ok := node.(*ast.Ident); ok {
			return ident.Name
		} else if a, ok := node.(*ast.ArrayType); ok {
			return "[]" + getTypeText(a.Elt)
		} else if s, ok := node.(*ast.StarExpr); ok {
			return "*" + getTypeText(s.X)
		} else {
			return getText(node)
		}
	}

	isList := func(node ast.Node) bool {
		if _, ok := node.(*ast.ArrayType); ok {
			return true
		} else {
			return false
		}
	}

	root, err := parser.ParseExpr(s)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	var queryParams []*queryParam
	funcTypeNode, ok := root.(*ast.FuncType)
	if !ok {
		return nil, sderr.New("query func type must be a GO function type")
	}
	for _, paramNode := range funcTypeNode.Params.List {
		for _, subNode := range paramNode.Names {
			queryParams = append(queryParams, &queryParam{
				name:   subNode.Name,
				typ:    getText(paramNode.Type),
				table:  getTableId(paramNode.Type),
				isList: isList(paramNode.Type),
			})
		}
	}
	if funcTypeNode.Results == nil || len(funcTypeNode.Results.List) <= 0 {
		return &queryFuncType{
			params: queryParams,
			result: nil,
		}, nil
	} else {
		if len(funcTypeNode.Results.List) != 1 {
			return nil, sderr.New("query func must be one result")
		}
		resTypeNode := funcTypeNode.Results.List[0].Type
		return &queryFuncType{
			params: queryParams,
			result: &queryResult{
				typ:    getTypeText(resTypeNode),
				table:  getTableId(resTypeNode),
				isList: getTableId(resTypeNode) != "" && isList(resTypeNode),
			},
		}, nil
	}
}

type idSet map[string]struct{}

func newIdSet() idSet {
	return idSet{}
}

func (ids idSet) add(id string) bool {
	_, ok := ids[id]
	if ok {
		return false
	}
	ids[id] = struct{}{}
	return true
}

func (ids idSet) has(id string) bool {
	_, ok := ids[id]
	return ok
}

func matchIds(ids, patterns []string) []string {
	isMatched := func(id, patt string) bool {
		matched, err := filepath.Match(patt, id)
		if err != nil {
			return false
		}
		return matched
	}

	if len(ids) <= 0 {
		return nil
	}
	if len(patterns) <= 0 {
		return slices.Clone(ids)
	}
	var matchedIds []string
	for _, patt := range patterns {
		for _, id := range ids {
			if isMatched(id, patt) {
				if !lo.Contains(matchedIds, id) {
					matchedIds = append(matchedIds, id)
				}
			}
		}
	}
	return matchedIds
}

var templateBuiltinFuncs = template.FuncMap{
	"toSnakeLower":        sdstrings.ToSnakeL,
	"toSnakeUpper":        sdstrings.ToSnakeU,
	"toCamelLower":        sdstrings.ToCamelL,
	"toCamelUpper":        sdstrings.ToCamelU,
	"goPackageByFilename": sdgengo.PackageByFilename,
}

func executeTemplate(tpl string, data any) (string, error) {
	t, err := template.New("").Funcs(templateBuiltinFuncs).Parse(tpl)
	if err != nil {
		return "", sderr.Wrap(err, "parse text template error")
	}
	buff := bytes.NewBufferString("")
	err = t.Execute(buff, data)
	if err != nil {
		return "", sderr.Wrap(err, "execute text template error")
	}
	return buff.String(), nil
}

func isPublic(id string) bool {
	if id == "" {
		return false
	}
	idUnicode := []rune(id)
	if len(idUnicode) <= 0 {
		return false
	}
	return unicode.IsUpper(idUnicode[0])
}
