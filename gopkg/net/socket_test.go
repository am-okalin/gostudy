package net

import (
	"net"
	"syscall"
	"testing"
	"time"
)

func TestSocket(t *testing.T) {
	var domain int // [AF_INET ipv4, AF_INET6 ipv6, AF_UNIX unix]
	var stype int  // [SOCK_DGRAM 无逻辑连接, SOCK_STREAM 无消息边界, SOCK_SEQPACKET, SOCK_RAW]
	var proto int  // 为0时会根据前2个参数自行选择
	socket, err := syscall.Socket(domain, stype, proto)
	t.Log(socket, err)

	// proto=0 会自动选择 udp
	socket, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	t.Log(socket, err)

	// 大概率选择 tcp
	socket, err = syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	t.Log(socket, err)
}

func TestDial(t *testing.T) {
	// 解析network, address
	// 创建sockethi里并建立网络连接
	conn, err := net.Dial("tcp6", "127.0.0.1")
	t.Log(conn, err)

	// 解析address时,可能访问DNS, 有可能解析出多个IP,采用第一个建立的IP
	conn, err = net.DialTimeout("tcp6", "127.0.0.1", time.Second)
	t.Log(conn, err)

	// 对dial正确的设置读写的超时时间

}
