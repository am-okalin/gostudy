package basic

import (
	"fmt"
	"testing"
)

type data struct {
	name string
}

//go:noinline
func e1() data {
	d := data{"downing"}
	return d
}

//go:noinline
func e2() *data {
	d := data{"downing"}
	fmt.Printf("%p\n", &d)
	return &d
}

//todo::执行命令做逃逸分析
func TestEscape(t *testing.T) {
	d1 := e1()
	fmt.Printf("%p\n", &d1)
	d2 := e2()
	fmt.Printf("%p\n", d2)
	t.Log(d1.name + d2.name)
}
