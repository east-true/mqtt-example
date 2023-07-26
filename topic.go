package mqtt

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Topic struct {
	Topic string
	Qos   int    // [1 || 2 || 3]
	Name  string // ["" || tag name || json key]
	Value string // ["" || tag name || json key]
}

func (topic *Topic) Subscribe(client mqtt.Client) {
	token := client.Subscribe(topic.Topic, byte(topic.Qos), topic.messageHandler(input))
	token.Wait()
}

func (topic *Topic) messageHandler() mqtt.MessageHandler {
	return func(aClient mqtt.Client, aMessage mqtt.Message) {
		var (
			msgTopic string = topic.WildCard(aMessage.Topic())
			payload  []byte = aMessage.Payload()
			// time     time.Time = time.Now()
		)

		if topic.Name == "" {
			topic.Name = strings.ReplaceAll(msgTopic[1:], "/", "_")
		}

		// payload is numeric
		payloadstr := string(payload)
		if topic.isNumeric(payloadstr) {
			// Name:     topic.Name,
			// Time:     time,
			// Value:    []interface{}{payloadstr},
			return
		}

		// json parsing
		var msgMap map[string]interface{}
		err := json.Unmarshal(payload, &msgMap)
		if err != nil {
			// TODO : ERR LOG
			//("[%s] %v : %v", topic.Topic, topic.Name, err)
			return
		}

		// get value
		name, ok := msgMap[topic.Name]
		if !ok {
			// TODO : ERR LOG
			//(`[%s] Not Found Key "%s" : %s`, topic.Topic, topic.Name, string(payload))
			return
		}

		value, ok := msgMap[topic.Value]
		if !ok {
			// TODO : ERR LOG
			//(`[%s] Not Found Key "%s" : %s`, topic.Topic, topic.Name, string(payload))
			return
		}

		// if array
		// switch t := value.(type) {
		// case []interface{}:
		// for _, idata := range t {
		// Name:     name.(string),
		// Time:     time,
		// Value:    []interface{}{idata},
		// }
		// case interface{}:
		// Name:     name.(string),
		// Time:     time,
		// Value:    []interface{}{t},
		// default:
		// 	return
		// }
	}
}

func (t *Topic) isNumeric(aInterface interface{}) bool {
	_, sErr := strconv.ParseFloat(fmt.Sprint(aInterface), 64)
	return sErr == nil
}

func (t *Topic) WildCard(topic string) string {
	var reg string
	if strings.Contains(topic, "+") {
		reg = strings.ReplaceAll(topic, "+", "*")
	} else if strings.Contains(topic, "#") {
		reg = strings.ReplaceAll(topic, "#", "*+")
	}

	r, err := regexp.Compile(reg)
	if err != nil {
		// TODO : ERR LOG
	}

	if r.MatchString(topic) {
		return topic
	}

	return ""
}
