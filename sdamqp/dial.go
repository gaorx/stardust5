package sdamqp

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/streadway/amqp"
)

type Address struct {
	Url string `json:"url" toml:"url"`
}

func DialConn(addr Address) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr.Url)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return conn, nil
}

func DialChan(addr Address) (*ChannelConn, error) {
	conn, err := DialConn(addr)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return &ChannelConn{Chan: channel, Conn: conn}, nil
}
