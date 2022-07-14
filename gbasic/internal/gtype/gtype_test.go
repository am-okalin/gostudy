package gtype

import "testing"

func TestDeclaration(t *testing.T) {
	var str1 string
	str2 := ""          //隐式的指定了string类型
	str3 := new(string) //返回指针类型
	t.Log(str1, str2, str3, *str3)
}

func TestFunc(t *testing.T) {
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
