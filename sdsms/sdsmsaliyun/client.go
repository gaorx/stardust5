package sdsmsaliyun

import (
	"context"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdjson"
	"github.com/gaorx/stardust5/sdsms"
	"github.com/samber/lo"
)

type Client struct {
	client *dysmsapi.Client
	config *Config
}

type Config struct {
	Endpoint        string `json:"endpoint" toml:"endpoint" yaml:"endpoint"`
	AccessId        string `json:"access_id" toml:"access_id" yaml:"access_id"`
	AccessKey       string `json:"access_key" toml:"access_key" yaml:"access_key"`
	DefaultSignName string `json:"default_sign_name" toml:"default_sign_name" yaml:"default_sign_name"`
}

var _ sdsms.Interface = &Client{}

func New(config *Config) (*Client, error) {
	var config1 = lo.FromPtr(config)
	if config1.Endpoint == "" {
		return nil, sderr.New("no endpoint")
	}
	if config1.AccessId == "" {
		return nil, sderr.New("no access id")
	}
	if config1.AccessKey == "" {
		return nil, sderr.New("no access key")
	}
	if config1.DefaultSignName == "" {
		return nil, sderr.New("no default sign name")
	}
	aliConfig := openapi.Config{
		AccessKeyId:     &config1.AccessId,
		AccessKeySecret: &config1.AccessKey,
		Endpoint:        tea.String(config1.Endpoint),
	}
	client, err := dysmsapi.NewClient(&aliConfig)
	if err != nil {
		return nil, sderr.Wrap(err, "create aliyun sms client")
	}
	return &Client{client: client, config: &config1}, nil
}

func (c *Client) Config() *Config {
	return c.config
}

func (c *Client) Send(ctx context.Context, req *sdsms.SendRequest) error {
	req1 := lo.FromPtr(req)
	switch len(req1.Messages) {
	case 0:
		return nil
	case 1:
		msg := req1.Messages[0]
		if msg.Phone == "" {
			return sderr.New("no phone number")
		}
		if msg.SignName == "" {
			msg.SignName = c.config.DefaultSignName
		}
		paramMap, ok := msg.ParamToMap()
		if !ok {
			return sderr.New("invalid message param")
		}
		aliReq := &dysmsapi.SendSmsRequest{
			PhoneNumbers:  tea.String(msg.Phone),
			SignName:      tea.String(msg.SignName),
			TemplateCode:  tea.String(req1.TemplateId),
			TemplateParam: messageToJson(paramMap),
		}
		aliResp, err := c.client.SendSms(aliReq)
		return newErrorForSendSms(aliResp, err)
	default:
		var msgs1 sdsms.Messages
		for _, msg := range req1.Messages {
			msg1 := msg
			if msg1.Phone == "" {
				return sderr.New("no phone number")
			}
			if msg1.SignName == "" {
				msg1.SignName = c.config.DefaultSignName
			}
			if paramMap, ok := msg.ParamToMap(); !ok {
				return sderr.New("invalid message param")
			} else {
				msg1.Param = paramMap
			}
			msgs1 = append(msgs1, msg1)
		}
		aliReq := &dysmsapi.SendBatchSmsRequest{
			PhoneNumberJson:   tea.String(lo.Must(sdjson.MarshalString(msgs1.Phones()))),
			SignNameJson:      tea.String(lo.Must(sdjson.MarshalString(msgs1.SignNames()))),
			TemplateCode:      tea.String(req1.TemplateId),
			TemplateParamJson: tea.String(lo.Must(sdjson.MarshalString(msgs1.Params()))),
		}
		aliResp, err := c.client.SendBatchSms(aliReq)
		return newErrorForSendBatchSms(aliResp, err)
	}
}

func messageToJson(paramMap map[string]string) *string {
	if len(paramMap) <= 0 {
		return nil
	}
	j := lo.Must(sdjson.MarshalString(paramMap))
	return tea.String(j)
}
