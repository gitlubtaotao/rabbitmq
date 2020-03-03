package main

import "rabbitmq/rabbit/RabbitMQ"

func main() {
	rabbit := RabbitMQ.NewRabbitMQSimple(
		"imoocSimple")
	rabbit.ConsumeSimple()
	defer rabbit.Destroy()
}
