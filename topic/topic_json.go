package topic

import (
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type JsonTopic struct {
	Topic string
	Qos   int
	Keys  []string // ["" || tag name || json key]
}

func (topic *JsonTopic) GetTopic() string {
	return topic.Topic
}

func (topic *JsonTopic) MessagePrintHandler() mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		now := time.Now()
		jsonMap := make(map[string]interface{})
		payload := msg.Payload()
		err := json.Unmarshal(payload, &jsonMap)
		if err != nil {
			mqtt.WARN.Printf("%s payload is not json : %v", topic.Topic, err)
			return
		}

		for i := range topic.Keys {
			if value, ok := jsonMap[topic.Keys[i]]; ok {
				res := TopicValue{
					Topic: msg.Topic(),
					Name:  topic.Keys[i],
					Time:  now,
					Value: value,
				}
				mqtt.DEBUG.Printf("%+v", res)
			} else {
				mqtt.ERROR.Println("not found '%d' in json payload", topic.Keys[i])
			}
		}
	}
}

func (topic *JsonTopic) MessageDeliveryHandler(out chan<- *TopicValue) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		now := time.Now()
		jsonMap := make(map[string]interface{})
		payload := msg.Payload()
		err := json.Unmarshal(payload, &jsonMap)
		if err != nil {
			mqtt.WARN.Printf("%s payload is not json : %v", topic.Topic, err)
			return
		}

		for i := range topic.Keys {
			if value, ok := jsonMap[topic.Keys[i]]; ok {
				out <- &TopicValue{
					Topic: msg.Topic(),
					Name:  topic.Keys[i],
					Time:  now,
					Value: value,
				}
			} else {
				mqtt.ERROR.Println("not found '%d' in json payload", topic.Keys[i])
			}
		}
	}
}
