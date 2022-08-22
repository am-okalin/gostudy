package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	go func() {
		sch := make(chan os.Signal)
		signal.Notify(sch, syscall.SIGINT, syscall.SIGTERM)
		v := <-sch
		fmt.Println("------", v)
	}()
	time.Sleep(time.Second * 5)

}