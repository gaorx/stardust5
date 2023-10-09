package sdgorm

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/samber/lo"
	"strings"
)

func MustMapper[M any](asPrimitive bool) func(sdjson.Object) (M, error) {
	return lo.Must(NewMapper[M](asPrimitive))
}

func NewMapper[M any](asPrimitive bool) (func(sdjson.Object) (M, error), error) {
	var model M
	s, err := ParseSchema(model, nil)
	if err != nil {
		return nil, sderr.Wrap(err, "parse model error")
	}
	return func(row sdjson.Object) (M, error) {
		row1 := sdjson.Object{}
		for _, f := range s.Fields {
			jsonCol := f.Name
			if jsonTag := f.StructField.Tag.Get("json"); jsonTag != "" {
				l := strings.SplitN(jsonTag, ",", 2)
				jsonCol = strings.TrimSpace(l[0])
			}
			row1[jsonCol] = row.Get(f.DBName).Interface()
		}
		if asPrimitive {
			return sdjson.ObjectToStruct[M](row1.TryPrimitive())
		} else {
			return sdjson.ObjectToStruct[M](row1)
		}
	}, nil
}

func MustMapperToAny[M any](asPrimitive bool) func(sdjson.Object) (any, error) {
	return lo.Must(NewMapperToAny[M](asPrimitive))
}

func NewMapperToAny[M any](asPrimitive bool) (func(sdjson.Object) (any, error), error) {
	mapper0, err := NewMapper[M](asPrimitive)
	if err != nil {
		return nil, err
	}
	return func(row sdjson.Object) (any, error) {
		if row1, err := mapper0(row); err != nil {
			return nil, err
		} else {
			return row1, nil
		}
	}, nil
}
