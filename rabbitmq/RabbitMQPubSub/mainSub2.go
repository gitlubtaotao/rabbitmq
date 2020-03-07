package main

import (
	
	"rabbitmq/rabbitmq/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" +
		"newProduct")
	rabbitmq.RecieveSub()
}
