// 标准库 ServeMux 版。中间件没有现成的，全部自己实现，顺序和 chi 那版一致。
package main

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"sync/atomic"
	"time"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

type Middleware = func(http.Handler) http.Handler

type ctxKey int

const ctxKeyRequestID ctxKey = 0

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{$}", createUser)    // {$} 精确匹配 "/"
	mux.HandleFunc("GET /panic", panicHandler) // 比 "GET /" 更具体，优先命中
	mux.HandleFunc("GET /", echoPath)          // "/" 是子树匹配，兜住其余 GET

	// 顺序：requestID → accessLog → recoverer → timeout
	handler := chain(mux,
		requestID(),            // ctx 里的 id
		accessLog(logger),      // 访问日志
		recoverer(logger),      // panic → 栈 + 500
		timeout(5*time.Second), // ctx 5s 预算，返回后补 504
	)

	srv := &http.Server{
		Addr:              "127.0.0.1:9101",
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,  // 读请求头的上限，防 slowloris 慢慢发 header
		ReadTimeout:       10 * time.Second, // 读完整个请求（头 + 体）的上限
		WriteTimeout:      10 * time.Second, // 从读完请求头到响应写完的上限，必须大于 5s
		IdleTimeout:       60 * time.Second, // keep-alive 空闲连接的保留时长
	}

	logger.Info("stdlib ServeMux", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("服务退出", "err", err)
		os.Exit(1)
	}
}

// chain 按参数顺序从外到内套。
func chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

var (
	reqIDPrefix  = newPrefix()
	reqIDCounter atomic.Uint64
)

func newPrefix() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		host = "localhost"
	}
	var buf [12]byte
	crand.Read(buf[:])
	return host + "/" + base64.RawURLEncoding.EncodeToString(buf[:])[:10]
}

// requestID 优先用上游传来的 X-Request-Id，没有就生成 host/随机串-000001。
func requestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-Id")
			if id == "" {
				id = fmt.Sprintf("%s-%06d", reqIDPrefix, reqIDCounter.Add(1))
			}
			ctx := context.WithValue(r.Context(), ctxKeyRequestID, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// accessLog 记一条访问日志。
func accessLog(l *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				l.LogAttrs(r.Context(), slog.LevelInfo, "http request",
					slog.String("request_id", reqIDFrom(r.Context())),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Duration("elapsed", time.Since(start)),
				)
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// recoverer 兜住 panic，栈走 slog，再写 500。
func recoverer(l *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rvr := recover()
				if rvr == nil {
					return
				}
				if rvr == http.ErrAbortHandler {
					panic(rvr) // 让连接直接断掉，不写响应
				}
				l.ErrorContext(r.Context(), "panic recovered",
					slog.Any("panic", rvr),
					slog.String("request_id", reqIDFrom(r.Context())),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("stack", string(debug.Stack())),
				)
				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// timeout 给 ctx 加预算。和 chi 一样是协作式的：handler 不看 ctx 就照跑，504 是返回后补写的。
func timeout(d time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer func() {
				cancel()
				if ctx.Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
			}()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func reqIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(ctxKeyRequestID).(string)
	return id
}
