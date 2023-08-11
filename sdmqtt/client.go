package sdmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gaorx/stardust5/sderr"
)

type Client struct {
	mqtt.Client
}

func Dial(opts *mqtt.ClientOptions) *Client {
	c := mqtt.NewClient(opts)
	return &Client{c}
}

func (c *Client) ConnectSync() error {
	token := c.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.Wrap(err, "connect MQTT sync error")
	}
	return nil
}

func (c *Client) SubscribeSync(topic string, qos byte, callback mqtt.MessageHandler) error {
	token := c.Subscribe(topic, qos, callback)
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.Wrap(err, "subscribe MQTT sync error")
	}
	return nil
}

func (c *Client) UnsubscribeSync(topics ...string) error {
	token := c.Unsubscribe(topics...)
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.Wrap(err, "unsubscribe MQTT sync error")
	}
	return nil
}

func (c *Client) PublishSync(topic string, qos byte, retained bool, payload any) error {
	token := c.Publish(topic, qos, retained, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.Wrap(err, "publish MQTT sync error")
	}
	return nil
}
