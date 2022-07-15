package sync

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"time"
)

//Producer, consumer
func InCh(chs []chan string, list []string) {
	// 数据输入完毕后关闭chs
	defer func() {
		for _, outch := range chs {
			close(outch)
		}
	}()

	// 循环喂数据(block)
	chslen := len(chs)
	inCases := make([]reflect.SelectCase, chslen)
	for i := 0; i < len(list); i++ {
		for j := 0; j < len(chs); j++ {
			//实现 case: chs[i%n]<-v
			inCases[j%chslen] = reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(chs[j]),
				Send: reflect.ValueOf(list[i]),
			}
		}
		reflect.Select(inCases)
	}
}

func OutCh(chs []chan string, i int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	val := <-chs[i]
	for val != "" {
		r, _ := rand.Int(rand.Reader, big.NewInt(5))
		t := time.Duration(r.Int64())
		time.Sleep(t * time.Second)
		fmt.Println(i, val)
		val = <-chs[i]
	}
}

func Consume(list []string, n int) error {
	// 创建无缓存的chan
	chs := make([]chan string, n)
	for i := 0; i < n; i++ {
		chs[i] = make(chan string)
	}

	// 往多个chs(消费队列)中喂数据
	go InCh(chs, list)

	// 多线程从chs中消费数据
	wg := sync.WaitGroup{}
	wg.Add(len(chs))
	for i := 0; i < n; i++ {
		go OutCh(chs, i, &wg)
	}

	// 阻塞
	wg.Wait()
	return nil
}
