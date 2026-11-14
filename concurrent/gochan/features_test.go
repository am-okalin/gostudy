package gochan

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// chan的初始化
func TestChInit(t *testing.T) {
	num := 3
	//创建仅接收chan(值为nil,不可使用)
	//ch1 := new(chan<- struct{})
	var ch1 chan<- struct{}
	//ch1 <- struct{}{} 	//给为nil的channel发送数据 会造成永远阻塞
	t.Log(ch1)

	//初始化仅发送chan(值为chan类型的空值)
	ch2 := make(<-chan struct{}, num)
	// len返回缓存中的元素数量， cap返回缓存容量
	t.Log(ch2, len(ch2), cap(ch2))
	//创建双向[]chan
	var chs1 []chan any
	//会在第一次扩容时初始化chs1
	for range num {
		chs1 = append(chs1, make(chan any))
	}
	t.Log(len(chs1), cap(chs1))
	//初始化双向[]chan
	chs2 := make([]chan any, num)
	t.Log(len(chs2), cap(chs2))
}

// todo::优雅退出select
func TestSelect(t *testing.T) {
	//select {
	//case <-closed:
	//case <-time.After(time.Second):
	//	fmt.Println("goroutine cleanup time out")
	//}
	//fmt.Println("goroutine cleanup success")
}

// chan使用for-select for-range
func TestChFor(t *testing.T) {
	// 使用select
	var ch = make(chan int, 10)
	for i := range 5 {
		select {
		case ch <- i:
			fmt.Println("in : ", i)
		case v := <-ch:
			fmt.Println("out: ", v)
		}
	}

	//for-select循环获取元素，当ch为空时阻塞
	go func() {
		for {
			select {
			case v, ok := <-ch:
				// 通道关闭后若只写 v := <-ch，会一直读到零值形成忙循环
				if !ok {
					return
				}
				fmt.Println("for-select:", v)
			}
		}
	}()

	// for-range与for-select功能相等
	go func() {
		time.Sleep(time.Second)
		//可用 for range ch {} 清空元素后阻塞
		for v := range ch {
			fmt.Println("for-range :", v)
		}
		fmt.Println("for range done")
	}()

	//往两个goruntine中投递数据
	for i := 5; i < 10; i++ {
		ch <- i
	}
	close(ch)

	//延迟2秒后退出主进程
	time.Sleep(time.Second * 2)
	fmt.Println("done")
}

// chan_close的特性
func TestChClose(t *testing.T) {
	//部分chan操作会导致panic
	defer func() {
		//close 为 nil 的 chan
		//send 已经 close 的 chan
		//close 已经 close 的 chan。
		p := recover()
		t.Log("panic:", p)
	}()

	//chan_close后 阻塞的<-ch 会返回 其传输数据类型 的 零值
	ch1 := make(chan struct{}, 5)
	go func() {
		time.Sleep(1 * time.Second)
		close(ch1)
	}()
	go func() {
		t.Log(<-ch1) //阻塞至close(ch1)执行
	}()

	//chan_close后的每次执行 <-ch 会按顺序返回队列中的数据，当队列空后返回零值
	time.Sleep(2 * time.Second)
	ch2 := make(chan int, 5)
	go func() {
		ch2 <- 2
		ch2 <- 1
		close(ch2)
	}()
	go func() {
		time.Sleep(1 * time.Second)
		//当队列空后返回零值
		for range 5 {
			t.Log(<-ch2)
		}
	}()

	//todo::测试forrange的清理逻辑

	time.Sleep(3 * time.Second)
}

func TestCh1(t *testing.T) {
	//通道作为`传递某种简单信号`的`介质`时，用struct{}最合适
	//ch := new(chan struct{})
	ch := make(chan struct{})
	go func() {
		<-ch
		fmt.Println("子goruntine结束")
	}()
	//struct{}类型`值的表示法`只有一个即 struct{}{}, 占用内存空间=0字节
	//这个值在整个`Go程序`中永远都只会存在一份。虽然我们可以无数次地使用这个`值字面量`，但是用到的却都是同一个值。
	ch <- struct{}{}
	fmt.Println("父goruntine结束")
}

func TestCh2(t *testing.T) {
	num := 10
	ch := make(chan string, num)
	done := make(chan struct{})
	go func() {
		for i := range num {
			ch <- fmt.Sprintf("a%d", i)
			t.Logf("in:a%d", i)
			time.Sleep(time.Second)
		}
		close(ch)
	}()

	go func() {
		defer close(done)
		for range num {
			v := <-ch
			t.Logf("out:%s", v)
			//time.Sleep(time.Second)
		}
	}()

	// 等消费者读完，避免固定 sleep 造成的等待或竞态
	<-done
}

func TestCh3(t *testing.T) {
	num := 10
	sign := make(chan struct{}, num)

	// 猜猜输出什么：go1.22 起循环变量每次迭代都会重新声明，因此输出 0..9（1.21 及之前通常全是 10）
	for i := range num {
		go func() {
			// 子 goroutine 运行时读取的是本次迭代的 i（go1.22+）
			fmt.Println(i)
			sign <- struct{}{}
		}()
	}

	//阻塞进程的方式有 select for-range time.sleep sync.wait chan
	//select {}
	//time.Sleep(time.Millisecond * 500)
	for range num {
		<-sign
	}
}

func TestCh4(t *testing.T) {
	// 所有 goroutine 共享同一个计数器，才能按 0..9 的顺序依次触发
	var count atomic.Uint32
	for i := range uint32(10) {
		// 如此的话会先对子goruntine的形参会先求值
		go func(i uint32) {
			fn := func() {
				fmt.Println(i)
			}
			trigger(&count, i, fn)
		}(i)
	}
	trigger(&count, 10, func() {})
}

func trigger(count *atomic.Uint32, i uint32, fn func()) {
	for {
		if n := count.Load(); n == i {
			fn()
			count.Add(1)
			return
		}
		time.Sleep(time.Nanosecond)
	}
}
