package cycleimport

import (
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	go func() {
		fmt.Println(1111)
	}()
	time.Sleep(3 * time.Second)
}
