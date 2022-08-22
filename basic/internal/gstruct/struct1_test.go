package gstruct

import "testing"

func TestStruct1(t *testing.T) {
	a := NewStruct1(10)
	t.Log(a)
	a.SetParam1("字符串q")
	t.Log(a)
	a.SetParam1N("字符串w")
	t.Log(a)
}
