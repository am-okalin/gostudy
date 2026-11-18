package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// inbound 是客户端发上来的消息。
type inbound struct {
	Text string `json:"text"`
}

// outbound 是广播给所有客户端的消息。
type outbound struct {
	Time   string `json:"time"`
	Name   string `json:"name"`
	Text   string `json:"text"`
	System bool   `json:"system,omitempty"`
}

// Hub 维护所有在线连接，负责广播。
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]struct{}
}

func NewHub() *Hub { return &Hub{clients: make(map[*Client]struct{})} }

func (h *Hub) add(c *Client) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) remove(c *Client) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
}

// broadcast 扇出给所有连接。
// 慢客户端（send 缓冲满）直接丢这一条，绝不让一个卡住的浏览器拖住整个房间。
func (h *Hub) broadcast(msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		select {
		case c.send <- msg:
		default:
		}
	}
}

func (h *Hub) broadcastSystem(text string) {
	msg, _ := json.Marshal(outbound{Time: time.Now().Format("15:04:05"), Text: text, System: true})
	h.broadcast(msg)
}

// Client 是一个 WebSocket 连接。
// 所有出站消息都走 send channel，由 writePump 单 goroutine 写出。
type Client struct {
	conn *websocket.Conn
	hub  *Hub
	name string
	send chan []byte
}

// ServeWS 处理 GET /ws。
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	// InsecureSkipVerify 只用于本地演示：不校验 Origin。
	// 线上必须配 OriginPatterns，否则任意站点都能连上来（相当于 WS 版 CSRF）。
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		logger.Warn("accept 失败", "err", err)
		return
	}
	defer conn.CloseNow()

	conn.SetReadLimit(1 << 20) // 单条消息上限 1MB

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "匿名"
	}

	c := &Client{conn: conn, hub: h, name: name, send: make(chan []byte, 32)}
	h.add(c)
	defer func() {
		h.remove(c)
		h.broadcastSystem(name + " 离开了")
	}()

	ctx := r.Context()

	go c.writePump(ctx)
	h.broadcastSystem(name + " 加入了")

	for {
		typ, data, err := conn.Read(ctx)
		if err != nil {
			return // 客户端断开 / 超时 / 被踢
		}
		if typ != websocket.MessageText {
			continue // 本 demo 只处理文本
		}

		var in inbound
		if err := json.Unmarshal(data, &in); err != nil || in.Text == "" {
			continue
		}
		out, _ := json.Marshal(outbound{
			Time: time.Now().Format("15:04:05"),
			Name: name,
			Text: in.Text,
		})
		h.broadcast(out)
	}
}

// writePump 是这条连接唯一的写者。
// WebSocket 不允许并发写，所以所有发送都收口到 c.send，
// 再从这里顺序写出——这是和普通 HTTP handler 最大的实现差异之一。
func (c *Client) writePump(ctx context.Context) {
	ping := time.NewTicker(30 * time.Second)
	defer ping.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.send:
			wctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := c.conn.Write(wctx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				return
			}
		case <-ping.C:
			// 心跳：浏览器会自动回 pong，前端不需要写任何代码
			pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := c.conn.Ping(pctx)
			cancel()
			if err != nil {
				return
			}
		}
	}
}
