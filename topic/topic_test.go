package topic_test

import (
	"mqttgo/broker"
	"mqttgo/broker/opt"
	"mqttgo/topic"
	"testing"
)

func TestTopicSubscribe(t *testing.T) {
	opt := new(opt.Option)
	broker := broker.New("127.0.0.1:503", nil)
	if err := broker.Connect(opt.None()); err != nil {
		t.Error(err)
	}
	defer broker.Disconnect()

	topic := &topic.Topic{
		Topic: "/machbase/test/#",
		Qos:   1,
	}
	token := topic.Subscribe(broker.Client, topic.MessagePrintHandler())
	_ = token.Wait()
	if token.Error() != nil {
		t.Error(token.Error())
	}

	token = topic.Unsubscribe(broker.Client)
	_ = token.Wait()
	if token.Error() != nil {
		t.Error(token.Error())
	}
}

func TestTopicPublish(t *testing.T) {
	opt := new(opt.Option)
	broker := broker.New("127.0.0.1:503", nil)
	if err := broker.Connect(opt.None()); err != nil {
		t.Error(err)
	}
	defer broker.Disconnect()

	payload := []byte("HI")
	topic := &topic.Topic{
		Topic: "/machbase/test/#",
		Qos:   1,
	}
	token := topic.Publish(broker.Client, payload)
	_ = token.Wait()
	if token.Error() != nil {
		t.Error(token.Error())
	}
}
