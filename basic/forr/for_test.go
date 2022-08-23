package forr

import (
	"fmt"
	"strconv"
	"testing"
)

type Item struct {
	Value1 string
	Value2 int
}

func getList(n int) []Item {
	fmt.Println("getList()")

	data := make([]Item, n)
	for i := 0; i < n; i++ {
		data[i].Value1 = strconv.Itoa(i) + "-string"
		data[i].Value2 = i
	}
	return data
}

func getLen() int {
	fmt.Println("getLen()")
	return 5
}

func TestFor(t *testing.T) {
	// forr range仅执行一次getList
	for i, item := range getList(5) {
		t.Log(i, item.Value1, item.Value2)
	}

	// fori 每次迭代都要执行getLen
	// 所以getLen要用形参计算出来后再放入迭代
	for i := 0; i < getLen(); i++ {
		t.Log(i)
	}
}

func Test1(t *testing.T) {
	for i, item := range getList(5) {
		t.Log(i, item.Value1, item.Value2)
	}
	//仅遍历计数
	for i := range getList(5) {
		t.Log(i)
	}
}
