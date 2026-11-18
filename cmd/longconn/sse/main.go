// SSE（Server-Sent Events）最小可运行示例。
//
// 四个关键点：
//  1. 响应头 text/event-stream + no-cache，nginx 还要加 X-Accel-Buffering: no；
//  2. 每写一条就 Flush，否则数据会攒在 bufio 里，客户端看着像"卡住"；
//  3. 用 http.NewResponseController 而不是 w.(http.Flusher)，中间件包了一层也不怕；
//  4. 不要给 Server 设 WriteTimeout —— 它从读完请求头开始计时，到点就把长连接掐断。
package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

func main() {
	addr := flag.String("addr", "127.0.0.1:9103", "监听地址")
	tick := flag.Duration("tick", 5*time.Second, "服务端定时推送间隔，0 表示关闭")
	flag.Parse()

	b := NewBroker()
	if *tick > 0 {
		go b.tick(*tick)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", index)              // 演示页面
	mux.HandleFunc("GET /events", b.Stream)        // SSE 流
	mux.HandleFunc("POST /publish", b.PublishHTTP) // 手动广播一条

	srv := &http.Server{
		Addr:              *addr,
		Handler:           withAccessLog(mux),
		ReadHeaderTimeout: 5 * time.Second,
		// 故意不设 WriteTimeout / IdleTimeout：SSE 是长连接，
		// 设了 WriteTimeout 就会在 N 秒后把流切断（这是最常见的"线上 SSE 十分钟必断"）。
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		sctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := srv.Shutdown(sctx); err != nil {
			_ = srv.Close() // 还挂着 SSE 长连接时，Shutdown 会一直等，超时后强关
		}
	}()

	logger.Info("SSE 服务启动", "addr", srv.Addr, "page", "http://"+srv.Addr+"/")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("服务退出", "err", err)
		os.Exit(1)
	}
}
