package reflection

import (
	"fmt"
	"reflect"
	"time"
)

func initChs(num int) []chan interface{} {
	var chs = make([]chan interface{},0)
	for i := 0; i < num; i++ {
		chs = append(chs, make(chan interface{}, 10))
	}
	return chs
}

func Select1() {
	var num = 3
	var chs = initChs(num)
	var cases []reflect.SelectCase
	for i, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: reflect.ValueOf(fmt.Sprintf("chs[%v]", i)),
		})
	}
	for _, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}

	for i := 0; i < 3*num; i++ {
		//chosen 表示cases的索引；recvOk 表示是否有接受返回值；recv 表示接收到的返回值
		chosen, recv, ok := reflect.Select(cases)
		switch cases[chosen].Dir { //cases[chosen]可获取case对象
		case reflect.SelectSend:
			fmt.Println("send", cases[chosen].Send)
			fallthrough
		case reflect.SelectRecv:
			fmt.Println("recv", ok, recv)
			fallthrough
		case reflect.SelectDefault:
			time.Sleep(time.Second)
		}
	}
}

func Select2() {
	var num = 3
	var chs = initChs(num)
	var cases []reflect.SelectCase
	for _, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}

	close(chs[0])
	fmt.Println(reflect.Select(cases))
}
