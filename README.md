# mqtt-example
https://github.com/eclipse/paho.mqtt.golang

### Subscribe
```go
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

	client.Publish(topic, qos, true, []byte("Hi"))
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
```
### SubscribeMultiple
```go
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
		client.Publish(topic, qos, true, []byte("Hi"))
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
```
### Publish
```go
func TestPublish(t *testing.T) {
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
	for i := 0; i < 10; i++ {
		client.Publish(topic, qos, true, []byte{byte(i)})
		go func(topic string) {
			_ = token.Wait()
			if token.Error() != nil {
				t.Errorf("%v %v", topic, token.Error())
			}
		}(topic)
	}
}
```
