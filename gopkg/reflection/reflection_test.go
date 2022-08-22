package reflection

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

// 定义一个Enum类型
type Enum int

const (
	Zero Enum = 0
)

func TestTypeOf(t *testing.T) {
	var w io.Writer = os.Stdout    // reflect.TypeOf总是返回具体的类型
	fmt.Println(reflect.TypeOf(w)) // "*os.File"
	fmt.Printf("%T\n", w)

	to := reflect.TypeOf(3)
	t.Log(to)
}

func TestValueOf(t *testing.T) {
	//v := reflect.ValueOf(3)

	v := reflect.ValueOf(3)              // reflect.Value
	fmt.Println(v)                       // "3"
	fmt.Printf("%v\n", v.Interface())    // "3"
	fmt.Println(v.Kind() == reflect.Int) //true
	fmt.Println(v.String())              // NOTE: "<int Value>"
}

func TestDemo2(t *testing.T) {
	defer func() {
		p := recover()
		t.Log("panic:", p)
	}()

	var x float64 = 3.4
	// (v Value) Elem() 将指针指向的值解析到v中
	v1 := reflect.ValueOf(&x).Elem()         //传递x的地址后解引用
	v2 := reflect.ValueOf(x)                 //传递x的拷贝
	t.Log("settability of v1:", v1.CanSet()) //true
	t.Log("settability of v2:", v2.CanSet()) //false
	v1.SetFloat(7.1)                         //ok
	t.Log("v‘s value is", v1)
	v2.SetFloat(7.1) //Error: will panic
}

func TestAny(t *testing.T) {
	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	fmt.Println(Any(x))          // "1"
	fmt.Println(Any(d))          // "1"
	fmt.Println(Any([]int64{x})) // "[]int64 0x8202b87b0"
	fmt.Println(Any([]time.Duration{d}))
}

func TestSelect1(t *testing.T) {
	Select1()
}

func TestSelect2(t *testing.T) {
	Select2()
}

//func TestStruct(t *testing.T) {
//	err := toMap()
//	t.Log(err)
//}
