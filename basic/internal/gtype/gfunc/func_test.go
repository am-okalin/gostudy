package gfunc

import (
	"testing"
)

func Test1(t *testing.T) {
	f1()
}

func Test2(t *testing.T) {
	f1 := getEatFruitFun()
	w1 := f1("apple")
	t.Log(w1)

	f2 := getEatMeatFun()
	w2 := f2("Beef")
	t.Log(w2)
}
