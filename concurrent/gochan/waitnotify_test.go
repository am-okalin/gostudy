package gochan

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

// TestWaitNotifyBySignal 使用chan实现信号通知
func TestWaitNotifyBySignal(t *testing.T) {
	// 测试里无法等待用户按 Ctrl+C，延迟向当前进程发送 SIGINT 触发优雅退出
	time.AfterFunc(200*time.Millisecond, func() {
		_ = syscall.Kill(os.Getpid(), syscall.SIGINT)
	})
	waitNotifyBySignal()
}

func waitNotifyBySignal() {
	var closing = make(chan struct{})
	var closed = make(chan struct{})

	go func() {
		// 模拟业务处理
		for {
			select {
			case <-closing:
				//处理一些退出业务
				return
			default:
				// ....... 业务计算
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// todo::标出执行test的命令
	// 处理CTRL+C等中断信号
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(termChan)
	//signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	close(closing)
	go doCleanup(closed)

	select {
	case <-closed:
	case <-time.After(time.Second):
		fmt.Println("清理超时，不等了")
	}
	fmt.Println("优雅退出")
}

func doCleanup(closed chan struct{}) {
	//time.Sleep(time.Minute)
	close(closed)
}
