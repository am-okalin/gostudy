package net

import (
	"bufio"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

func TestHttpGet(t *testing.T) {
	url := "http://127.0.0.1:8081/hi"
	var httpClient1 http.Client
	resp, err := httpClient1.Get(url)
	//resp1, err := http.Get(url) // 等价于这一行
	if err != nil {
		t.Logf("request sending error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 取出返回body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Logf("request sending error: %v\n", err)
		return
	}

	t.Logf("The first line of response:\n%s %s\n%q", resp.Proto, resp.Status, body)
}

func Test1(t *testing.T) {
	reqStr := `GET /hi HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: 127.0.0.1:8081
User-Agent: Dialer/go1.18.10

`

	// 创建请求连接
	conn1, err := net.Dial("tcp4", "127.0.0.1:8081")
	if err != nil {
		t.Log(err)
		return
	}
	defer conn1.Close()

	// 写入请求参数
	_, err = io.WriteString(conn1, reqStr)
	if err != nil {
		t.Log(err)
		return
	}

	// 获取返回信息
	reader1 := bufio.NewReader(conn1)
	//body, err := io.ReadAll(reader1) // 这样取数据流, 接口返回失败, 这是为啥?
	body, err := reader1.ReadString('\n')

	t.Log(body, err)
}
