package list

import (
	"testing"
)

func initList() *List {
	v3 := Elem{Value: "v3"}
	v2 := Elem{Value: "v2", next: &v3}
	//v2 := Elem{Value: "v2"}
	v1 := Elem{Value: "v1", next: &v2}
	v0 := Elem{Value: "v0", next: &v1}
	l := List{head: &v0}
	return &l
}

func Test1(t *testing.T) {
	l := initList()
	v1 := l.get(-1)
	flag := l.removeAt(v1)
	t.Log(flag)
	t.Log(v1)
	t.Log(l)
}

func TestReverse(t *testing.T) {
	l := initList()
	t.Log(l)
	l.reverse()
	t.Log(l)
}

func TestMiddle(t *testing.T) {
	l := initList()
	if l.hasCycle() {
		t.Error("该链表有环")
	}
	t.Log(l)
	e := l.middle()
	t.Log(e)
}
