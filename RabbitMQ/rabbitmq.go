package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const MQURL = "amqp://guest:guest@127.0.0.1:5672/imooc"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl string
}

//创建结构体实例
func NewRabbitMQ(QueueName, Exchange, Key string) *RabbitMQ {
	return &RabbitMQ{QueueName: QueueName, Exchange: Exchange, Key: Key, Mqurl: MQURL}
}

//断开channel 和 connection
func (r *RabbitMQ) Destroy() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

//自定义错误，统一处理
func failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	failOnErr(err, "failed to connect rabb"+
		"iamb!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	failOnErr(err, "failed to open a channel")
	return rabbitmq
}
func (r *RabbitMQ) PublishSimple(message string) {
	//创建队列
	_, err := r.channel.QueueDeclare(r.QueueName,
		false,  //是否持久化
		false, //是否自动删除
		false, //是否具有排他性
		false, //是否阻塞处理
		nil) //额外的属性
	if err != nil {
		log.Fatalf("%s", err)
	}
	//发布消息
	err = r.channel.Publish(r.Exchange,
		r.QueueName,
		false, //如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		
		false, //如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	failOnErr(err, "publish is error")
}

func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		log.Fatalf("%s", err)
	}
	//接收消息
	megs, err := r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		true, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("%s", err)
	}
	forever := make(chan bool)
	go func() {
		for d := range megs {
			//消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
