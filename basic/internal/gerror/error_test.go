package gerror

import (
	"errors"
	"fmt"
	"testing"
)

var ErrorNotFound = errors.New("record not found")

func doSomethingPanic() string {
	//panic会中断程序执行
	panic(ErrorNotFound)
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
	doSomethingPanic()
}
