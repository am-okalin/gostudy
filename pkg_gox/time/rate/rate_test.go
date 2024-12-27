package rate

import (
	"fmt"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

// 速率限制其
func TestLimiter(t *testing.T) {
	// 创建一个速率限制器，每3秒创建1个令牌, 最多3个令牌
	limiter := rate.NewLimiter(1.0/3.0, 3)

	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			t.Log("success", i)
		} else {
			t.Log("failed", i)
		}
		time.Sleep(time.Second) // 控制循环速度
	}

	// 等待10秒，恢复至最多3个令牌
	time.Sleep(10 * time.Second)

	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			t.Log("success", i)
		} else {
			t.Log("failed", i)
		}
		time.Sleep(time.Second) // 控制循环速度
	}
}

// 每N次/每个时间间隔 执行操作
func TestSometime(t *testing.T) {
	s := &rate.Sometimes{
		First:    5,               // 第5次执行
		Every:    2,               // 每2次执行一次
		Interval: 1 * time.Second, // 时间间隔为1秒, 每隔1秒执行一次
	}

	for i := 0; i < 100; i++ {
		s.Do(func() {
			fmt.Println("Executing operation at", time.Now())
		})
		time.Sleep(300 * time.Millisecond) // 控制循环速度
	}
}
