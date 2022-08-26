package basic

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDeclaration(t *testing.T) {
	var str1 string     //str1的空值
	str2 := ""          //隐式的指定了类型为string
	str3 := new(string) //返回指针
	t.Log(str1, str2, str3, *str3)
}

func TestPointer(t *testing.T) {
	f1 := func(s *string) {
		*s = "f1"
	}
	f2 := func(s string) {
		s = "f2"
	}

	str := "something"

	//传入地址
	f1(&str)
	t.Log(str)

	//传入副本
	f2(str)
	t.Log(str)
}

func TestByte(t *testing.T) {
	var a []byte
	b := new([]byte)
	c := make([]byte, 8)

	t.Log(a, *b, c, nil)

	if a == nil {
		t.Log("[]byte == nil")
	}
	if *b == nil {
		t.Log("&[]byte == nil")
	}
	if c == nil {
		t.Log("c == nil")
	}
}

func TestStr(t *testing.T) {
	jsonStr := "[\"str1\", \"str2\"]" //json字符串 ["str1", "str2"]
	var list []string

	err := json.Unmarshal([]byte(jsonStr), &list)
	t.Log(err, list)
}

func TestInt(t *testing.T) {
	var ui32 uint32
	t.Log(ui32, ui32-1)

	i64 := time.Now().UnixMilli()
	t.Log(i64, uint32(i64))
}
