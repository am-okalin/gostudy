package send

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"testing"
)

func TestDeclare1(t *testing.T) {
	// 初始化连接
	ch, err := conn.Channel()
	if err != nil {
		t.Error(err)
		return
	}
	defer ch.Close()

	// 定义队列
	queue, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	t.Log(queue, err)
}

func TestSend(t *testing.T) {
	// 初始化连接
	ch, err := conn.Channel()
	if err != nil {
		t.Error(err)
		return
	}

	// 推送数据至已存在的队列
	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",      // exchange
		"hello", // routing key (可指定为队列名)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		t.Error(err)
	}

	log.Printf(" [x] Sent %s\n", body)
}
