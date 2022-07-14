package gtype

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
)

type Cat struct {
	Name  string
	Color string `json:"cat_color;"`
	age   int
}

func CatString() string {
	cat := Cat{
		Name:  "小四",
		Color: "三花",
		age:   1,
	}

	//json序列化
	marshal, err := json.Marshal(cat)
	if err != nil {
		//panic会中断程序执行
		panic(err)
	}

	//类型强转
	return string(marshal)
}

func TestCatString(t *testing.T) {
	//defer在程序退出时执行 由panic导致的异常程序退出也会执行defer
	defer func() {
		if err := recover(); err != nil {
			// 这里的err其实就是panic传入的内容
			fmt.Println(err)
		}
	}()

	//panic("downing test panic")

	str := CatString()
	t.Log(str)
}

//runtime.Stack(buf, all),buf是保存信息的字符串数字
//all表示是否要打印全部的堆栈信息。all=true就会先调用stopTheWorld
func TestGetStackInfo(t *testing.T) {
	var (
		size = 64 * 1024
		all  = false
	)
	buf := make([]byte, size)
	runtime.Stack(buf, all)
	t.Logf("%s", string(buf))
}

//TestStack panic时会打印出goroutine的ID, 状态，函数，调用栈的信息。
func TestStack(t *testing.T) {
	panic("get stack info")
}
