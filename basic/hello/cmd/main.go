package main

import (
	"basic/hello/internal/hello"
)

//main 命令源码文件：必定包含main方法
func main() {
	name := "downing"
	hello.SayHello(name)
	//fmt.Println(hello.A)
}
