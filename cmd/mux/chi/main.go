// 与 ../std 逻辑一致，只换 router。中间件全部来自 chi 的 middleware 包。
package main

import (
	"encoding/json"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

func main() {
	r := chi.NewRouter()

	// r.Use 洋葱从外到内, 先注册的在最外层.
	r.Use(middleware.RequestID)                                     // 生成/透传 request id，放进 ctx
	r.Use(middleware.RequestLogger(&slogFormatter{logger: logger})) // slogFormatter 把 chi 的日志 entry 接到 slog 上。
	r.Use(middleware.Recoverer)                                     // panic → 栈 + 500
	r.Use(middleware.Timeout(5 * time.Second))                      // ctx 5s 预算，返回后补 504

	// /panic 必须注册在 /* 之前，否则会被通配路由先吃掉。
	r.Post("/", createUser)
	r.Get("/", echoPath)
	r.Get("/panic", panicHandler) // 只为演示 Recoverer
	r.Get("/*", echoPath)

	srv := &http.Server{
		Addr:              "127.0.0.1:9102",
		Handler:           r,
		ReadHeaderTimeout: 3 * time.Second,  // 读请求头的上限，防 slowloris 慢慢发 header
		ReadTimeout:       10 * time.Second, // 读完整个请求（头 + 体）的上限
		WriteTimeout:      10 * time.Second, // 从读完请求头到响应写完的上限，必须大于 5s
		IdleTimeout:       60 * time.Second, // keep-alive 空闲连接的保留时长
	}

	logger.Info("chi.Mux", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("服务退出", "err", err)
		os.Exit(1)
	}
}

// 实现 NewLogEntry(r *http.Request) LogEntry
type slogFormatter struct{ logger *slog.Logger }

func (f *slogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &slogEntry{logger: f.logger, r: r}
}

type slogEntry struct {
	logger *slog.Logger
	r      *http.Request
}

func (e *slogEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ any) {
	ctx := e.r.Context()
	e.logger.LogAttrs(ctx, slog.LevelInfo, "http request",
		slog.String("request_id", middleware.GetReqID(ctx)),
		slog.String("method", e.r.Method),
		slog.String("path", e.r.URL.Path),
		slog.Int("status", status),
		slog.Int("bytes", bytes),
		slog.Duration("elapsed", elapsed),
	)
}

func (e *slogEntry) Panic(v any, stack []byte) {
	e.logger.ErrorContext(e.r.Context(), "panic recovered",
		slog.Any("panic", v),
		slog.String("request_id", middleware.GetReqID(e.r.Context())),
		slog.String("method", e.r.Method),
		slog.String("path", e.r.URL.Path),
		slog.String("stack", string(stack)),
	)
}

/* - handler 返回 - */

type user struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var req user
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体不是合法 JSON")
		return
	}
	req.ID = strconv.FormatUint(rand.Uint64(), 10)
	writeJSON(w, http.StatusCreated, req)
}

func echoPath(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"path": r.URL.Path,
		"rest": strings.TrimPrefix(r.URL.Path, "/"),
	})
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("故意 panic，演示 Recoverer")
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Error("写响应失败", "err", err)
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
