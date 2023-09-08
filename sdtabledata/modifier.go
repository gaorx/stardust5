package sdtabledata

import (
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdstrings"
)

type Modifier interface {
	ModifyRow(sdjson.Object) sdjson.Object
}

var (
	_ Modifier = ModifierFunc(nil)
	_ Modifier = Modifiers{}
)

type ModifierFunc func(sdjson.Object) sdjson.Object
type Modifiers []Modifier

func (f ModifierFunc) ModifyRow(row sdjson.Object) sdjson.Object {
	return f(row)
}

func (modifiers Modifiers) ModifyRow(row sdjson.Object) sdjson.Object {
	if len(modifiers) <= 0 {
		return row
	}
	row1 := row
	for _, m := range modifiers {
		if m != nil {
			row1 = m.ModifyRow(row1)
		}
	}
	return row1
}

func SetColumn(col string, v any) ModifierFunc {
	return func(row sdjson.Object) sdjson.Object {
		if f, ok := v.(func(any) any); ok {
			row[col] = f(row[col])
		} else {
			row[col] = v
		}
		return row
	}
}

func SetColumns(colVals map[string]any) ModifierFunc {
	if len(colVals) <= 0 {
		return nil
	}
	return func(row sdjson.Object) sdjson.Object {
		for col, v := range colVals {
			if f, ok := v.(func(any) any); ok {
				row[col] = f(row[col])
			} else {
				row[col] = v
			}
		}
		return row
	}
}

func ExpandColumnShellLikeV(col string, vars map[string]string) ModifierFunc {
	return func(row sdjson.Object) sdjson.Object {
		row[col] = sdstrings.ExpandShellLikeV(row[col].(string), vars)
		return row
	}
}

func ExpandColumnsShellLikeV(cols []string, vars map[string]string) ModifierFunc {
	return func(row sdjson.Object) sdjson.Object {
		for _, col := range cols {
			row[col] = sdstrings.ExpandShellLikeV(row[col].(string), vars)
		}
		return row
	}
}
