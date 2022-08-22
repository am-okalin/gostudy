package gochan

import (
	"fmt"
	"testing"
	"time"
)

// chan的初始化
func TestCh0(t *testing.T) {
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

// chan使用for-select for-range
func TestCh1(t *testing.T) {
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
	}()

	// for-range与for-select功能相等
	go func() {
		//可用 for range ch {} 清空元素后阻塞
		for v := range ch {
			fmt.Println("for-range :", v)
		}
	}()

	//往两个goruntine中投递数据
	for i := 5; i < 10; i++ {
		ch <- i
	}

	//延迟2秒后退出主进程
	time.Sleep(time.Second * 2)
	fmt.Println("done")
}

//chan_close的特性
func TestCh2(t *testing.T) {
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

	time.Sleep(3 * time.Second)
}

func Test1(t *testing.T) {
	//ch := new(chan struct{})
	ch := make(chan struct{})
	go func() {
		<-ch
		fmt.Println("111")
	}()
	ch <- struct{}{}
}
