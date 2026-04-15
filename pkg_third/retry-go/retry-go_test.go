package retry_go

import (
	"fmt"
	"github.com/avast/retry-go"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// Test 测试重试
func Test1(t *testing.T) {
	url := "http://exampl111e.com"
	var body []byte
	var count int

	err := retry.Do(
		func() error {
			resp, err := http.Get(url)
			if err != nil {
				if count == 2 { //假设第三次重试成功
					return nil
				}
				count++
				t.Errorf("create discount code error_%d: %v", count, err)
				return err
			}
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			return err
		},
		//重试次数
		retry.Attempts(3),
		//最大重试时间
		retry.MaxDelay(3*time.Second),
	)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(body)
}
