package opentelemetry

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

// TestInitValidation 演示初始化时应该做的参数校验。
func TestInitValidation(t *testing.T) {
	ctx := context.Background()

	// service.name 是 resource 的核心字段，不能为空。
	if _, err := Init(ctx, Options{}, testLogger()); err == nil || !strings.Contains(err.Error(), "ServiceName") {
		t.Fatalf("空 ServiceName 应该报错，实际: %v", err)
	}

	// 采样比例必须在 [0,1]。
	for _, ratio := range []float64{-0.1, 1.1} {
		if _, err := Init(ctx, Options{ServiceName: "svc", SampleRatio: ratio}, testLogger()); err == nil {
			t.Errorf("SampleRatio=%v 应该报错", ratio)
		}
	}
}

// TestInitLogsAndResource 演示初始化日志与 resource 属性，以及 Shutdown 会把 span 冲刷出去。
func TestInitLogsAndResource(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 去掉时间，日志断言更稳定。
			if len(groups) == 0 && a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	exp := newMemExporter()
	p, err := Init(context.Background(), Options{
		ServiceName:    "checkout",
		ServiceVersion: "1.2.3",
		Environment:    "test",
		SampleRatio:    1,
		SpanExporter:   exp,
	}, logger)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	_, span := p.TracerProvider().Tracer("demo").Start(context.Background(), "hello")
	span.End()
	shutdown(t, p)

	got := logBuf.String()
	if !strings.Contains(got, "telemetry initialized") {
		t.Errorf("初始化日志缺失:\n%s", got)
	}
	if !strings.Contains(got, "service=checkout") {
		t.Errorf("初始化日志里应带 service 字段:\n%s", got)
	}

	spans := exp.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("Shutdown 后应导出 1 个 span，实际 %d 个", len(spans))
	}
	if got := resourceValue(t, spans[0].Resource, "service.name").AsString(); got != "checkout" {
		t.Errorf("resource service.name = %q, want checkout", got)
	}
	if got := resourceValue(t, spans[0].Resource, "deployment.environment").AsString(); got != "test" {
		t.Errorf("resource deployment.environment = %q, want test", got)
	}
}

// TestConsoleOutput 演示 ConsoleOutput：本地没有 collector 时，把 span 打到 stdout 观察。
func TestConsoleOutput(t *testing.T) {
	var out bytes.Buffer
	p, err := Init(context.Background(), Options{
		ServiceName:   "console-demo",
		SampleRatio:   1,
		ConsoleOutput: true,
		ConsoleWriter: &out,
	}, testLogger())
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	_, span := p.TracerProvider().Tracer("demo").Start(context.Background(), "console-span")
	span.End()
	shutdown(t, p) // Shutdown 会 flush BatchSpanProcessor，把 span 交给 stdout exporter

	if !strings.Contains(out.String(), "console-span") {
		t.Fatalf("console 输出里没有 span 名:\n%s", out.String())
	}
}

// TestNilProviderShutdown 演示 nil receiver 的容错：未初始化的场景调用 Shutdown 不应 panic。
func TestNilProviderShutdown(t *testing.T) {
	var p *Provider
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatalf("nil Provider Shutdown: %v", err)
	}
}
