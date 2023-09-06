package sdblueprint

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdreflect"
)

// 生成一个mapper，用于将一个json内容映射到一个MarkAsTable的结构体上

func MustMapper[PROTO any]() func(sdjson.Object) (PROTO, error) {
	mapper, err := NewMapper[PROTO]()
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return mapper
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
			dbCol := c.Get("db").AsStr()
			if dbCol == "" {
				dbCol = c.id
			}
			row1[c.id] = row.Get(dbCol)
		}
		return sdjson.ObjectToStruct[PROTO](row1)
	}, nil
}

func scanProtoToTable(proto any) (*table, error) {
	bp := &Blueprint{
		finalized:        false,
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
