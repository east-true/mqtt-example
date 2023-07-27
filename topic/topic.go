package topic

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type TopicDelegate interface {
	GetTopic() string
	MessagePrintHandler() mqtt.MessageHandler
	MessageDeliveryHandler(out chan<- *TopicValue) mqtt.MessageHandler
}

type TopicValue struct {
	Topic string
	Name  string
	Time  time.Time
	Value interface{}
}

type Topic struct {
	Topic string
	Qos   int
}

func (topic *Topic) GetTopic() string {
	return topic.Topic
}

func (topic *Topic) Subscribe(c mqtt.Client, handler mqtt.MessageHandler) mqtt.Token {
	return c.Subscribe(topic.Topic, byte(topic.Qos), handler)
}

func (topic *Topic) Unsubscribe(c mqtt.Client) mqtt.Token {
	return c.Unsubscribe(topic.Topic)
}

func (topic *Topic) MessagePrintHandler() mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		now := time.Now()
		payload := msg.Payload()
		mqtt.DEBUG.Printf("%v %v %v", msg.Topic(), now, string(payload))
	}
}

func (topic *Topic) MessageDeliveryHandler(out chan<- *TopicValue) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		now := time.Now()
		out <- &TopicValue{
			Topic: msg.Topic(),
			Time:  now,
			Value: msg.Payload(),
		}
	}
}
