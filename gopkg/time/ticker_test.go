package time

import (
	"fmt"
	"testing"
	"time"
)

//ticker 定时器，不断的按照设定的间隔时间触发，除非主动终止运行
func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second * 1)
	done := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case <-ticker.C:
				fmt.Println(i)
			}
		}
		ticker.Stop()
		done <- struct{}{}
	}()
	<-done
}

//timer 定时器，仅执行一次，可在执行后调用 timer.Reset() 让定时器再次工作，并可以更改时间间隔
func TestTimer(t *testing.T) {
	done := make(chan struct{})
	timer := time.NewTimer(time.Second * 2)
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case <-timer.C:
				fmt.Println(i, time.Now().Format("2006-01-02 15:04:05"))
			}
			timer.Reset(time.Second * 1)
		}
		done <- struct{}{}
	}()
	<-done

	//After() 方法其实是Timer的一个语法糖。
	fmt.Println("4秒后退出...")
	<-time.After(time.Second * 4)
}
