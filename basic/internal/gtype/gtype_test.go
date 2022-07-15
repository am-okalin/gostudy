package gtype

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDeclaration(t *testing.T) {
	var str1 string
	str2 := ""          //隐式的指定了string类型
	str3 := new(string) //返回指针类型
	t.Log(str1, str2, str3, *str3)
}

func TestStr(t *testing.T) {
	a := "[\"1\", \"2\"]"
	t.Log(a)
	var s []string

	err := json.Unmarshal([]byte(a), &s)
	t.Log(err, s)
}

func TestInt(t *testing.T) {
	var a uint32
	a = a - 1
	t.Log(a)

	b := time.Now().UnixMilli()
	c := uint32(b)
	t.Log(b, c)
}

func TestPointer(t *testing.T) {
	f1 := func(s *string) {
		*s = "q"
	}
	f2 := func(s string) {
		s = "w"
	}

	str := "a"
	//传入地址
	f1(&str)
	t.Log(str)
	//传入副本
	f2(str)
	t.Log(str)
}
