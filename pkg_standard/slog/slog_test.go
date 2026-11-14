// Package slog 用一组示例测试演示标准库 log/slog（Go 1.21+）的用法。
//
// 一条日志由 时间(time)、级别(level)、消息(msg)、若干属性(Attr，即 key=value) 组成：
//
//   - Logger：负责组装日志，提供 Debug/Info/Warn/Error/Log/LogAttrs 等方法；
//
//   - Handler：决定日志输出成什么样（TextHandler / JSONHandler / 自定义），
//     同时决定级别过滤（Enabled）以及如何携带 With/WithGroup 的上下文。
//
//     1. MultiHandler 同时输出多种格式       TestMultiHandler
//     2. DiscardHandler 完全丢弃日志         TestDiscardHandler
//     3. 自定义 Handler 1：包装已有 Handler  TestContextHandler
//     4. 自定义 Handler 2：直接处理 Record   TestRecordHandler
//
// 另外，可直接运行的示例函数（Example_basic）见 example_test.go。
// 运行方式（pkg_standard 已在 go.work 中，可以直接在模块目录下跑）：
// cd pkg_standard && go test ./slog/ -v
// 参考文档：https://pkg.go.dev/log/slog
package slog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
)

// ----------------------------------------------------------------------------
// 公共辅助
// ----------------------------------------------------------------------------

// textOptions / jsonOptions 返回带固定时间的 Handler 配置；level 为 nil 时默认 INFO。
func textOptions(level slog.Leveler) *slog.HandlerOptions {
	return &slog.HandlerOptions{Level: level, ReplaceAttr: Replacer}
}

func jsonOptions(level slog.Leveler) *slog.HandlerOptions {
	return &slog.HandlerOptions{Level: level, ReplaceAttr: Replacer}
}

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("输出中缺少 %q\n实际输出:\n%s", want, got)
	}
}

