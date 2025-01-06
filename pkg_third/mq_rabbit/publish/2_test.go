package send

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"testing"
)

func TestDeclare2(t *testing.T) {
	// 初始化连接
	ch, err := conn.Channel()
	if err != nil {
		t.Error(err)
		return
	}
	defer ch.Close()

	// 定义队列
	queue, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	t.Log(queue, err)
}

func TestTask(t *testing.T) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(ctx,
		"",           // exchange
		"task_queue", // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	t.Logf(" [x] Sent %s", body)
}
