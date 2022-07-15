package gstruct

import (
	"fmt"
	"testing"
)

//Test1 以测试方法作为客户端代码
func Test1(t *testing.T) {
	//用构造方法实例化A
	a := NewA(1, "Downing")
	//聚合，将a注入b
	b := NewB(a)
	b.AInterface.SayHi("")
	//类型断言
	b.AInterface.(*A).SetSayStr("Gary")
	//语法糖-隐式调用AInterface的方法
	b.SayHi("Gary")
	//按string接口的实现输出b
	fmt.Println(b)
	fmt.Println(a)
}
