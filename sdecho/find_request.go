package sdecho

import (
	"github.com/gaorx/stardust5/sdjson"
)

type AntdFindReq[PARAMS any, SORT any, FLT any] struct {
	Params PARAMS `json:"params,omitempty"`
	Sort   SORT   `json:"sort,omitempty"`
	Filter FLT    `json:"filter,omitempty"`
	Flags  string `json:"string,omitempty"`
}

type AntdCommonParams struct {
	Page     int    `json:"current,omitempty"`
	PageSize int    `json:"pageSize,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
}

type AntdJsonFindReq struct {
	Params sdjson.Object `json:"params,omitempty"`
	Sort   sdjson.Object `json:"sort,omitempty"`
	Filter sdjson.Object `json:"filter,omitempty"`
	Flags  string        `json:"string,omitempty"`
}

func (req AntdJsonFindReq) Paging1(defaultSize int) (int, int) {
	page := req.Params.TryGet("current", "page", "pageNum", "page_num").AsIntDef(0)
	if page <= 0 {
		page = 1
	}
	pageSize := req.Params.TryGet("pageSize", "size", "page_size").AsIntDef(0)
	if pageSize <= 0 {
		pageSize = defaultSize
	}
	return page, pageSize
}
