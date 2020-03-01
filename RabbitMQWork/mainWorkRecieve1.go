package main

import "rabbitmq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
		"imoocWork")
	defer rabbitmq.Destroy()
	rabbitmq.ConsumeSimple()
	
}