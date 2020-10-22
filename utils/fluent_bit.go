package utils

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"time"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	logrus.Infof("TOPIC: %s\n", msg.Topic())
	logrus.Infof("MSG: %s\n", msg.Payload())
}

func PushToFluentBit(url string, topic, clientID, message string) error {
	opts := mqtt.NewClientOptions().AddBroker("tcp://0.0.0.0:1883").SetClientID(clientID)
	opts.SetKeepAlive(5 * time.Second)

	// 这里需要注入一个client收到消息后对消息处理的方法
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	// 启动一个链接
	c := mqtt.NewClient(opts)
	defer c.Disconnect(1000)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	token := c.Publish(topic, 0, false, message)
	token.Wait()
	return nil
}
