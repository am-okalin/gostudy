package rbmq

import (
	//"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

const (
	JSON = "application/json"
	UTF8 = "utf-8"
)

type Conn struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	locker  sync.Locker
	url     string
}

func NewConn(url string) *Conn {
	conn, err := amqp.Dial(url)
	if err != nil {
		println(err.Error())
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		println(err.Error())
		panic(err)
	}

	return &Conn{
		conn:    conn,
		channel: channel,
		url:     url,
		locker:  &sync.Mutex{},
	}
}

// Channel 尽可能的复用channel
func (c *Conn) Channel() (*amqp.Channel, error) {
	// 直接返回
	if !c.channel.IsClosed() {
		return c.channel, nil
	}

	// 懒加载channel
	var err error
	c.locker.Lock()
	defer c.locker.Unlock()

	if c.conn.IsClosed() {
		if c.conn, err = amqp.Dial(c.url); err != nil {
			return nil, err
		}
	}

	if c.channel.IsClosed() {
		if c.channel, err = c.conn.Channel(); err != nil {
			return nil, err
		}
	}

	return c.channel, nil
}

//func (c *Conn) Publish(ctx context.Context, routekey string, body []byte, opts ...ProducerOption) error {
//
//	channel, err := c.Channel()
//	if err != nil {
//		return err
//	}
//
//	publishing := amqp.Publishing{
//		Headers:         amqp.Table{},
//		ContentType:     JSON,
//		ContentEncoding: UTF8,
//		Body:            body,
//		DeliveryMode:    amqp.Persistent, // 消息持久化模式
//		Priority:        0,
//	}
//
//	for _, o := range opts {
//		o(&publishing)
//	}
//
//	return channel.PublishWithContext(ctx, p.exchange, routekey, true, false, publishing)
//}
