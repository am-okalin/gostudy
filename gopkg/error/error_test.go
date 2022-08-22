package error

import (
	"fmt"
	"testing"
)

func doSomethingPanic() string {
	//panic会中断程序执行
	panic(ErrNotFound)
}

func TestPanic(t *testing.T) {
	//defer在程序退出时执行 由panic导致的异常程序退出也会执行defer
	defer func() {
		if err := recover(); err != nil {
			// 这里的err其实就是panic传入的内容
			fmt.Println(err)
		}
	}()

	//panic("downing test panic")
	doSomethingPanic()

	fmt.Println("hi")
	if err := recover(); err != nil {
		// 这里的err其实就是panic传入的内容
		fmt.Println(err)
	}
}
