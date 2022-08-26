package basic

import (
	"fmt"
	"testing"
)

//AInterface 定义A接口
type AInterface interface {
	//SayHi 根据实现输出str
	SayHi(str string)
}

//A 实现AInterface的结构体
type A struct {
	//uni A的唯一标识
	uni int64
	//sayStr 要说的内容
	sayStr string
}

//NewA 构造方法
func NewA(uni int64, sayStr string) *A {
	return &A{uni: uni, sayStr: sayStr}
}

//SetSayStr 设置要说的内容
//此时的a是带指针的，若改为不带指针，请问可以设置成功吗？
//func (a Struct1) SetSayStr(sayStr string) {
func (a *A) SetSayStr(sayStr string) {
	a.sayStr = sayStr
}

//SayHi 输出sayStr
func (a A) SayHi(str string) {
	if str == "" {
		str = a.sayStr
	}

	fmt.Println(str)
}

//String
func (a A) String() string {
	return fmt.Sprintf("%d: will say %s", a.uni, a.sayStr)
}

type B struct {
	AInterface
	//Struct1	//常用于继承
	//AInterface AInterface //常用于组合
}

//NewB 聚合
func NewB(AInterface AInterface) *B {
	return &B{AInterface: AInterface}
}

//NewB 组合
//func NewB() *B {
//	return &B{AInterface: NewStruct1(100, "inline-build")}
//}

//Test1 以测试方法作为客户端代码
func TestStruct2(t *testing.T) {
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
