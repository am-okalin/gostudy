package taskscheduling

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func fanInReflect(inchs ...<-chan interface{}) <-chan interface{} {
	var cases []reflect.SelectCase
	for _, inch := range inchs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(inch),
		})
	}

	outch := make(chan interface{})
	go func() {
		defer close(outch)
		for len(cases) > 0 {
			chose, recv, ok := reflect.Select(cases)
			if !ok && cases[chose].Dir == reflect.SelectRecv { //channel已经close
				cases = append(cases[:chose], cases[chose+1:]...)
				continue
			}
			outch <- recv.Interface()
		}
	}()
	return outch
}

func Input(val int, d time.Duration) <-chan interface{} {
	var inchan = make(chan interface{})
	go func() {
		defer close(inchan)
		time.Sleep(d)
		inchan <- val
	}()
	return inchan
}

func TestFanInReflect(t *testing.T) {
	outch := fanInReflect(
		Input(00, 0*time.Second),
		Input(01, 1*time.Second),
		Input(02, 2*time.Second),
		Input(03, 3*time.Second),
		Input(04, 4*time.Second),
	)
	fmt.Println("已获取输出管道")
	for {
		val := <-outch
		if val == nil {
			time.Sleep(time.Second)
			fmt.Println("输出结束")
			return
		}
		fmt.Println("output:", val)
	}

}
