package main

import (
	"fmt"
	"rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	imoocOne := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_one")
	imoocTwo := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_two")
	for i := 0; i <= 20; i++ {
		go func(i int) {
			str := "Hello imooc one!" + strconv.Itoa(i)
			imoocOne.PublishRouting(str)
			//runtime.Gosched()
			fmt.Println(str)
		}(i)
	}
	for i := 0; i <= 10; i++ {
		go func(i int) {
			for {
				str := "Hello imooc Two!" + strconv.Itoa(i)
				imoocTwo.PublishRouting(str)
				//runtime.Gosched()
				fmt.Println(str)
			}
			
		}(i)
	}
	time.Sleep(time.Minute * 1)
}
