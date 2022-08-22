package observer

import (
	"context"
	"fmt"
)

//Subscriber 订阅者接口
type Subscriber interface {
	update(ctx context.Context) //更新订阅者内容
	getID() string              //获取订阅者ID
}

type customer struct {
	id string
}

func (c *customer) update(ctx context.Context) {
	fmt.Printf("Sending email to customer %s\n", c.id)
}

func (c *customer) getID() string {
	return c.id
}
