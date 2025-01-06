package cache

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

func TestLimiter(t *testing.T) {
	// 每3秒生成1个令牌, 最多10个, 1分钟过期
	l := NewTBLimiter(1.0/3, 10, time.Minute)
	for i := 0; i < 30; i++ {
		if !l.GetLimiter("127.0.0.1").Allow() {
			t.Errorf("%d too many requests", i+1)
		}
		time.Sleep(time.Second)
	}
}

func TestCache1(t *testing.T) {
	l := NewTBLimiter(1.0/3, 10, 10*time.Second)
	l.GetLimiter("127.0.0.1").Allow()
	i := 0
	for l.GetLimiter("127.0.0.2").Allow() {
		i++
		t.Log(i)
	}
	time.Sleep(10 * time.Second)
	t.Log(l.GetLimiter("127.0.0.2").Allow())
}

// TestCache 测试缓存过期
func TestCache(t *testing.T) {
	l := NewTBLimiter(1.0/3, 10, time.Minute)
	for i1 := 0; i1 < 255; i1++ {
		for i2 := 0; i2 < 255; i2++ {
			ip := fmt.Sprintf("%d.%d.%d.%d", i1, i2, 0, 0)
			if !l.GetLimiter(ip).Allow() {
				t.Errorf("%s too many requests", ip)
			}
		}
	}
}

// TestMemory mock全量IP, 占用多少内存?
func TestMemory(t *testing.T) {
	setup()

	myLimiter := NewTBLimiter(1.0/3, 20, time.Minute)

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
	t.Logf("done: %v", myLimiter.Len())
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
