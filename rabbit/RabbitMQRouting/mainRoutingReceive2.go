package main

import "rabbitmq/rabbit/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_two")
	defer rabbitmq.Destroy()
	rabbitmq.ConsumeRouting()
}
