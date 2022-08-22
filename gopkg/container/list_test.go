package container

import (
	"container/list"
	"testing"
)

func TestList(t *testing.T) {
	//var l list.List //返回结构体空值(使用时延迟初始化)
	var l = list.New() // 返回引用类型(立刻初始化)
	v1 := l.PushBack(1)
	v2 := l.PushFront(2)
	v3 := l.InsertAfter(3, v2)
	v4 := l.InsertBefore(4, v1)
	l.MoveToFront(v1)
	l.MoveAfter(v2, v1)
	front := l.Front()
	backVal := l.Remove(l.Back())
	length := l.Len()
	l.Init() // 初始化为空链表
	t.Log(v1.Value,v2.Value,v3.Value,v4.Value)
	t.Log(front.Value,backVal,length)
}