// outputLines 把日志输出按行切开，去掉首尾空行。
func outputLines(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// ----------------------------------------------------------------------------
// 1. MultiHandler：一条日志同时交给多个 Handler（Go 1.26+）
// ----------------------------------------------------------------------------

func TestMultiHandler(t *testing.T) {
	var textBuf, jsonBuf bytes.Buffer
	logger := slog.New(slog.NewMultiHandler(
		slog.NewTextHandler(&textBuf, textOptions(nil)),
		slog.NewJSONHandler(&jsonBuf, jsonOptions(nil)),
	))
	logger.Info("multi", "k", "v")

	t.Logf("TextHandler 输出: %s", strings.TrimSpace(textBuf.String()))
	t.Logf("JSONHandler 输出: %s", strings.TrimSpace(jsonBuf.String()))
	assertContains(t, textBuf.String(), "level=INFO msg=multi k=v")
	assertContains(t, jsonBuf.String(), `"msg":"multi","k":"v"`)

	// MultiHandler.Enabled 的语义是“任一子 Handler 能处理即可”。
	ctx := context.Background()
	if logger.Enabled(ctx, slog.LevelDebug) {
		t.Error("两个 Handler 默认都是 INFO，DEBUG 不应该启用")
	}
	if !logger.Enabled(ctx, slog.LevelWarn) {
		t.Error("WARN 应该启用")
	}
}

// ----------------------------------------------------------------------------
// 2. DiscardHandler：完全丢弃日志（Go 1.27+），适合测试或彻底关闭日志
// ----------------------------------------------------------------------------

func TestDiscardHandler(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	ctx := context.Background()

	// DiscardHandler 对所有级别都返回 false，Logger 连 Record 都不会创建。
	if logger.Enabled(ctx, slog.LevelInfo) || logger.Enabled(ctx, slog.LevelError) {
		t.Fatal("DiscardHandler 不应启用任何级别")
	}
	logger.Info("should-not-appear") // 无输出、不 panic
}

// ----------------------------------------------------------------------------
// 3. 自定义 Handler 1：包装已有 Handler，从 context 追加属性（中间件模式）
// ----------------------------------------------------------------------------

type traceKey struct{}

// traceHandler 包装一个已有 Handler：在 Handle 里从 context 取出 trace_id 并加进记录。
// Handler 的四个方法都要转发，否则 With/WithGroup 会丢失。
type traceHandler struct {
	inner slog.Handler
}

func (h *traceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := ctx.Value(traceKey{}).(string); ok && id != "" {
		r.AddAttrs(slog.String("trace_id", id))
	}
	return h.inner.Handle(ctx, r)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{inner: h.inner.WithAttrs(attrs)}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{inner: h.inner.WithGroup(name)}
}

func TestContextHandler(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(&traceHandler{inner: slog.NewTextHandler(&buf, textOptions(nil))})

	// InfoContext / DebugContext 等会把 context 传给 Handler。
	ctx := context.WithValue(context.Background(), traceKey{}, "trace-123")
	logger.InfoContext(ctx, "with-ctx", "step", 1)
	// 没有 context 的 Info 拿不到 trace_id。
	logger.Info("without-ctx", "step", 2)

	lines := outputLines(buf.String())
	t.Logf("Context Handler 输出:\n%s", buf.String())
	assertContains(t, lines[0], "msg=with-ctx step=1 trace_id=trace-123")
	if strings.Contains(lines[1], "trace_id") {
		t.Fatalf("没有 context 时不应出现 trace_id: %s", lines[1])
	}
}

// ----------------------------------------------------------------------------
// 4. 自定义 Handler 2：直接处理 Record，观察 Handler 接口的组成细节
// ----------------------------------------------------------------------------

// recordHandler 把 Record 格式化成：LEVEL "msg" key=value ...
// 它演示了 Handler 需要自己处理：
//   - Enabled：级别过滤；
//   - Handle：读取 Record 的时间/级别/消息/属性；
//   - WithAttrs / WithGroup：保存并应用 Logger.With / WithGroup 带下来的上下文。
//
// 真实 Handler 必须并发安全，这里用互斥锁保护 Buffer。
type recordHandler struct {
	mu     *sync.Mutex
	w      io.Writer
	level  slog.Level
	prefix []string // 已经按当时的 groups 渲染好的 key=value 片段
	groups []string // 当前分组，作用于之后 Record 里的属性
}

func newRecordHandler(w io.Writer, level slog.Level) *recordHandler {
	return &recordHandler{mu: &sync.Mutex{}, w: w, level: level}
}

func (h *recordHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *recordHandler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %q", r.Level, r.Message)
	for _, p := range h.prefix {
		fmt.Fprintf(&b, " %s", p)
	}
	r.Attrs(func(a slog.Attr) bool {
		appendAttr(&b, h.groups, a)
		return true
	})

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := fmt.Fprintln(h.w, b.String())
	return err
}

func (h *recordHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// 先按“加入时的分组”把属性渲染好再保存：
	// 这样之后再调用 WithGroup 时，之前加入的属性不会被错误地放进新分组。
	prefix := append([]string{}, h.prefix...)
	for _, a := range attrs {
		if s := renderAttr(h.groups, a); s != "" {
			prefix = append(prefix, s)
		}
	}
	return &recordHandler{
		mu:     h.mu,
		w:      h.w,
		level:  h.level,
		prefix: prefix,
		groups: h.groups,
	}
}

func (h *recordHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &recordHandler{
		mu:     h.mu,
		w:      h.w,
		level:  h.level,
		prefix: h.prefix,
		groups: append(append([]string{}, h.groups...), name),
	}
}

// renderAttr 把单个 Attr 渲染成 "key=value"（分组递归展开成 group.key=value）。
func renderAttr(groups []string, a slog.Attr) string {
	var b strings.Builder
	appendAttr(&b, groups, a)
	return strings.TrimSpace(b.String())
}

// appendAttr 把 Attr 展开成 key=value；分组递归展开成 group.key=value。
// LogValuer 通过 Resolve 生效，零值 Attr 忽略——这些都是 Handler 的约定。
func appendAttr(b *strings.Builder, groups []string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}
	if a.Value.Kind() == slog.KindGroup {
		gs := groups
		if a.Key != "" {
			gs = append(append([]string{}, groups...), a.Key)
		}
		for _, ga := range a.Value.Group() {
			appendAttr(b, gs, ga)
		}
		return
	}
	key := a.Key
	if len(groups) > 0 {
		key = strings.Join(append(append([]string{}, groups...), a.Key), ".")
	}
	fmt.Fprintf(b, " %s=%v", key, a.Value.Any())
}

func TestRecordHandler(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(newRecordHandler(&buf, slog.LevelDebug)).With("app", "demo")
	logger.Info("saved", slog.Int("id", 7), slog.Group("db", slog.String("table", "orders")))
	logger.WithGroup("http").Info("request", slog.Int("status", 200))

	lines := outputLines(buf.String())
	t.Logf("自定义 Handler 输出:\n%s", buf.String())
	if len(lines) != 2 {
		t.Fatalf("期望 2 行，实际 %d 行", len(lines))
	}
	if want := `INFO "saved" app=demo id=7 db.table=orders`; lines[0] != want {
		t.Fatalf("第一行 = %q, want %q", lines[0], want)
	}
	// 注意：app=demo 是 WithGroup 之前 With 进去的，所以它不在 http 分组里，
	// 只有之后的属性（status）带 http. 前缀——这正是 Handler 契约要求的行为。
	if want := `INFO "request" app=demo http.status=200`; lines[1] != want {
		t.Fatalf("第二行 = %q, want %q", lines[1], want)
	}
}

// 另外：TopLevel 的 DebugContext/InfoContext/WarnContext/ErrorContext、
// Logger.Log、Logger.Handler、slog.With 等 API 也值得浏览一遍
// （go doc log/slog），本文件覆盖了日常使用中最常用的部分。
