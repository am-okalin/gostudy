package reflection

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

// TestDemo 反射的基本使用
func TestDemo(t *testing.T) {
	var w io.Writer = os.Stdout
	wt := reflect.TypeOf(w)
	wv := reflect.ValueOf(w)
	//wvt := reflect.ValueOf(w).Type() //等同于wt

	fmt.Printf("%T\n", w) // "*os.File"
	fmt.Println(wt)       // reflect.TypeOf返回具体的类型

	fmt.Printf("%v\n", wv.Interface()) // "&{0xc0000c6280}"
	fmt.Println(wv)                    // reflect.Value返回具体的值(若是指针类型,则返回其指向的地址)
	fmt.Println(wv.String())           // "<*os.File Value>"

	fmt.Println(wv.Kind() == reflect.Int) //true
}

// TestStruct 反射中对结构体的解析
func TestStruct(t *testing.T) {
	s := User{Name: "user1"}
	st := reflect.TypeOf(s)

	//获取对象中的属性信息
	stf0 := st.Field(0)
	stfAge, ok := st.FieldByName("Age")
	if !ok {
		t.Error("Age 属性不存在")
	}
	fmt.Println(st.Name(), stf0, stfAge)

	//获取属性相关的Tag信息
	color := stf0.Tag.Get("color")
	species := stf0.Tag.Get("species")
	fmt.Println(color, species)
}

// TestSet 反射中对变量的重新赋值
func TestSet(t *testing.T) {
	var err error
	u := User{Name: "user1"}
	uv := reflect.ValueOf(u)           //传递u的拷贝
	uvpr := reflect.ValueOf(&u).Elem() //传递u的地址后解引用

	//值变量传递了拷贝无法被寻址，所以set失败
	err = setValue(uv, reflect.ValueOf(User{Name: "user2"}))
	if err != nil {
		t.Log(err)
	}

	//成功
	err = setValue(uvpr, reflect.ValueOf(User{Name: "user2"}))
	if err != nil {
		t.Log(err)
	}

	//设置的值类型不同，所以set失败
	err = setValue(uvpr, reflect.ValueOf(1))
	if err != nil {
		t.Log(err)
	}
}
