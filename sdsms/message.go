package sdsms

import (
	"github.com/samber/lo"
	"maps"
)

type Message struct {
	Phone    string
	SignName string
	Param    map[string]string
	Content  string
}

type Messages []Message

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

func (msgs Messages) Params() []map[string]string {
	return lo.Map(msgs, func(msg Message, _ int) map[string]string {
		return maps.Clone(msg.Param)
	})
}
