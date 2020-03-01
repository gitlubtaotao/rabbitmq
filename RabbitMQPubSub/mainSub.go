package main

import "rabbitmq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" +
		"newProduct")
	defer rabbitmq.Destroy()
	rabbitmq.ConsumePub()
}
