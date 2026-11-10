package opentelemetry

import (
	"context"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TestSpanLifecycle 演示一条最普通的 trace：
// 创建 span → 加属性/事件 → 创建子 span → 设置状态 → 结束后被 exporter 收到。
func TestSpanLifecycle(t *testing.T) {
	p, exp := newTestProvider(t, Options{ServiceName: "trace-demo", ServiceVersion: "1.0.0", SampleRatio: 1})
	tracer := p.TracerProvider().Tracer("trace-demo/instrumentation")

	ctx, root := tracer.Start(context.Background(), "GET /orders",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("http.method", "GET"),
			attribute.String("http.route", "/orders"),
		),
	)
	// 事件适合记录"某个时间点发生了什么"（日志式），属性适合可查询的维度。
	root.AddEvent("cache.miss", trace.WithAttributes(attribute.String("key", "orders:42")))
	root.SetAttributes(attribute.Int("http.status_code", 200))

	// 子 span 会自动成为 root 的 child（Context 里带着当前 span）。
	_, child := tracer.Start(ctx, "SELECT orders",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.String("db.system", "mysql")),
	)
	child.End()

	root.SetStatus(codes.Ok, "")
	root.End()

	shutdown(t, p)

	spans := exp.GetSpans()
	if len(spans) != 2 {
		t.Fatalf("期望导出 2 个 span，实际 %d 个", len(spans))
	}

	rootStub := findSpan(t, spans, "GET /orders")
	childStub := findSpan(t, spans, "SELECT orders")

	if rootStub.SpanKind != trace.SpanKindServer {
		t.Errorf("root span kind = %v, want server", rootStub.SpanKind)
	}
	if len(rootStub.Events) != 1 || rootStub.Events[0].Name != "cache.miss" {
		t.Errorf("root span 事件不对: %+v", rootStub.Events)
	}
	if v, ok := attrValue(rootStub.Attributes, "http.status_code"); !ok || v.AsInt64() != 200 {
		t.Errorf("root span http.status_code = %v", v)
	}
	if rootStub.Status.Code != codes.Ok {
		t.Errorf("root span status = %v, want Ok", rootStub.Status.Code)
	}

	// 子 span 通过 Context 继承了 parent 的 SpanContext。
	if childStub.Parent.TraceID() != rootStub.SpanContext.TraceID() {
		t.Errorf("子 span 的 TraceID 与父 span 不一致")
	}
	if childStub.Parent.SpanID() != rootStub.SpanContext.SpanID() {
		t.Errorf("子 span 的 ParentSpanID 应该指向 root span")
	}
	if v, ok := attrValue(childStub.Attributes, "db.system"); !ok || v.AsString() != "mysql" {
		t.Errorf("child span db.system = %v", v)
	}

	// span 上带着初始化时配置的 resource 属性。
	if got := resourceValue(t, rootStub.Resource, "service.name").AsString(); got != "trace-demo" {
		t.Errorf("resource service.name = %q", got)
	}
	if got := resourceValue(t, rootStub.Resource, "service.version").AsString(); got != "1.0.0" {
		t.Errorf("resource service.version = %q", got)
	}
}

// TestSampling 演示 ParentBased + TraceIDRatioBased 的效果：
// 比例 0 不采样（span 根本不会创建），比例 1 全采样。
func TestSampling(t *testing.T) {
	cases := []struct {
		ratio float64
		want  int
	}{
		{ratio: 0, want: 0},
		{ratio: 1, want: 1},
	}

	for _, tc := range cases {
		p, exp := newTestProvider(t, Options{ServiceName: "sampling-demo", SampleRatio: tc.ratio})
		_, span := p.TracerProvider().Tracer("t").Start(context.Background(), "op")
		span.End()
		shutdown(t, p)

		if got := len(exp.GetSpans()); got != tc.want {
			t.Errorf("SampleRatio=%v 时导出 %d 个 span，期望 %d 个", tc.ratio, got, tc.want)
		}
	}
}

// TestContextPropagation 演示跨服务传播：Inject/Extract W3C tracecontext + baggage，
// 对端服务从 HTTP header 里恢复出同一个 trace，并继续往下游发请求。
func TestContextPropagation(t *testing.T) {
	p, exp := newTestProvider(t, Options{ServiceName: "gateway", SampleRatio: 1})
	tracer := p.TracerProvider().Tracer("gateway")

	// 服务 A：处理请求，并把 trace 上下文 + baggage 写进 HTTP header。
	ctx, server := tracer.Start(context.Background(), "server",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	member, err := baggage.NewMember("tenant", "acme")
	if err != nil {
		t.Fatalf("baggage.NewMember: %v", err)
	}
	bag, err := baggage.New(member)
	if err != nil {
		t.Fatalf("baggage.New: %v", err)
	}
	ctx = baggage.ContextWithBaggage(ctx, bag)

	carrier := propagationMapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	if got := carrier.Get("traceparent"); !strings.HasPrefix(got, "00-") {
		t.Errorf("traceparent 格式不对: %q", got)
	}
	if got := carrier.Get("baggage"); !strings.Contains(got, "tenant=acme") {
		t.Errorf("baggage 没有传播: %q", got)
	}

	// 服务 B：从 header 恢复上下文，再创建自己的 span（Parent 会指向服务 A 的 span）。
	ctx2 := otel.GetTextMapPropagator().Extract(context.Background(), carrier)
	remote := trace.SpanContextFromContext(ctx2)
	if remote.TraceID() != server.SpanContext().TraceID() {
		t.Errorf("Extract 后 TraceID 不一致: %v vs %v", remote.TraceID(), server.SpanContext().TraceID())
	}
	if got := baggage.FromContext(ctx2).Member("tenant").Value(); got != "acme" {
		t.Errorf("Extract 后 baggage tenant = %q", got)
	}

	_, client := tracer.Start(ctx2, "call-downstream",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	client.End()
	server.End()

	shutdown(t, p)

	spans := exp.GetSpans()
	clientStub := findSpan(t, spans, "call-downstream")
	serverStub := findSpan(t, spans, "server")
	if clientStub.Parent.TraceID() != serverStub.SpanContext.TraceID() {
		t.Errorf("下游 span 的 TraceID 与服务端不一致")
	}
	if clientStub.Parent.SpanID() != serverStub.SpanContext.SpanID() {
		t.Errorf("下游 span 的 Parent 应该指向服务端 span")
	}
}

// propagationMapCarrier 演示 TextMapCarrier 的最小实现；实际项目里可以用 propagation.MapCarrier，
// 接入 HTTP 时通常用 otelhttp / otelgin 等库自动完成 Inject/Extract。
type propagationMapCarrier map[string]string

func (c propagationMapCarrier) Get(key string) string { return c[key] }

func (c propagationMapCarrier) Set(key, value string) { c[key] = value }

func (c propagationMapCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
