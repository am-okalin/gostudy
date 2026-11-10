package opentelemetry

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// testLogger 让示例测试保持安静，不把 OTel 的初始化日志打到终端。
func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// memExporter 是一个 Shutdown 时不会清空数据的 span 导出器。
//
// 为什么不用 tracetest.NewInMemoryExporter()？因为它的 Shutdown 会 Reset：
// TracerProvider.Shutdown 会 flush BatchSpanProcessor 并调用每个 exporter 的 Shutdown，
// 于是 "span.End() → p.Shutdown() → exp.GetSpans()" 会拿到空列表（这是个常见坑）。
// 想保留数据做断言，就换成这个导出器（或者改用 SimpleSpanProcessor 在 Shutdown 前读）。
type memExporter struct {
	mu    sync.Mutex
	spans []sdktrace.ReadOnlySpan
}

func newMemExporter() *memExporter { return &memExporter{} }

func (e *memExporter) ExportSpans(_ context.Context, spans []sdktrace.ReadOnlySpan) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.spans = append(e.spans, spans...)
	return nil
}

func (e *memExporter) Shutdown(context.Context) error { return nil }

func (e *memExporter) GetSpans() tracetest.SpanStubs {
	e.mu.Lock()
	defer e.mu.Unlock()
	return tracetest.SpanStubsFromReadOnlySpans(e.spans)
}

// newTestProvider 用内存导出器创建 Provider，避免测试依赖真实 collector。
// 注意 SampleRatio 不传默认值：0 是合法值（不采样），需要各测试自己显式指定。
func newTestProvider(t *testing.T, opts Options) (*Provider, *memExporter) {
	t.Helper()
	exp := newMemExporter()
	opts.SpanExporter = exp
	if opts.ServiceName == "" {
		opts.ServiceName = "otel-sample"
	}
	p, err := Init(context.Background(), opts, testLogger())
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	return p, exp
}

// shutdown 用带超时的 ctx 关闭 Provider，模拟真实服务退出流程。
func shutdown(t *testing.T, p *Provider) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func findSpan(t *testing.T, spans tracetest.SpanStubs, name string) tracetest.SpanStub {
	t.Helper()
	for _, s := range spans {
		if s.Name == name {
			return s
		}
	}
	names := make([]string, 0, len(spans))
	for _, s := range spans {
		names = append(names, s.Name)
	}
	t.Fatalf("没有找到 span %q，实际有: %v", name, names)
	return tracetest.SpanStub{}
}

func attrValue(attrs []attribute.KeyValue, key string) (attribute.Value, bool) {
	for _, kv := range attrs {
		if string(kv.Key) == key {
			return kv.Value, true
		}
	}
	return attribute.Value{}, false
}

func resourceValue(t *testing.T, res *resource.Resource, key string) attribute.Value {
	t.Helper()
	if res == nil {
		t.Fatalf("resource 为 nil")
	}
	v, ok := res.Set().Value(attribute.Key(key))
	if !ok {
		t.Fatalf("resource 缺少属性 %q", key)
	}
	return v
}
