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
	var chs1 []chan interface{}
	//会在第一次扩容时初始化chs1
	for i := 0; i < num; i++ {
		chs1 = append(chs1, make(chan interface{}))
	}
	t.Log(len(chs1), cap(chs1))
	//初始化双向[]chan
	chs2 := make([]chan interface{}, num)
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
	for i := 0; i < 5; i++ {
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
			case v := <-ch:
				fmt.Println("for-select:", v)
			}
		}
		//todo::本协程未退出
		fmt.Println("for select done")
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
		for i := 0; i < 5; i++ {
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
	go func() {
		for i := 0; i < num; i++ {
			ch <- fmt.Sprintf("a%d", i)
			t.Logf("in:a%d", i)
			time.Sleep(time.Second)
		}
		close(ch)
	}()

	go func() {
		for i := 0; i < num; i++ {
			v := <-ch
			t.Logf("out:%s", v)
			//time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Minute)
}

func TestCh3(t *testing.T) {
	num := 10
	sign := make(chan struct{}, num)

	//猜猜输出什么
	for i := 0; i < num; i++ {
		go func() {
			//子goruntine在运行到这一行时才拿i的值进行输出
			fmt.Println(i)
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
