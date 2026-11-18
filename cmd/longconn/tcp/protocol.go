package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// TCP 是字节流，没有"消息边界"：你连着发两条，对方可能一次读到两条粘在一起，
// 也可能一条只读到一半。所以应用层必须自己划边界，常见三种做法：
//
//	1. 长度前缀（本 demo）：4 字节大端长度 + 负载，支持任意二进制；
//	2. 分隔符：\n 或 \0 结尾，简单但不能传裸二进制，需要转义；
//	3. 固定长度：简单但浪费，不适合变长消息。
//
// 选长度前缀，是因为它是 HTTP/2、WebSocket、gRPC 底下共同的做法。

const (
	maxFrame  = 1 << 20 // 单帧上限 1MB，防止对方发个 4GB 长度把内存打爆
	headerLen = 4
)

// Message 是应用层协议。Type 区分消息种类。
type Message struct {
	Type string `json:"type"`           // hello / chat / ping / pong / system
	Name string `json:"name,omitempty"` // 发送者
	Text string `json:"text,omitempty"`
	Time string `json:"time,omitempty"`
}

// writeFrame 一次 Write 把「头 + 体」发出去。
// 分两次写会产生两个 TCP 段，在 Nagle 算法下可能被延迟合并，没必要。
func writeFrame(w io.Writer, m Message) error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if len(body) > maxFrame {
		return fmt.Errorf("帧太大: %d 字节", len(body))
	}

	buf := make([]byte, headerLen+len(body))
	binary.BigEndian.PutUint32(buf[:headerLen], uint32(len(body)))
	copy(buf[headerLen:], body)

	_, err = w.Write(buf)
	return err
}

// readFrame 先读 4 字节长度，再按长度读满负载。
// io.ReadFull 会一直阻塞到读满为止，正好解决"半包"。
func readFrame(r io.Reader) (Message, error) {
	var head [headerLen]byte
	if _, err := io.ReadFull(r, head[:]); err != nil {
		return Message{}, err
	}

	n := binary.BigEndian.Uint32(head[:])
	if n > maxFrame {
		return Message{}, errors.New("帧长度超过上限，拒绝分配内存")
	}

	body := make([]byte, n)
	if _, err := io.ReadFull(r, body); err != nil {
		return Message{}, err
	}

	var m Message
	if err := json.Unmarshal(body, &m); err != nil {
		return Message{}, err
	}
	return m, nil
}
