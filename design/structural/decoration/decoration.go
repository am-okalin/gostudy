package decoration

import (
	"fmt"
)

type 通知 interface {
	echosomething() map[string]bool
}

type 应用 struct {
	param string
}

func (i *应用) echosomething() map[string]bool {
	//var m map[string]bool
	m := make(map[string]bool)
	m["应用"] = true
	return m
}

type 短信 struct {
	param    string
	instance 通知
}

func (i *短信) echosomething() map[string]bool {
	m := i.instance.echosomething()
	m["短信"] = true
	return m
}

func main() {
	a := 应用{param: "这是参数"}
	b := 短信{param: "这是参数", instance: &a}
	m := b.echosomething()
	fmt.Printf("%v", m)
}
