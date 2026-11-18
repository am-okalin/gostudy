// 裸 TCP 长连接示例：长度前缀分帧 + 心跳 + 单写者。
//
// 和 ws 版是同一个"聊天室"场景，但这里没有任何协议帮你：
// 消息边界、心跳、断线、名字握手，全都得自己定义。
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9301", "监听地址；客户端模式下是服务端地址")
	asClient := flag.Bool("client", false, "以客户端身份连接")
	name := flag.String("name", "", "客户端名字")
	flag.Parse()

	if *asClient {
		if err := runClient(*addr, *name); err != nil {
			log.Fatal(err)
		}
		return
	}

	srv, err := NewServer(*addr)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("TCP 服务启动: %s", srv.Addr())
	if err := srv.Serve(ctx); err != nil {
		log.Fatal(err)
	}
}
