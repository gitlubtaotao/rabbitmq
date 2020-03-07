package main

import (
	
	"rabbitmq/rabbitmq/RabbitMQ"
)

func main()  {
	imoocOne:=RabbitMQ.NewRabbitMQRouting("exImooc","imooc_one")
	imoocOne.RecieveRouting()
}
