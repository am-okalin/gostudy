package basic

import "testing"

type Struct1 struct {
	param1 string
	Param2 int
}

func NewStruct1(param2 int) *Struct1 {
	return &Struct1{
		param1: "内部赋值",
		Param2: param2,
	}
}

func (a Struct1) Param1() string {
	return a.param1
}

//SetParam1 设置变量时的`receiver`必须为指针类型
func (a *Struct1) SetParam1(param1 string) {
	a.param1 = param1
}

func (a Struct1) SetParam1N(param1 string) {
	a.param1 = param1
}

func TestStruct1(t *testing.T) {
	a := NewStruct1(10)
	t.Log(a)
	a.SetParam1("字符串q")
	t.Log(a)
	a.SetParam1N("字符串w")
	t.Log(a)
}
