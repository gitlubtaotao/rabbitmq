package main

import (
	"fmt"
	"rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
		"imoocWork")
	defer rabbitmq.Destroy()
	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("Hello imooc!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
	fmt.Println("发送成功！")
	
}
