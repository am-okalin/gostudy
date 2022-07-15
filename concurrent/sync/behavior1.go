package sync

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type T struct {
	a int
	b string
}

func (t T) isempty() bool {
	return t.a == 0 && t.b == ""
}

func Consume1(list []T, n int) error {
	// list的数据写入队列
	ch := make(chan T, len(list))
	for i := 0; i < len(list); i++ {
		ch <- list[i]
	}
	close(ch)

	// 多线程从chs中消费数据
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			val := <-ch
			for !val.isempty() {
				r, _ := rand.Int(rand.Reader, big.NewInt(5))
				t := time.Duration(r.Int64())
				time.Sleep(t * time.Millisecond)
				fmt.Println(i, val)
				val = <-ch
				if val.isempty() {
					println('a')
				}
			}
		}(i)
	}

	// 阻塞
	wg.Wait()
	return nil
}
