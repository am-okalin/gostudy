package container

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestHeap(t *testing.T) {
	h := &IntHeap{}
	h.Push(5)
	h.Push(2)
	h.Push(1)
	h.Push(8)
	h.Push(7)
	h.Push(6)
	h.Push(2)
	h.Push(4)
	fmt.Println(h)
	heap.Init(h)
	fmt.Println(h)
	b := heap.Remove(h, 1)
	fmt.Println(h)
	fmt.Println(b)
	heap.Push(h, 12)
	fmt.Println(h)
	heap.Pop(h)
	fmt.Println(h)
}