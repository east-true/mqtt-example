package broker

import (
	"mqttgo/topic"
	"regexp"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Broker struct {
	addr   string // format : [tcp | ssl | ws]://[IP]:[PORT]
	topics []topic.Subscriber
	Client mqtt.Client
}

func New(address string, topics ...topic.Subscriber) *Broker {
	return &Broker{
		addr:   address,
		topics: topics,
	}
}

// opt := broker.Option.get(broker.Addr)
func (broker *Broker) Connect(opt *mqtt.ClientOptions) error {
	broker.Client = mqtt.NewClient(opt)
	token := broker.Client.Connect()
	_ = token.Wait()
	if token.Error() != nil {
		mqtt.WARN.Println(token.Error())
		broker.Client = nil
	}

	return token.Error()
}

func (broker *Broker) IsConnected() bool {
	return broker.Client != nil
}

func (broker *Broker) Disconnect() {
	if broker.Client.IsConnected() {
		broker.Client.Disconnect(10)
		broker.Client = nil
	}
}

func (broker *Broker) SubscribeMultiple(filter map[string]byte) mqtt.Token {
	for key, value := range filter {
		mqtt.DEBUG.Printf("Subscribe Multiple(topic:%s, qos:%d)", key, value)
	}

	return broker.Client.SubscribeMultiple(filter, broker.multipleMessageHandler())
}

func (broker *Broker) multipleMessageHandler() mqtt.MessageHandler {
	return func(c mqtt.Client, msg mqtt.Message) {
		if !c.IsConnected() {
			panic("mqtt not connected")
		}

		for _, topic := range broker.topics {
			topicName := topic.GetTopic()
			msgTopicName := msg.Topic()
			if isWildCard(topicName, msgTopicName) || topicName == msgTopicName {
				topic.MessagePrintHandler()
			}
		}
	}
}

func isWildCard(reg, taget string) bool {
	if strings.Contains(reg, "+") {
		reg = strings.ReplaceAll(reg, "+", ".*")
	} else if strings.Contains(reg, "+") {
		reg = strings.ReplaceAll(reg, "+", ".*")
	} else {
		return false
	}

	if r, err := regexp.Compile(reg); err == nil {
		if r.MatchString(taget) {
			return true
		}
	}

	return false
}
