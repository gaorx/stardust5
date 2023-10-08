package sdblueprint

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdreflect"
	"github.com/samber/lo"
	"strings"
)

// 生成一个mapper，用于将一个json内容映射到一个MarkAsTable的结构体上

func MustMapper[PROTO any]() func(sdjson.Object) (PROTO, error) {
	return lo.Must(NewMapper[PROTO]())
}

func NewMapper[PROTO any]() (func(sdjson.Object) (PROTO, error), error) {
	var proto PROTO
	t, err := scanProtoToTable(proto)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return func(row sdjson.Object) (PROTO, error) {
		row1 := sdjson.Object{}
		for _, c := range t.columns {
			dbCol := c.NameForDB()
			jsonCol := c.id
			if jsonTag := c.Get("json").AsStr(); jsonTag != "" {
				l := strings.SplitN(jsonTag, ",", 2)
				jsonCol = strings.TrimSpace(l[0])
			}
			row1[jsonCol] = row.Get(dbCol).Interface()
		}
		return sdjson.ObjectToStruct[PROTO](row1.ToPrimitive())
	}, nil
}

func scanProtoToTable(proto any) (*table, error) {
	bp := &Blueprint{
		finalized:        false,
		disableMethod:    true,
		disableDummyData: true,
	}
	st, ok := structTypeOf(proto)
	if !ok {
		return nil, sderr.New("prototype is not struct")
	}
	sv := sdreflect.ValueOf(proto)
	if err := scanProto(bp, sv, st, ""); err != nil {
		return nil, sderr.WithStack(err)
	}
	if len(bp.tables) <= 0 {
		return nil, sderr.New("scan struct to table error")
	}
	return bp.tables[0], nil
}
