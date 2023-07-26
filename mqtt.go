package mqtt

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	Addr   string // format : [tcp | ssl | ws]://[IP]:[PORT]
	Option Option
	Topics []Topic
	client mqtt.Client
}

func (mq *MQTT) SetServerLogger() {
	mqtt.ERROR = log.New(os.Stdout, "ERR   ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "DEBUG ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "CRITIC", 0)
	mqtt.WARN = log.New(os.Stdout, "WARN  ", 0)
}

func (mq *MQTT) Conn() error {
	opt := mq.Option.get(mq.Addr)
	client := mqtt.NewClient(opt)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	mq.client = client
	return nil
}

func (mq *MQTT) Close() {
	if mq.client.IsConnected() {
		for _, topic := range mq.Topics {
			token := mq.client.Unsubscribe(topic.Name)
			if !token.Wait() {
				// TODO : ERR LOG
			}
		}
		mq.client.Disconnect(10)
	}
}

func (mq *MQTT) MultipleSubscribe() mqtt.Token {
	// for key, value := range mq.GetFilter() {
	// 	// TODO :  Multiple Subscribe LOG
	// }

	return mq.client.SubscribeMultiple(mq.GetFilter(), mq.multipleMessageHandler())
}

func (mq *MQTT) GetFilter() (filter map[string]byte) {
	filter = make(map[string]byte)

	for _, topic := range mq.Topics {
		if topic.Topic == "" {
			// log
			continue
		}

		filter[topic.Topic] = byte(topic.Qos)
	}

	return
}

func (mq *MQTT) multipleMessageHandler() mqtt.MessageHandler {
	return func(aClient mqtt.Client, aMessage mqtt.Message) {
		var (
			msgTopic string = aMessage.Topic()
		)

		if !aClient.IsConnected() {
			panic("mqtt not connected")
		}

		topic := mq.GetTopic(msgTopic)
		if topic == nil {
			return
		}

		handler := topic.messageHandler()
		handler(aClient, aMessage)
	}
}

func (mq *MQTT) GetTopic(msgTopic string) *Topic {
	for _, topic := range mq.Topics {
		if msgTopic == topic.Topic || topic.WildCard(msgTopic) != "" {
			return &topic
		}
	}

	return nil
}
