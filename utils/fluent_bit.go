package utils

import (
	"github.com/sirupsen/logrus"
)

func PushToFluentBit(url, topic, username, password, message string) (err error) {
	var mc MQTTClient
	if err, mc = NewMQClient(url, username, password); err != nil {
		logrus.Info(err)
		return err
	}
	mc.Publish(topic, message)
	return nil
}
