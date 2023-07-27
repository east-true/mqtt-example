package broker_test

import (
	. "mqttgo/broker"
	"mqttgo/broker/opt"
	"mqttgo/topic"
	"testing"
)

func TestBroker(t *testing.T) {
	topic := &topic.Topic{
		Topic: "/machbase/test/#",
		Qos:   1,
	}
	opt := new(opt.Option)

	broker := New("127.0.0.1:503", topic)
	if err := broker.Connect(opt.None()); err != nil {
		t.Error(err)
	}

	broker.Disconnect()
}
