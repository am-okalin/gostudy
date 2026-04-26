package leetcode

import (
	"fmt"
	"testing"
	"time"
)

func worker(key string, now, next chan int) {
	for {
		// 获取令牌
		t := <-now

		// 出口
		if t == 101 || t == 0 {
			break
		}

		// 打印
		time.Sleep(100 * time.Microsecond)
		fmt.Printf("%s: %d\n", key, t)

		// 令牌给下个ch
		next <- t + 1
	}
}

// 交叉输出1-100
func TestGo2(t *testing.T) {
	// 令牌
	tc1 := make(chan int)
	tc2 := make(chan int)

	go worker("go1", tc1, tc2)
	go worker("go2", tc2, tc1)

	// 令牌给第一个worker
	tc1 <- 1

	// 阻塞主协程
	time.Sleep(3 * time.Second)
	close(tc1)
	close(tc2)
}
