package gochan

import (
	"fmt"
	"testing"
	"time"
)

//Token 令牌
type Token struct{}

//TestDataTransferByToken 通过令牌进行数据传输
func TestDataTransferByToken(t *testing.T) {
	dataTransferByToken()
}

func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch // 取得令牌
		fmt.Println(id)
		time.Sleep(time.Second)
		nextCh <- token
	}
}

func dataTransferByToken() {
	chs := []chan Token{make(chan Token), make(chan Token), make(chan Token), make(chan Token)}

	// 创建4个worker
	for i := 0; i < 4; i++ {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}

	//首先把令牌交给第一个worker
	chs[0] <- struct{}{}

	// 阻塞主进程
	select {}
}
