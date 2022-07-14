package ants

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"runtime"
	"testing"
	"time"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	// GiB // 1073741824
	// TiB // 1099511627776             (超过了int32的范围)
	// PiB // 1125899906842624
	// EiB // 1152921504606846976
	// ZiB // 1180591620717411303424    (超过了int64的范围)
	// YiB // 1208925819614629174706176
)

const (
	Param    = 100
	AntsSize = 1000
	TestSize = 10000
	n        = 100000
)

var curMem uint64

func demoFunc() {
	time.Sleep(time.Duration(10) * time.Millisecond)
}

func demoPoolFunc(args interface{}) {
	n := args.(int)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func Test1(t *testing.T) {
	p, _ := ants.NewPool(10)
	defer p.Release()

	for i := 0; i < 1000; i++ {
		_ = p.Submit(func() {
			time.Sleep(time.Second)
			fmt.Println("do submit")
		})
	}

	time.Sleep(time.Minute)
}

// TestAntsPoolGetWorkerFromCache is used to sql getting worker from sync.Pool.
func TestAntsPoolGetWorkerFromCache(t *testing.T) {
	p, _ := ants.NewPool(AntsSize)
	defer p.Release()

	for i := 0; i < TestSize; i++ {
		_ = p.Submit(demoFunc)
	}
	time.Sleep(2 * ants.DefaultCleanIntervalTime)
	_ = p.Submit(demoFunc)
	t.Logf("pool, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestAntsPoolWithFuncGetWorkerFromCachePreMalloc(t *testing.T) {
	dur := 10
	p, _ := ants.NewPoolWithFunc(TestSize, demoPoolFunc, ants.WithPreAlloc(true))
	defer p.Release()

	for i := 0; i < AntsSize; i++ {
		_ = p.Invoke(dur)
	}
	time.Sleep(2 * ants.DefaultCleanIntervalTime)
	_ = p.Invoke(dur)
	t.Logf("pool with func, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}
