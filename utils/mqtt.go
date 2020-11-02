package utils

import (
	"github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	MC mqtt.Client
}

func NewMQClient(broker, user, passwd string) (err error, client MQTTClient) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetUsername(user)
	opts.SetPassword(passwd)

	client.MC = mqtt.NewClient(opts)
	if token := client.MC.Connect(); token.Wait() && token.Error() != nil {
		return token.Error(), client
	}
	return
}

func (mc *MQTTClient) Subscribe(topic string, subCallBackFunc func(client mqtt.Client, msg mqtt.Message)) {
	mc.MC.Subscribe(topic, 0x00, subCallBackFunc)
}

func (mc *MQTTClient) Publish(topic string, content string) {
	mc.MC.Publish(topic, 0x00, true, content)
}
