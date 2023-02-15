package tableconv

import "testing"

type people struct {
	name string
	age  uint32
	sex  string
}

func TestObjGroup(t *testing.T) {
	list := []*people{
		{name: "Downing", age: 26, sex: "male"},
		{name: "Tony", age: 20, sex: "male"},
		{name: "Vela", age: 28, sex: "female"},
	}

	//todo处理从泛型获取实际类型, 并对其赋值
	ObjGroup(list, "sex")
}
