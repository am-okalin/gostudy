package ginterface

import (
	"fmt"
	"testing"
)

type A struct {
	Name string
	Age  uint
}

//func (a A) String() string {
//	return fmt.Sprintf("my name is %s", a.Name)
//}

func TestStringer(t *testing.T) {
	a := A{
		Name: "downing",
		Age:  10,
	}
	fmt.Println(a)
	fmt.Printf("%v\n", a)
}
