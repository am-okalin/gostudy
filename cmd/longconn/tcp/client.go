package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// runClient 是一个极简命令行客户端：stdin 每行发一条 chat，收到的消息打到 stdout。
func runClient(addr, name string) error {
	nc, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return err
	}
	defer nc.Close()

	// 写收口：读 goroutine 要回 pong，主 goroutine 要发 chat，
	// 两处同时写会交错成乱帧，所以必须加锁。
	// 服务端那边则是用"单写者 goroutine + channel"，两种做法都对。
	var writeMu sync.Mutex
	send := func(m Message) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		nc.SetWriteDeadline(time.Now().Add(5 * time.Second))
		return writeFrame(nc, m)
	}

	if err := send(Message{Type: "hello", Name: name}); err != nil {
		return err
	}

	go func() {
		for {
			m, err := readFrame(nc)
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return // 本端主动关闭，安静退出
				}
				fmt.Println("\n连接断开:", err)
				os.Exit(0)
			}
			switch m.Type {
			case "chat":
				fmt.Printf("[%s] %s: %s\n", m.Time, m.Name, m.Text)
			case "system":
				fmt.Printf("[%s] * %s\n", m.Time, m.Text)
			case "ping":
				send(Message{Type: "pong"})
			}
		}
	}()

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		text := strings.TrimSpace(sc.Text())
		if text == "" {
			continue
		}
		if err := send(Message{Type: "chat", Text: text}); err != nil {
			return err
		}
	}
	// stdin 结束：给对端消息一点打印时间，再主动断开
	time.Sleep(200 * time.Millisecond)
	nc.Close()
	return nil
}
