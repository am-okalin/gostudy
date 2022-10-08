package sync

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestCAS(t *testing.T) {

}

//Test01 猜猜输出啥
func Test01(t *testing.T) {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("go1: ", i) //这是一次io操作
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("go2: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
