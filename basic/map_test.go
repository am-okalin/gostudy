package basic

import (
	"sync/atomic"
	"testing"
)

func TestNewMap(t *testing.T) {
	// 赋空值初始化
	// new 返回指针 指向所有字段全为0的值map对象 需要初始化后才能使用
	ma := new(map[string]int)
	//(*ma)["a"] = 1 //报错
	*ma = map[string]int{}
	(*ma)["a"] = 1
	t.Log(ma)

	// make初始化
	mb := make(map[string]int, 0)
	mb["c"] = 3
	// 赋值初始化(mb与ma操作同一个底层hash)
	*ma = mb
	t.Log(ma, mb)
	// 删除map的键
	delete(*ma, "c")
	t.Log(ma, mb)
}

func TestCopyM(t *testing.T) {
	a := map[int]bool{0: true, 1: true}

	// 浅拷贝，共用底层
	b := a
	b[2] = true
	t.Log(a, b)

	// 深拷贝，不共用底层
	c := make(map[int]bool)
	for i, f := range a {
		c[i] = f
	}
	c[3] = true
	t.Log(a, b, c)
}

func TestMM(t *testing.T) {
	//二维映射
	mm := make(map[int]map[int]float64)
	mm[1] = make(map[int]float64)
	mm[1][1] = 100.1
	//判断元素是否存在
	f1, ok1 := mm[1][1] //100.1 true
	f2, ok2 := mm[1][2] //0 flase
	t.Log(f1, ok1)
	t.Log(f2, ok2)
}

func TestAtomicAdd(t *testing.T) {
	m := make(map[string]*uint64)
	//报错: m["a"]的值是nil
	atomic.AddUint64(m["a"], 1)
	t.Log(m)
}
