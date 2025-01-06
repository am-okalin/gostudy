package main

import (
	"log"
	"pkg_third/mq_rabbit/rbmq"
)

var conn = rbmq.NewConn("amqp://root:123456@localhost:5672/")

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 声明队列, 独占模式(该临时队列只允许本次连接使用)
	q, err := ch.QueueDeclare(
		"",    // 生成随机队列名
		false, // durable 非持久化队列
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 将独占队列bind到交换机
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to bind a queue")

	// 消费消息
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack 关闭自动确认
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf(" [x] %s", d.Body)
	}
}
