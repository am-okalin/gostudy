// WebSocket 双向通信示例：一个极简聊天室。
//
// 和 ../sse 对照着看：
//
//	sse —— 一条 HTTP 响应永不结束，服务端单向推，客户端只能另发普通请求
//	ws  —— 一次 HTTP Upgrade 之后就不再是 HTTP，双向、有帧、能 ping/pong
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
	addr := flag.String("addr", "127.0.0.1:9201", "监听地址")
	flag.Parse()

	hub := NewHub()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", page)       // 演示页面
	mux.HandleFunc("GET /ws", hub.ServeWS) // 升级入口

	srv := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		// 和 SSE 一样，不要设 WriteTimeout：WebSocket 也是长连接。
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		sctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := srv.Shutdown(sctx); err != nil {
			_ = srv.Close() // 还挂着 WS 连接时，Shutdown 会一直等
		}
	}()

	logger.Info("WebSocket 服务启动", "addr", srv.Addr, "page", "http://"+srv.Addr+"/")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("服务退出", "err", err)
		os.Exit(1)
	}
}
