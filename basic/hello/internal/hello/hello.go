package hello

import "fmt"

//A 公有变量
var A = "string_A"
var a = "string_a"

//SayHello 输出hello world %s
//小写表示包级私有，外部不可调用
func sayHello(name string) {
	fmt.Printf("hello %s\n", name)
}

//SayHello 输出hello world %s
//internal目录下的包仅其父级可访问
func SayHello(name string) {
	sayHello(name)
}
