package mqttgo

import (
	"errors"
	"log"
	"mqttgo/broker"
	"mqttgo/broker/opt"
	"mqttgo/topic"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	addr   string        // [tcp | ssl | ws]://[IP]:[PORT]
	topics []topic.Topic // topic.Topic | topic.ValueTopic | topic.JsonTopic
	broker *broker.Broker
}

// TODO : param : mqtt.config
func New(address string) *MQTT {
	// TODO : mqtt.Config(topics) - > make topics(topic, value, json) -> broker.New(address, topics...)
	return &MQTT{
		addr:   address,
		broker: broker.New(address, nil),
	}
}

func (mq *MQTT) SetLogger() {
	mqtt.ERROR = log.New(os.Stdout, "ERR   ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "DEBUG ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "CRITIC", 0)
	mqtt.WARN = log.New(os.Stdout, "WARN  ", 0)
}

func (mq *MQTT) Connect() error {
	opt := new(opt.Option)
	return mq.broker.Connect(opt.Get(mq.addr))
}

func (mq *MQTT) IsConnected() bool {
	return mq.broker.IsConnected()
}

func (mq *MQTT) Disconnect() {
	if mq.broker.IsConnected() {
		if err := mq.UnsubscribeAllBlock(); err != nil {
			mqtt.ERROR.Println(err)
		}
		mq.broker.Disconnect()
	}
}

func (mq *MQTT) Subscribe(out chan<- *topic.TopicValue) {
	if len(mq.topics) > 1 {
		mqtt.ERROR.Println("It has several topics, you can use func SubscribeAll()")
		return
	}

	token := mq.topics[0].Subscribe(mq.broker.Client, mq.topics[0].MessageDeliveryHandler(out))
	go func() {
		_ = token.Wait()
		if token.Error() != nil {
			mqtt.ERROR.Println(token.Error())
		}
	}()
}

func (mq *MQTT) SubscribeBlock(out chan<- *topic.TopicValue) error {
	if len(mq.topics) > 1 {
		return errors.New("It has several topics, you can use func SubscribeAll()")
	}

	token := mq.topics[0].Subscribe(mq.broker.Client, mq.topics[0].MessageDeliveryHandler(out))
	_ = token.Wait()
	return token.Error()
}

func (mq *MQTT) SubscribeAll(out chan<- *topic.TopicValue) {
	token := mq.broker.SubscribeMultipleDelivery(mq.getFilter(), out)
	go func() {
		_ = token.Wait()
		if token.Error() != nil {
			mqtt.ERROR.Println(token.Error())
		}
	}()
}

func (mq *MQTT) SubscribeAllBlock(out chan<- *topic.TopicValue) error {
	token := mq.broker.SubscribeMultipleDelivery(mq.getFilter(), out)
	_ = token.Wait()
	return token.Error()
}

func (mq *MQTT) getFilter() map[string]byte {
	filter := make(map[string]byte)
	for _, topic := range mq.topics {
		filter[topic.Topic] = byte(topic.Qos)
	}

	return filter
}

func (mq *MQTT) Unsubscribe(topic string) {
	for i := range mq.topics {
		if mq.topics[i].Topic == topic {
			token := mq.topics[i].Unsubscribe(mq.broker.Client)
			go func() {
				_ = token.Wait()
				if token.Error() != nil {
					mqtt.ERROR.Println(token.Error())
				}
			}()
		}
	}
}

func (mq *MQTT) UnsubscribeAll() {
	for i := range mq.topics {
		token := mq.topics[i].Unsubscribe(mq.broker.Client)
		go func() {
			_ = token.Wait()
			if token.Error() != nil {
				mqtt.ERROR.Println(token.Error())
			}
		}()
	}
}

func (mq *MQTT) UnsubscribeBlock(topic string) error {
	for i := range mq.topics {
		if mq.topics[i].Topic == topic {
			token := mq.topics[i].Unsubscribe(mq.broker.Client)
			_ = token.Wait()
			if token.Error() != nil {
				return token.Error()
			}
		}
	}

	return nil
}

func (mq *MQTT) UnsubscribeAllBlock() error {
	for i := range mq.topics {
		token := mq.topics[i].Unsubscribe(mq.broker.Client)
		_ = token.Wait()
		if token.Error() != nil {
			return token.Error()
		}
	}

	return nil
}

func (mq *MQTT) Publish(topic string, payload []byte) error {
	for i := range mq.topics {
		if mq.topics[i].Topic == topic {
			token := mq.topics[i].Publish(mq.broker.Client, payload)
			_ = token.Wait()
			if token.Error() != nil {
				return token.Error()
			}
		}
	}

	return nil
}
