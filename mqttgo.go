package mqttgo

import (
	"log"
	"mqttgo/broker"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	addr   string // format : [tcp | ssl | ws]://[IP]:[PORT]
	topics []string
	broker *broker.Broker
}

func (_ MQTT) SetLogger() {
	mqtt.ERROR = log.New(os.Stdout, "ERR   ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "DEBUG ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "CRITIC", 0)
	mqtt.WARN = log.New(os.Stdout, "WARN  ", 0)
}

// func (mq *MQTT) getFilter() (filter map[string]byte) {
// 	filter = make(map[string]byte)
// 	for _, topic := range mq.topics {
// 		if topic.Topic == "" {
// 			// log
// 			continue
// 		}

// 		filter[topic.Topic] = byte(topic.Qos)
// 	}

// 	return
// }
