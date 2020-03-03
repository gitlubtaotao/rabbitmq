package main

import "rabbitmq/rabbit/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" +
		"newProduct")
	defer rabbitmq.Destroy()
	rabbitmq.ConsumePub()
}
