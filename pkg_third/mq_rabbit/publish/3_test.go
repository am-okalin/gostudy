package send

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"testing"
)

func TestDeclare3(t *testing.T) {
	// 初始化连接
	ch, err := conn.Channel()
	if err != nil {
		t.Error(err)
		return
	}
	defer ch.Close()

	// 定义交换机
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

}

func TestLog(t *testing.T) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// exchange 不存在时会报错
	// 没有 queue 绑定至 exchange 时, 丢弃消息
	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(ctx,
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	t.Logf(" [x] Sent %s", body)
}
