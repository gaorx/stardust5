package sdamqp

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/streadway/amqp"
)

type ChannelConn struct {
	Chan *amqp.Channel
	Conn *amqp.Connection
}

func (cc *ChannelConn) Close() error {
	var chanErr, connErr error
	if cc.Chan != nil {
		chanErr = cc.Chan.Close()
		cc.Chan = nil
	}
	if cc.Conn != nil {
		connErr = cc.Conn.Close()
		cc.Conn = nil
	}
	if chanErr != nil {
		return sderr.Wrap(chanErr, "close AMQP channel error")
	} else if connErr != nil {
		return sderr.Wrap(connErr, "close AMQP connection error")
	} else {
		return nil
	}
}
