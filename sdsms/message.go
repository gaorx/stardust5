package sdsms

import (
	"encoding/json"
	"github.com/samber/lo"
)

type Message struct {
	Phone    string
	SignName string
	Param    any
	Content  string
}

type Messages []Message

func (msg Message) ParamToMap() (map[string]string, bool) {
	param := msg.Param
	if param == nil {
		return nil, true
	} else if p, ok := param.(map[string]string); ok {
		return p, true
	} else {
		j, err := json.Marshal(param)
		if err != nil {

			return nil, false
		}
		var paramMap map[string]string
		err = json.Unmarshal(j, &paramMap)
		if err != nil {
			return nil, false
		}
		return paramMap, true
	}
}

func (msgs Messages) SetSignName(signName string) {
	for i := 0; i < len(msgs); i++ {
		msgs[i].SignName = signName
	}
}

func (msgs Messages) Phones() []string {
	return lo.Map(msgs, func(msg Message, _ int) string {
		return msg.Phone
	})
}

func (msgs Messages) SignNames() []string {
	return lo.Map(msgs, func(msg Message, _ int) string {
		return msg.SignName
	})
}

func (msgs Messages) Params() []any {
	return lo.Map(msgs, func(msg Message, _ int) any {
		return msg.Param
	})
}

func (msgs Messages) ParamsToMap() ([]map[string]string, bool) {
	var params []map[string]string
	for _, msg := range msgs {
		param, ok := msg.ParamToMap()
		if !ok {
			return nil, false
		}
		params = append(params, param)
	}
	return params, true
}
