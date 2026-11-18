package main

import (
	"net/http"
	"time"
)

// statusRecorder 包一层 ResponseWriter，用来记录状态码和响应字节数。
//
// 这正是 SSE 最容易被踩的地方：包完之后 w 不再是 *http.response，
// 也不再实现 http.Flusher，w.(http.Flusher) 会直接失败。
// Unwrap 让 http.NewResponseController 能穿透这层包装，找到底层 writer。
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (w *statusRecorder) WriteHeader(code int) {
	if w.status != 0 { // 和标准库一样，重复 WriteHeader 忽略
		return
	}
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusRecorder) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytes += int64(n)
	return n, err
}

func (w *statusRecorder) Unwrap() http.ResponseWriter { return w.ResponseWriter }

// withAccessLog 记录访问日志。SSE 的这条日志会等到连接断开才打出来。
func withAccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(rec, r)

		if rec.status == 0 {
			rec.status = http.StatusOK
		}
		logger.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"bytes", rec.bytes,
			"elapsed", time.Since(start).Round(time.Millisecond),
			"remote", r.RemoteAddr,
		)
	})
}
