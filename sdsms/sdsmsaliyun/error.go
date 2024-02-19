package sdsmsaliyun

import (
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gaorx/stardust5/sderr"
)

func newErrorForSendSms(resp *dysmsapi.SendSmsResponse, err error) error {
	if err != nil {
		return sderr.Newf("send sms error(%s)", err.Error())
	}
	if resp == nil {
		return sderr.New("send sms error(resp=nil)")
	}
	if tea.Int32Value(resp.StatusCode) != 200 {
		return sderr.Newf("send sms error(status=%d)", tea.Int32Value(resp.StatusCode))
	}
	body := resp.Body
	if body == nil {
		return sderr.New("send sms error(body=nil)")
	}
	if tea.StringValue(body.Code) != "OK" {
		return sderr.Newf(
			"send sms error(code=%s, msg=%s, bizId=%s, reqId=%s)",
			tea.StringValue(body.Code),
			tea.StringValue(body.Message),
			tea.StringValue(body.BizId),
			tea.StringValue(body.RequestId),
		)
	}
	return nil
}

func newErrorForSendBatchSms(resp *dysmsapi.SendBatchSmsResponse, err error) error {
	if err != nil {
		return sderr.Newf("send batch sms error(%s)", err.Error())
	}
	if resp == nil {
		return sderr.New("send batch sms error(resp=nil)")
	}
	if tea.Int32Value(resp.StatusCode) != 200 {
		return sderr.Newf("send batch sms error(status=%d)", tea.Int32Value(resp.StatusCode))
	}
	body := resp.Body
	if body == nil {
		return sderr.New("send batch sms error(body=nil)")
	}
	if tea.StringValue(body.Code) != "OK" {
		return sderr.Newf(
			"send batch sms error(code=%s, msg=%s, bizId=%s, reqId=%s)",
			tea.StringValue(body.Code),
			tea.StringValue(body.Message),
			tea.StringValue(body.BizId),
			tea.StringValue(body.RequestId),
		)
	}
	return nil
}
