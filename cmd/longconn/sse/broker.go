package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Event 是一条 SSE 消息。
// Type 对应 SSE 的 event: 字段，留空就是默认的 "message"。
type Event struct {
	ID   int64
	Type string
	Data string
}

// Broker 是极简的发布订阅：每个订阅者一个带缓冲的 channel。
type Broker struct {
	mu      sync.RWMutex
	clients map[chan Event]struct{}
	seq     atomic.Int64 // 全局事件序号，用作 SSE 的 id:
}

func NewBroker() *Broker {
	return &Broker{clients: make(map[chan Event]struct{})}
}

// Subscribe 返回消息 channel 和取消订阅函数。
// 缓冲 16：客户端稍微慢一点没关系，慢太多就在 Publish 里被丢弃。
func (b *Broker) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 16)

	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()

	var once sync.Once
	unsub := func() {
		once.Do(func() {
			b.mu.Lock()
			delete(b.clients, ch)
			close(ch) // 持写锁关闭，保证 Publish 的读锁不会撞上"发往已关闭 channel"
			b.mu.Unlock()
		})
	}
	return ch, unsub
}

// Publish 把事件扇出给所有订阅者。
// 某个客户端读得慢（channel 满了）就丢掉它这一条，绝不阻塞发布方。
func (b *Broker) Publish(ev Event) {
	if ev.ID == 0 {
		ev.ID = b.seq.Add(1)
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.clients {
		select {
		case ch <- ev:
		default: // 缓冲满：慢客户端，丢
		}
	}
}

// tick 周期性推一条消息，模拟服务端主动推送。
func (b *Broker) tick(d time.Duration) {
	t := time.NewTicker(d)
	defer t.Stop()
	for now := range t.C {
		b.Publish(Event{Type: "message", Data: "服务端定时推送 " + now.Format("15:04:05")})
	}
}

// PublishHTTP 处理 POST /publish：把请求体当成一条消息广播出去。
func (b *Broker) PublishHTTP(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Text string `json:"text"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Text == "" {
		http.Error(w, `需要 JSON: {"text":"..."}`, http.StatusBadRequest)
		return
	}
	b.Publish(Event{Data: in.Text})
	w.WriteHeader(http.StatusNoContent)
}
