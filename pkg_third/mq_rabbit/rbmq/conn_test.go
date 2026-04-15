package rbmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

var conn = NewConn("amqp://root:123456@localhost:5672/")

func TestDeclare(t *testing.T) {
	// 初始化连接
	ch, err := conn.Channel()
	if err != nil {
		t.Error(err)
		return
	}

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

func Test1(t *testing.T) {
	// 初始化连接
	url := "amqp://root:123456@localhost:5672/"
	conn := NewConn(url)
	ch, err := conn.Channel()
	if err != nil {
		t.Error(err)
		return
	}

	// 生产消息
	err = ch.Publish("", "hello", false, false, amqp.Publishing{})
	if err != nil {
		t.Error(err)
		return
	}

	// 消费消息
	t.Log("done")
}
