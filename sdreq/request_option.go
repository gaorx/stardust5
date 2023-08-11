package sdreq

import (
	"github.com/imroc/req/v3"
)

type RequestOption func(*req.Request) *req.Request

func applyOptions(request *req.Request, opts []RequestOption) *req.Request {
	for _, opt := range opts {
		request = opt(request)
	}
	return request
}

func QueryParam(k string, v any) RequestOption {
	return func(request *req.Request) *req.Request {
		return request.AddQueryParam(k, ToQueryParam(v))
	}
}

func QueryParams(params map[string]any) RequestOption {
	return func(request *req.Request) *req.Request {
		for k, v := range params {
			request = request.AddQueryParam(k, ToQueryParam(v))
		}
		return request
	}
}

func Header(k, v string) RequestOption {
	return func(request *req.Request) *req.Request {
		return request.SetHeader(k, v)
	}
}

func Headers(headers map[string]string) RequestOption {
	return func(request *req.Request) *req.Request {
		for k, v := range headers {
			return request.SetHeader(k, v)
		}
		return request
	}
}
