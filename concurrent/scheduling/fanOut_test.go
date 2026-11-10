package taskscheduling

import (
	"fmt"
	"testing"
	"time"
)

func TestFanOut(t *testing.T) {
	var inch = make(chan any)
	var outchs []chan any
	for range 5 {
		outchs = append(outchs, make(chan any))
	}
	fanOut(true, inch, outchs)

	go func() {
		time.Sleep(2 * time.Second)
		inch <- 10
	}()

	val3 := <-outchs[3]
	val4 := <-outchs[4]
	fmt.Println(val3, val4)
}

func fanOut(sync bool, inch chan any, outchs []chan any) {
	go func() {
		defer func() {
			for i := range outchs {
				close(outchs[i])
			}
		}()

		for val := range inch {
			for i := range outchs {
				// go1.22 起循环变量每次迭代都会重新声明，无需再写 i := i 取新地址
				go func() {
					outchs[i] <- val
				}()
			}
		}
	}()
}
