package main

import "rabbitmq/rabbit/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
		"imoocWork")
	defer rabbitmq.Destroy()
	rabbitmq.ConsumeSimple()
	
}