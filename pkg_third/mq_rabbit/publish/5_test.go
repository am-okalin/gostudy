package send

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"testing"
)

func TestDeclare5(t *testing.T) {
	// 初始化连接
	ch, err := conn.Channel()
	if err != nil {
		t.Error(err)
		return
	}
	defer ch.Close()

	// 定义交换机
	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

}

func TestLogTopic(t *testing.T) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// exchange 不存在时会报错
	// 没有 queue 绑定至 exchange 时, 丢弃消息
	os.Args = []string{"", "cron.critical", "A cron error"}
	//os.Args = []string{"", "kern.critical", "A critical kernel error"}
	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(ctx,
		"logs_topic",           // exchange
		severityFrom5(os.Args), // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	t.Logf(" [x] Sent %s", body)
}

func severityFrom5(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}
