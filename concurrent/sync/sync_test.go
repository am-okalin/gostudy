package sync

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestCAS(t *testing.T) {

}

// Test01 猜猜输出啥
func Test01(t *testing.T) {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	// go1.22 起循环变量每次迭代都会重新声明，go1 组输出 0..9，与 go2 组一致（1.21 及之前 go1 通常全是 10）
	for i := range 10 {
		go func() {
			fmt.Println("go1: ", i) //这是一次io操作
			wg.Done()
		}()
	}
	for i := range 10 {
		go func(i int) {
			fmt.Println("go2: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
