package main_test

import (
	"fmt"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestSubscribe(t *testing.T) {
	addr := "127.0.0.1:1883"
	id := "mqtt-client"
	client := mqtt.NewClient(mqtt.NewClientOptions().SetClientID(id).AddBroker(addr))

	token := client.Connect()
	_ = token.Wait()
	if token.Error() != nil {
		t.Error(token.Error())
	} else {
		defer client.Disconnect(10)
	}

	topic := "/test/a"
	qos := byte(0)
	client.Subscribe(topic, qos, func(c mqtt.Client, m mqtt.Message) {
		name := m.Topic()
		payload := m.Payload()

		fmt.Printf("%v %v\n", name, string(payload))
	})

	client.Publish(topic, qos, true, "Hi")
	_ = token.Wait()
	if token.Error() != nil {
		t.Error(token.Error())
	}

	token = client.Unsubscribe(topic)
	_ = token.Wait()
	if token.Error() != nil {
		t.Error(token.Error())
	}
}

func TestSubscribeMultiple(t *testing.T) {
	addr := "127.0.0.1:1883"
	id := "mqtt-client"
	client := mqtt.NewClient(mqtt.NewClientOptions().SetClientID(id).AddBroker(addr))

	token := client.Connect()
	_ = token.Wait()
	if token.Error() != nil {
		t.Error(token.Error())
	} else {
		defer client.Disconnect(10)
	}

	filter := map[string]byte{
		"/test/a": 0,
		"/test/b": 0,
		"/test/c": 0,
	}
	client.SubscribeMultiple(filter, func(c mqtt.Client, m mqtt.Message) {
		name := m.Topic()
		payload := m.Payload()

		fmt.Printf("%v %v\n", name, string(payload))
	})

	for topic, qos := range filter {
		client.Publish(topic, qos, true, "Hi")
		_ = token.Wait()
		if token.Error() != nil {
			t.Error(token.Error())
		}
	}

	for topic := range filter {
		token = client.Unsubscribe(topic)
		_ = token.Wait()
		if token.Error() != nil {
			t.Error(token.Error())
		}
	}
}
