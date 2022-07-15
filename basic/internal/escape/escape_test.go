package escape

import (
	"testing"
)

type data struct {
	name string
}

//go:noinline
func f1() data {
	d := data{"downing"}
	return d
}

//go:noinline
func f2() *data {
	d := data{"downing"}
	return &d
}

//todo::执行命令做逃逸分析
func TestEscape(t *testing.T) {
	d1 := f1()
	d2 := f2()
	t.Log(d1.name + d2.name)
}
