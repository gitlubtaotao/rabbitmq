package main

import "rabbitmq/rabbit/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_one")
	defer rabbitmq.Destroy()
	rabbitmq.ConsumeRouting()
}
