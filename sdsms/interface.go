package sdsms

import (
	"context"
)

type Interface interface {
	Send(ctx context.Context, req *SendRequest) error
}

type SendRequest struct {
	TemplateId string
	Messages   Messages
}
