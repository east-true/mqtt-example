package mqttgo

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	addr := "127.0.0.1:1883"
	id := "mqtt-client"
	client := mqtt.NewClient(mqtt.NewClientOptions().SetClientID(id).AddBroker(addr))

	token := client.Connect()
	_ = token.Wait()
	if token.Error() != nil {
		fmt.Println(token.Error())
		return
	} else {
		defer client.Disconnect(10)
	}

	topic := "/test/a"
	client.Subscribe(topic, byte(0), func(c mqtt.Client, m mqtt.Message) {
		name := m.Topic()
		payload := m.Payload()

		fmt.Printf("%v %v\n", name, string(payload))
	})

	token = client.Unsubscribe(topic)
	_ = token.Wait()
	if token.Error() != nil {
		fmt.Println(token.Error())
		return
	}
}
