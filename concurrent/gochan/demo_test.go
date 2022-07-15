package gochan

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCh3(t *testing.T) {
	num := 10
	//通道作为`传递某种简单信号`的`介质`时，用struct{}最合适
	sign := make(chan struct{}, num)

	//猜猜输出什么
	for i := 0; i < num; i++ {
		go func() {
			//子goruntine在运行到这一行时才拿i的值进行输出
			fmt.Println(i)
			//struct{}类型`值的表示法`只有一个即 struct{}{}, 占用内存空间=0字节
			//这个值在整个`Go程序`中永远都只会存在一份。虽然我们可以无数次地使用这个`值字面量`，但是用到的却都是同一个值。
			sign <- struct{}{}
		}()
	}

	//阻塞进程的方式有 select for-range time.sleep sync.wait chan
	//select {}
	//time.Sleep(time.Millisecond * 500)
	for j := 0; j < num; j++ {
		<-sign
	}
}

func TestCh4(t *testing.T) {
	for i := uint32(0); i < 10; i++ {
		// 如此的话会先对子goruntine的形参会先求值
		go func(i uint32) {
			fn := func() {
				fmt.Println(i)
			}
			trigger(i, fn)
		}(i)
	}
	trigger(10, func() {})
}

func trigger(i uint32, fn func()) {
	var count uint32
	for {
		if n := atomic.LoadUint32(&count); n == i {
			fn()
			atomic.AddUint32(&count, 1)
			break
		}
		time.Sleep(time.Nanosecond)
	}
}

//Test5 猜猜输出啥
func Test5(t *testing.T) {
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
