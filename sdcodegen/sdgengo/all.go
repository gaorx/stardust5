package sdgengo

import (
	"fmt"
	"github.com/gaorx/stardust5/sdcodegen"
	"github.com/gaorx/stardust5/sdstrings"
	"github.com/samber/lo"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

func Header(w sdcodegen.Writer, filename string, pkg string, importPackages []string) sdcodegen.Writer {
	if pkg == "" {
		pkg = PackageByFilename(w.Filename())
	}
	w.FL("package %s", pkg)
	w.NL()
	w.L("// AUTO GENERATED, DO NOT EDIT")
	w.L("// AUTO GENERATED, DO NOT EDIT")
	w.L("// AUTO GENERATED, DO NOT EDIT")
	w.NL()
	w.WritePlaceholder(&sdcodegen.Placeholder{
		Name: "go_imports",
		Data: slices.Clone(importPackages),
		Expand: func(w sdcodegen.Writer, data any) {
			importPackages1 := data.([]string)
			if len(importPackages1) > 0 {
				sort.Strings(importPackages1)
				w.L("import (")
				for _, importPkg := range importPackages1 {
					if importPkg == "" {
						continue
					}
					if strings.Contains(importPkg, `"`) {
						w.I(1).L(importPkg)
					} else {
						w.I(1).FL(`"%s"`, importPkg)
					}
				}
				w.L(")")
			}
		},
	})
	return w
}

func AddImportPackages(w sdcodegen.Writer, pkgs []string) sdcodegen.Writer {
	if len(pkgs) <= 0 {
		return w
	}
	return w.UsePlaceholder("go_imports", func(p *sdcodegen.Placeholder) {
		importPackages := p.Data.([]string)
		for _, pkg := range pkgs {
			if !lo.Contains(importPackages, pkg) {
				importPackages = append(importPackages, pkg)
			}
		}
		p.Data = importPackages
	})
}

type NamedType struct {
	Name, Type string
}

func Return(types ...string) []NamedType {
	if len(types) <= 0 {
		return nil
	}
	return lo.Map(types, func(typ string, _ int) NamedType { return NamedType{Type: typ} })
}

func Func(
	w sdcodegen.Writer,
	name string,
	params []NamedType,
	returns []NamedType,
	body func(sdcodegen.Writer),
) sdcodegen.Writer {
	w.FL("func %s(%s) %s {", name, getInStr(params), getOutStr(returns))
	if body != nil {
		body(w)
	}
	return w.L("}")
}

func Method(
	w sdcodegen.Writer,
	name string,
	self NamedType,
	params []NamedType,
	returns []NamedType,
	body func(sdcodegen.Writer),
) sdcodegen.Writer {
	w.FL("func (%s) %s(%s) %s {", self.String(), name, getInStr(params), getOutStr(returns))
	if body != nil {
		body(w)
	}
	return w.L("}")
}

type FieldTag struct {
	K, V string
}

type Field struct {
	Name, Type string
	Tags       []FieldTag
	Comment    string
}

func Struct(
	w sdcodegen.Writer,
	name string,
	fields []Field,
) sdcodegen.Writer {
	if name != "" {
		w.F("type %s ", name)
	}
	if len(fields) > 0 {
		w.L("struct {")
		for _, field := range fields {
			w.I(1).L(field.String())
		}
		w.L("}")
	} else {
		w.L("struct {}")
	}
	return w
}

func (nt NamedType) String() string {
	if nt.Name != "" && nt.Type != "" {
		return nt.Name + " " + nt.Type
	} else if nt.Name != "" {
		return nt.Name
	} else if nt.Type != "" {
		return nt.Type
	} else {
		return ""
	}
}

func (f Field) String() string {
	s := NamedType{Name: f.Name, Type: f.Type}.String()
	if len(f.Tags) > 0 {
		tag := sdstrings.JoinFunc(f.Tags, " ", func(ft FieldTag, _ int) string { return ft.String() })
		if tag != "" {
			s += fmt.Sprintf("`%s`", tag)
		}
	}
	if f.Comment != "" {
		s += fmt.Sprintf("// %s", f.Comment)
	}
	return s
}

func (ft FieldTag) String() string {
	if ft.K == "" {
		return ft.V
	}
	return fmt.Sprintf(`%s:"%s"`, ft.K, ft.V)
}

func PackageByFilename(fn string) string {
	d := filepath.Dir(fn)
	if d == "" {
		return ""
	}
	_, p := filepath.Split(d)
	return p
}

func getInStr(params []NamedType) string {
	if len(params) <= 0 {
		return ""
	}
	return sdstrings.JoinFunc(params, ", ", func(param NamedType, _ int) string {
		return param.String()
	})
}

func getOutStr(returns []NamedType) string {
	if len(returns) <= 0 {
		return ""
	} else if len(returns) == 1 {
		return returns[0].String()
	} else {
		return fmt.Sprintf("(%s)", sdstrings.JoinFunc(returns, ", ", func(ret NamedType, _ int) string {
			return ret.String()
		}))
	}
}
