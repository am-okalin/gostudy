package goruntine

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestProcessor(t *testing.T) {
	// 获取当前processor数量
	println(runtime.GOMAXPROCS(0))

	// 设置processor数量
	runtime.GOMAXPROCS(2)

	println(runtime.GOMAXPROCS(0))
}

func TestGoruntine(t *testing.T) {
	go func() {
		defer func() {
			t.Log("go_func_defer")
		}()

		t.Log("go_func")

		runtime.Goexit() //主动退出goruntine

		t.Log("go_func_exit")
	}()

	t.Log("func")

	time.Sleep(3 * time.Second)

	t.Log("func_done")
}

//Test1 猜猜输出什么
func Test1(t *testing.T) {
	num := 10
	for i := 0; i < num; i++ {
		go func() {
			fmt.Println(i)
		}()
	}

	//阻塞goruntine的方式有 select for-range time.sleep sync.wait chan
	//select {}
	time.Sleep(time.Millisecond * 3000)
}
