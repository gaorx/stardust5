package sdecho

import (
	"github.com/gaorx/stardust5/sdjson"
	"slices"
)

type RequestPage struct {
	Page     int
	PageSize int
}

func (page RequestPage) Tuple() (int, int) {
	return page.Page, page.PageSize
}

type BaseRequest interface {
	Page() RequestPage
	SetPage(page RequestPage)
	Fields() []string
	SetFields(fields []string)
	Flags() []string
	SetFlags(flags []string)
}

var _ BaseRequest = &AntdJsonRequest{}

type AntdJsonRequest struct {
	Params       sdjson.Object `json:"params,omitempty"`
	Sort         sdjson.Object `json:"sort,omitempty"`
	Filter       sdjson.Object `json:"filter,omitempty"`
	UpdateFields []string      `json:"fields,omitempty"`
	ShowFlags    []string      `json:"flags,omitempty"`
}

func (req *AntdJsonRequest) Page() RequestPage {
	page := req.Params.TryGet("current", "page", "pageNum", "page_num").AsIntDef(0)
	pageSize := req.Params.TryGet("pageSize", "size", "page_size").AsIntDef(0)
	return RequestPage{Page: page, PageSize: pageSize}
}

func (req *AntdJsonRequest) SetPage(page RequestPage) {
	if req.Params == nil {
		req.Params = sdjson.Object{}
	}
	req.Params.Set("current", page.Page)
	req.Params.Set("pageSize", page.Page)
}

func (req *AntdJsonRequest) Fields() []string {
	return req.UpdateFields
}

func (req *AntdJsonRequest) SetFields(fields []string) {
	req.UpdateFields = slices.Clone(fields)
}

func (req *AntdJsonRequest) Flags() []string {
	return req.ShowFlags
}

func (req *AntdJsonRequest) SetFlags(flags []string) {
	req.ShowFlags = slices.Clone(flags)
}

func (req *AntdJsonRequest) Keyword() string {
	return req.Params.TryGet("keyword", "keyWord").AsStringDef("")
}

func (req *AntdJsonRequest) PageDef(defaultSize int) RequestPage {
	reqPage := req.Page()
	if reqPage.Page <= 0 {
		reqPage.Page = 1
	}
	if reqPage.PageSize <= 0 {
		reqPage.PageSize = defaultSize
	}
	return reqPage
}
