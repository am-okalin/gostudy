package gochan

import (
	"fmt"
	"testing"
	"time"
)

func process(timeout time.Duration) bool {
	// 将此chan改为带缓存的chan就可避免
	ch := make(chan bool)

	go func() {
		// 模拟处理耗时的业务
		time.Sleep((timeout + time.Second))
		ch <- true // block(阻塞) 导致goroutine泄露
		fmt.Println("exit goroutine")
	}()
	select {
	// 超时时永远无法执行到此行
	case result := <-ch:
		return result
	case <-time.After(timeout):
		return false
	}
}

//模拟chan导致的goroutine泄露
func TestLeak(t *testing.T) {
	process(2 * time.Second)
}
