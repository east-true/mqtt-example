package mqtt

func main() {
	mq := MQTT{}
	if err := mq.Conn(); err != nil {
		// TODO : ERR LOG
	}

	defer mq.Close()
	token := mq.MultipleSubscribe()
	if !token.Wait() {
		// TODO : ERR LOG : token.Error()
	}
}
