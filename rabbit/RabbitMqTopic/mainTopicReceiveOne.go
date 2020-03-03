package main

import "rabbitmq/rabbit/RabbitMQ"

func main()  {
	imoocOne:= RabbitMQ.NewRabbitMQTopic("exImoocTopic","imooc.*.two")
	imoocOne.RecieveTopic()
}

