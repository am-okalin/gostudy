package taskscheduling

import (
	"fmt"
	"testing"
	"time"
)

func TestFanOut(t *testing.T) {
	var inch = make(chan interface{})
	var outchs []chan interface{}
	for i := 0; i < 5; i++ {
		outchs = append(outchs, make(chan interface{}))
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

func fanOut(sync bool, inch chan interface{}, outchs []chan interface{}) {
	go func() {
		defer func() {
			for i := 0; i < len(outchs); i++ {
				close(outchs[i])
			}
		}()

		for val := range inch {
			for i := 0; i < len(outchs); i++ {
				i := i //给局部变量新的地址
				go func() {
					outchs[i] <- val
				}()
			}
		}
	}()
}
