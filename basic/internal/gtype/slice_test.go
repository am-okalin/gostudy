package gtype

import (
	"testing"
)

func TestArr(t *testing.T) {
	//数组 定长不可变
	var arr1 [10]int
	for i := 0; i < 10; i++ {
		arr1[i] = 100 - i
	}

	for i, v := range arr1 {
		t.Log(i, v)
	}

	t.Log(arr1)

	t.Log(len(arr1), cap(arr1))

	arr2 := [10]int{1, 2, 3, 4, 5, 6, 7}
	t.Log(arr2)
}

func TestSlice(t *testing.T) {
	var slice1 []int
	for i := 0; i < 10; i++ {
		//当append后的数量超过cap后，会对底层数组进行扩容
		//扩容规则为之前容量的2倍，当容量大于1000时，则扩容为1/4
		slice1 = append(slice1, 100-i)
	}

	for i, v := range slice1 {
		t.Log(i, v)
	}

	t.Log(len(slice1), cap(slice1))

	slice2 := []int{1, 2, 3, 4, 5, 6, 7}
	t.Log(slice2)
}

func TestEmptyNil(t *testing.T) {
	var a []int
	b := new([]int)
	c := []int{}
	//make(type, len, cap) make进行初始化的操作
	d := make([]int, 2, 10)

	// a,b没有分配空间 所以为nil
	// c,d分配了空间，实际上为空的结构
	t.Log(a, *b, c, d)
	t.Log(len(a), len(*b), len(c), len(d))

	if a == nil {
		t.Log("a is nil")
	}
	if *b == nil {
		t.Log("b is nil")
	}
	if c == nil {
		t.Log("c is nil")
	}
	if d == nil {
		t.Log("d is nil")
	}
}

func TestExpand(t *testing.T) {
	slice := make([]int, 0, 0)
	//当append后的数量超过cap后，会对底层数组进行扩容
	slice = append(slice, 1, 2, 3, 4, 5, 6, 7)
	t.Log(len(slice), cap(slice))

	//扩容规则为之前容量的2倍，当容量大于1000时，则扩容为1/4
	slice = append(slice, 8, 9)
	t.Log(len(slice), cap(slice))
}

func TestAppend(t *testing.T) {
	//切片 指向数组的指针，可扩容;
	var b []int
	a := []int{1, 2, 3}
	b = append(b, a[0], a[1], a[2])
	b = append(b, a...)
	t.Log(b)
}

func TestSliceIntercept(t *testing.T) {
	a := make([]int, 0, 10)
	a = append(a, 0)
	a = append(a, 1)
	t.Log(a) // [0, 1]
	//左闭右开区间，左边包括，右边不包括
	//a[start:end]
	a = a[:1] // [0]
	t.Log(a)
	a = a[:0] // []
	t.Log(a)
}

func TestCopy(t *testing.T) {
	a := []int{0, 1}
	a = append(a, 2)

	// 扩容也会深拷贝
	b := append(a, 4)
	//t.Log(b, a[:4])
	c := append(a, 5)

	// 拷贝前要设置好length
	//b := make([]int, 2)
	//b := make([]int, len(a))
	//c := make([]int, len(a), len(a)+1)
	//copy(b, a)
	//copy(c, a)
	//b = append(b, 4)
	//c = append(c, 5)

	t.Log(a, b, c)
}
