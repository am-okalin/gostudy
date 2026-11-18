package main

import (
	"context"
	"net"
	"sync"
	"time"
)

func now() string { return time.Now().Format("15:04:05") }

// Server 维护所有客户端连接，负责广播。
type Server struct {
	ln    net.Listener
	mu    sync.RWMutex
	conns map[*Conn]struct{}
}

func NewServer(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{ln: ln, conns: make(map[*Conn]struct{})}, nil
}

func (s *Server) Addr() string { return s.ln.Addr().String() }

func (s *Server) Serve(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.ln.Close() // 让下面的 Accept 立刻返回
	}()

	for {
		nc, err := s.ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil // 正常退出
			}
			return err
		}
		// 每来一个连接起一个 goroutine —— 这就是"Go 写长连接服务"的核心模型。
		go s.handle(ctx, nc)
	}
}

// Conn 是一条客户端连接。
type Conn struct {
	nc   net.Conn
	srv  *Server
	name string
	send chan Message // 单写者：所有出站消息都走这里
}

func (s *Server) add(c *Conn) {
	s.mu.Lock()
	s.conns[c] = struct{}{}
	s.mu.Unlock()
}

func (s *Server) remove(c *Conn) {
	s.mu.Lock()
	delete(s.conns, c)
	s.mu.Unlock()
}

// broadcast 扇出；某个客户端读得慢（send 满了）就丢它这一条，绝不阻塞别人。
func (s *Server) broadcast(m Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for c := range s.conns {
		select {
		case c.send <- m:
		default:
		}
	}
}

func (s *Server) handle(ctx context.Context, nc net.Conn) {
	defer nc.Close()

	// 第一条必须是 hello，用来拿名字，同时也算一次最简握手。
	nc.SetReadDeadline(time.Now().Add(10 * time.Second))
	hello, err := readFrame(nc)
	if err != nil || hello.Type != "hello" {
		return
	}

	name := hello.Name
	if name == "" {
		name = "匿名"
	}
	c := &Conn{nc: nc, srv: s, name: name, send: make(chan Message, 64)}

	s.add(c)
	defer func() {
		s.remove(c)
		s.broadcast(Message{Type: "system", Text: name + " 离开了", Time: now()})
	}()

	go c.writePump(ctx)

	s.broadcast(Message{Type: "system", Text: name + " 加入了", Time: now()})

	for {
		// 每条消息都重置读超时：60 秒内一条数据（含 pong）都没有，就认为连接已死。
		// 这是应用层心跳；TCP keepalive 默认要 7200 秒，太慢了。
		nc.SetReadDeadline(time.Now().Add(60 * time.Second))

		m, err := readFrame(nc)
		if err != nil {
			return // 断开 / 超时 / 协议错
		}

		switch m.Type {
		case "chat":
			s.broadcast(Message{Type: "chat", Name: name, Text: m.Text, Time: now()})
		case "pong":
			// 心跳回应，读超时上面已经重置，这里不用做别的
		}
	}
}

// writePump 是这条连接唯一的写者：业务消息和心跳 ping 都从这里顺序写出。
func (c *Conn) writePump(ctx context.Context) {
	ping := time.NewTicker(25 * time.Second)
	defer ping.Stop()

	write := func(m Message) bool {
		c.nc.SetWriteDeadline(time.Now().Add(5 * time.Second))
		return writeFrame(c.nc, m) == nil
	}

	for {
		select {
		case m := <-c.send:
			if !write(m) {
				c.nc.Close()
				return
			}
		case <-ping.C:
			if !write(Message{Type: "ping"}) {
				c.nc.Close() // 写不出去说明对端已经没了
				return
			}
		case <-ctx.Done():
			c.nc.Close()
			return
		}
	}
}
