package limiter

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	// 每秒生成1个令牌, 最多100个
	myLimiter := NewIPRateLimiter(1, 100)
	if !myLimiter.GetLimiter("127.0.0.1").Allow() {
		t.Error("too many requests")
	}
}

func Test2(t *testing.T) {
	// 每分钟最多请求20次: 每3秒生成1个令牌, 最多20次
	var myLimiter = NewIPRateLimiter(0.3, 20)

	n := 100
	for i := 0; i < n; i++ {
		if !myLimiter.GetLimiter("127.0.0.1").Allow() {
			t.Errorf("%d too many requests", i)
			time.Sleep(time.Second)
		}
	}
}

// mock了全量IP, 占用多少内存?
func Test3(t *testing.T) {
	setup()

	var myLimiter = NewIPRateLimiter(0.3, 20)

	saveMem("memprofile_1.prof")

	for i1 := 0; i1 < 255; i1++ {
		for i2 := 0; i2 < 255; i2++ {
			for i3 := 0; i3 < 255; i3++ {
				//for i4 := 0; i4 < 255; i4++ {
				ip := fmt.Sprintf("%d.%d.%d.%d", i1, i2, i3, 0)
				if !myLimiter.GetLimiter(ip).Allow() {
					t.Errorf("%s too many requests", ip)
				}
				//}
			}
		}
	}

	time.Sleep(time.Second * 5)
	saveMem("memprofile_2.prof")
	t.Logf("done: %v", len(myLimiter.ips))
}

func saveMem(filename string) {
	// 获取内存快照
	resp, err := http.Get("http://localhost:6060/debug/pprof/heap")
	if err != nil {
		fmt.Println("could not get heap profile:", err)
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("could not read response body:", err)
	}

	err = os.WriteFile(filename, bs, 0644)
	if err != nil {
		fmt.Println("could not read response body:", err)
	}
}

func setup() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
