package sdreq

import (
	"context"
	"github.com/gaorx/stardust5/sderr"
	"github.com/imroc/req/v3"
)

func GetForResponse(ctx context.Context, client *req.Client, url string, opts ...RequestOption) (*req.Response, error) {
	if client == nil {
		client = req.DefaultClient()
	}
	request := applyOptions(client.R().SetContext(ctx), opts)
	response, err := request.Get(url)
	if err != nil {
		return nil, sderr.Wrap(err, "get response error")
	}
	return response, nil
}

func GetForBytes(ctx context.Context, client *req.Client, url string, opts ...RequestOption) (int, []byte, error) {
	response, err := GetForResponse(ctx, client, url, opts...)
	if err != nil {
		return 0, nil, err
	}
	data, err := response.ToBytes()
	if err != nil {
		return response.StatusCode, nil, sderr.Wrap(err, "response to bytes error")
	}
	return response.StatusCode, data, nil
}

func GetForText(ctx context.Context, client *req.Client, url string, opts ...RequestOption) (int, string, error) {
	response, err := GetForResponse(ctx, client, url, opts...)
	if err != nil {
		return 0, "", err
	}
	data, err := response.ToString()
	if err != nil {
		return response.StatusCode, "", sderr.Wrap(err, "response to text error")
	}
	return response.StatusCode, data, nil
}

func GetForJson[R any](ctx context.Context, client *req.Client, url string, opts ...RequestOption) (int, R, error) {
	var r R
	response, err := GetForResponse(ctx, client, url, opts...)
	if err != nil {
		return 0, r, err
	}
	err = response.UnmarshalJson(&r)
	if err != nil {
		return response.StatusCode, r, sderr.Wrap(err, "response unmarshal json error")
	}
	return response.StatusCode, r, nil
}
