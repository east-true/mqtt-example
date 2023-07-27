package topic

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ValueTopic struct {
	Topic string
	Qos   int
	Name  string
}

func (topic *ValueTopic) GetTopic() string {
	return topic.Topic
}

func (topic *ValueTopic) MessagePrintHandler() mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		now := time.Now()
		payload := msg.Payload()
		mqtt.DEBUG.Printf("%v %v %v", msg.Topic(), now, string(payload))
	}
}

func (topic *ValueTopic) MessageDeliveryHandler(out chan<- *TopicValue) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		now := time.Now()
		payload := msg.Payload()
		out <- &TopicValue{
			Topic: msg.Topic(),
			Name:  topic.Name,
			Time:  now,
			Value: string(payload),
		}
	}
}
