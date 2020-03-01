package main

import "rabbitmq/RabbitMQ"

func main() {
	rabbit := RabbitMQ.NewRabbitMQSimple(
		"imoocSimple")
	rabbit.ConsumeSimple()
}
