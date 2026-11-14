package gochan

import (
	"fmt"
	"testing"
	"time"
)

// Token 令牌
type Token struct{}

// TestDataTransferByToken 通过令牌进行数据传输
func TestDataTransferByToken(t *testing.T) {
	// 让令牌环跑一轮多，然后通知各个 worker 退出
	stop := make(chan struct{})
	time.AfterFunc(4500*time.Millisecond, func() { close(stop) })
	dataTransferByToken(stop)
}

func newWorker(id int, ch chan Token, nextCh chan Token, stop <-chan struct{}) {
	for {
		select {
		case token := <-ch: // 取得令牌
			fmt.Println(id)
			time.Sleep(time.Second)
			select {
			case nextCh <- token:
			case <-stop:
				return
			}
		case <-stop:
			return
		}
	}
}

func dataTransferByToken(stop <-chan struct{}) {
	chs := []chan Token{make(chan Token), make(chan Token), make(chan Token), make(chan Token)}

	// 创建4个worker
	for i := range 4 {
		go newWorker(i, chs[i], chs[(i+1)%4], stop)
	}

	//首先把令牌交给第一个worker
	chs[0] <- struct{}{}

	// 阻塞主协程，直到收到停止信号
	<-stop
}
