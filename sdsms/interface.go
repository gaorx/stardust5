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

func Single(templateId string, phone string, signName string, param any) *SendRequest {
	return &SendRequest{
		TemplateId: templateId,
		Messages: Messages{{
			Phone:    phone,
			SignName: signName,
			Param:    param,
		}},
	}
}
