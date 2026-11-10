package opentelemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// TestMetrics 演示四类常用指标：Counter、UpDownCounter、Histogram、ObservableGauge。
// 用 ManualReader 手动采集，所以不需要 collector，也不依赖后台定时上报。
func TestMetrics(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	p, _ := newTestProvider(t, Options{
		ServiceName:  "metrics-demo",
		SampleRatio:  1,
		MetricReader: reader,
	})

	ctx := context.Background()
	meter := p.MeterProvider().Meter("metrics-demo/instrumentation")

	// 1) Counter：只增不减，适合请求数、错误数。
	requests, err := meter.Int64Counter("http.requests", metric.WithDescription("HTTP 请求总数"))
	if err != nil {
		t.Fatalf("创建 counter: %v", err)
	}
	requests.Add(ctx, 3, metric.WithAttributes(attribute.String("method", "GET")))
	requests.Add(ctx, 2, metric.WithAttributes(attribute.String("method", "POST")))

	// 2) Histogram：记录分布，适合耗时、响应体大小；关心的是分位数而不是单值。
	duration, err := meter.Float64Histogram("request.duration", metric.WithUnit("ms"))
	if err != nil {
		t.Fatalf("创建 histogram: %v", err)
	}
	duration.Record(ctx, 12.5)
	duration.Record(ctx, 30)

	// 3) UpDownCounter：可增可减，适合队列长度、正在处理的请求数。
	queue, err := meter.Int64UpDownCounter("queue.depth")
	if err != nil {
		t.Fatalf("创建 up-down counter: %v", err)
	}
	queue.Add(ctx, 7)

	// 4) ObservableGauge：回调式，采集时读取当前值，适合缓存条目数、连接池大小。
	cacheEntries, err := meter.Int64ObservableGauge("cache.entries")
	if err != nil {
		t.Fatalf("创建 observable gauge: %v", err)
	}
	if _, err := meter.RegisterCallback(func(_ context.Context, observer metric.Observer) error {
		observer.ObserveInt64(cacheEntries, 42)
		return nil
	}, cacheEntries); err != nil {
		t.Fatalf("注册回调: %v", err)
	}

	// 手动触发一次采集（生产环境由 PeriodicReader 定时做这件事）。
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		t.Fatalf("Collect: %v", err)
	}
	shutdown(t, p)

	// resource 同样会跟着指标一起上报。
	if got := resourceValue(t, rm.Resource, "service.name").AsString(); got != "metrics-demo" {
		t.Errorf("resource service.name = %q", got)
	}

	// Counter：两个 method 各一个数据点，合计 5。
	sum, ok := findMetric(t, rm, "http.requests").Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("http.requests 的数据类型不是 Sum[int64]")
	}
	var total int64
	for _, dp := range sum.DataPoints {
		total += dp.Value
	}
	if total != 5 {
		t.Errorf("http.requests 合计 = %d, want 5", total)
	}
	if !sum.IsMonotonic {
		t.Errorf("Counter 的 Sum 应该是 monotonic")
	}

	// Histogram：记录 2 次，总和 42.5。
	hist, ok := findMetric(t, rm, "request.duration").Data.(metricdata.Histogram[float64])
	if !ok {
		t.Fatalf("request.duration 的数据类型不是 Histogram[float64]")
	}
	if len(hist.DataPoints) != 1 || hist.DataPoints[0].Count != 2 || hist.DataPoints[0].Sum != 42.5 {
		t.Errorf("request.duration 数据点异常: %+v", hist.DataPoints)
	}

	// UpDownCounter：当前值 7。
	updown, ok := findMetric(t, rm, "queue.depth").Data.(metricdata.Sum[int64])
	if !ok || len(updown.DataPoints) != 1 || updown.DataPoints[0].Value != 7 {
		t.Errorf("queue.depth 数据点异常: %+v", updown.DataPoints)
	}

	// ObservableGauge：回调里返回 42。
	gauge, ok := findMetric(t, rm, "cache.entries").Data.(metricdata.Gauge[int64])
	if !ok || len(gauge.DataPoints) != 1 || gauge.DataPoints[0].Value != 42 {
		t.Errorf("cache.entries 数据点异常: %+v", gauge.DataPoints)
	}
}

func findMetric(t *testing.T, rm metricdata.ResourceMetrics, name string) metricdata.Metrics {
	t.Helper()
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == name {
				return m
			}
		}
	}
	t.Fatalf("没有找到指标 %q", name)
	return metricdata.Metrics{}
}
