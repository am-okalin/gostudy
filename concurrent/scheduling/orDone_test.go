package taskscheduling

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// 返回c并在执行函数后关闭c
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func TestOrr(t *testing.T) {
	start := time.Now()

	// 在没有函数完成前一直阻塞
	<-orr(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
		sig(4*time.Second),
		sig(5*time.Second),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}

func TestOrs(t *testing.T) {
	start := time.Now()

	<-ors(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
	)

	fmt.Printf("done after %v", time.Since(start))
}

//orr 递归方式的orDone
func orr(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有零个或者1个chan
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2: // 2个也是一种特殊情况
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: //超过两个，二分法递归处理
			m := len(channels) / 2
			select {
			case <-orr(channels[:m]...):
			case <-orr(channels[m:]...):
			}
		}
	}()

	return orDone
}

//ors select方式的orDone
func ors(chans ...<-chan interface{}) <-chan interface{} {
	if len(chans) == 0 {
		return nil
	}

	orDone := make(chan interface{})

	go func() {
		defer close(orDone)
		var cases []reflect.SelectCase
		for _, ch := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			})
		}
		reflect.Select(cases)
	}()

	return orDone
}
