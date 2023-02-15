package net

import (
	"net/url"
	"testing"
)

const (
	url1   = "https://admin:password@www.baidu.com:80/search?mq=test#12345"
	query1 = "mq=test&queue=people"
)

func TestUrl(t *testing.T) {
	//不解析#后缀(网页浏览器会在去掉该后缀后才将网址发送到网页服务器)
	u2, err := url.ParseRequestURI(url1)
	if err != nil {
		t.Error(err)
	}
	t.Log(u2)

	//解析#后缀
	u1, err := url.Parse(url1)
	if err != nil {
		t.Error(err)
	}

	//param解析(要去除'?')
	q1 := u1.Query()
	q1.Add("mq", "rabbit")
	u1.RawQuery = q1.Encode()
	uri1 := u1.RequestURI()
	t.Log(u1, uri1)

	//解析query字符串
	q2, err := url.ParseQuery(query1)
	if err != nil {
		t.Error(err)
	}
	q2.Add("mq", "rabbit")
	q2s := q2.Encode() //mq=test&mq=rabbit&queue=people
	t.Log(q1, q2, q2s)
}
