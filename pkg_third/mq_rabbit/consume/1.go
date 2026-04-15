package main

import (
	"log"
)

//var conn = rbmq.NewConn("amqp://root:123456@localhost:5672/")
//
//func failOnError(err error, msg string) {
//	if err != nil {
//		log.Panicf("%s: %s", msg, err)
//	}
//}

func main() {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		"hello", // queue
		"",      // consumer
		true,    // auto-ack 开启自动确认
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
	}
}
