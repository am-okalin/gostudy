package basic

import (
	"errors"
	"fmt"
	"testing"
)

//哨兵
var ErrNotFound = errors.New("not found")

func doSomethingPanic() string {
	//panic会中断程序执行，通常传入error类型的参数
	panic(ErrNotFound)
}

func TestPanic(t *testing.T) {
	//defer在程序退出时执行 由panic导致的异常程序退出也会执行defer
	defer func() {
		//这个↓err其实就是panic传入的内容
		if err := recover(); err != nil {
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
