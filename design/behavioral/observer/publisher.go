package observer

import "context"

//PublisherInterface 发布者接口
type PublisherInterface interface {
	Subscribe(s Subscriber)
	Unsubscribe(s Subscriber)
	Notify()
}

type Publisher struct {
	sl []Subscriber
}

//Subscribe 添加订阅者至 p.sl
func (p *Publisher) Subscribe(s Subscriber) {
	p.sl = append(p.sl, s)
}

//Unsubscribe 从 p.sl 删除订阅者
func (p *Publisher) Unsubscribe(s Subscriber) {
	sll := len(p.sl)
	for i, v := range p.sl {
		if s.getID() == v.getID() {
			p.sl[sll-1], p.sl[i] = p.sl[i], p.sl[sll-1]
			p.sl = p.sl[:sll-1]
			return
		}
	}
}

//Notify 通知所有订阅者
func (p *Publisher) Notify() {
	for _, s := range p.sl {
		s.update(context.Background())
	}
}
