package main

import (
	"fmt"
	"rabbitmq/RabbitMQ"
)

func main() {
	rabbit := RabbitMQ.NewRabbitMQSimple("imoocSimple")
	for i := 0; i < 16; i++ {
		rabbit.PublishSimple("Hello imooc!")
	}
	defer rabbit.Destroy()
	fmt.Println("发送成功！")
}
