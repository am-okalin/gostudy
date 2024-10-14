package goruntine

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestProcessorNum(t *testing.T) {
	// 获取当前processor数量
	println(runtime.GOMAXPROCS(-1))
	println(runtime.GOMAXPROCS(0))

	// 设置processor数量
	runtime.GOMAXPROCS(2)

	println(runtime.GOMAXPROCS(0))
}

func Test1(t *testing.T) {
	go func() {
		defer func() {
			t.Log("go_func_defer")
		}()

		t.Log("go_func")

		//runtime.Goexit() //主动退出goruntine
		panic("downing test") //panic会导致上层goruntine也panic

		t.Log("go_func_exit")
	}()

	t.Log("func")

	time.Sleep(3 * time.Second)

	t.Log("func_done")
}

func Test2(t *testing.T) {
	runtime.GOMAXPROCS(1)

	num := 9
	for i := 0; i < num; i++ {
		go func() {
			fmt.Println("go1: ", i) //block
		}()
	}

	//for i := 0; i < num; i++ {
	//	go func(i int) {
	//		fmt.Println("go2: ", i)
	//	}(i)
	//}

	time.Sleep(time.Millisecond * 3000) //block
}
