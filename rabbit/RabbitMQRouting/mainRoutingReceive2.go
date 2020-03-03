package main

import "rabbitmq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_two")
	defer rabbitmq.Destroy()
	rabbitmq.ConsumeRouting()
}
