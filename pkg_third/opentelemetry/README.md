# OpenTelemetry 样例（pkg_third/opentelemetry）

本目录是真实项目 `internal/pkg/telemetry` 的简化版示例，用来学习 Go 里如何接入
[OpenTelemetry](https://opentelemetry.io/) 的 **trace（链路）** 与 **metric（指标）**。

依赖版本与真实项目一致：`go.opentelemetry.io/otel v1.38.0`。

## 文件说明

| 文件 | 内容 |
|---|---|
| `telemetry.go` | `Init` / `Provider.Shutdown`：resource、TracerProvider（Sampler + BatchSpanProcessor）、MeterProvider（PeriodicReader）、全局注册、W3C 传播、错误日志 |
| `trace_test.go` | span 生命周期（属性/事件/状态/子 span）、采样比例、跨服务 Inject/Extract 传播 |
| `metric_test.go` | Counter / UpDownCounter / Histogram / ObservableGauge，以及 ManualReader 手动采集 |
| `init_test.go` | 参数校验、resource 属性、ConsoleOutput、Shutdown 行为 |
| `helpers_test.go` | 测试辅助：内存导出器、静默 logger、断言工具 |

## 运行

`pkg_third` 模块没有加入根目录的 `go.work`，所以要关闭 workspace 模式：

```bash
cd pkg_third && GOWORK=off go test ./opentelemetry/ -v
```

只想跑某一个主题：

```bash
cd pkg_third && GOWORK=off go test ./opentelemetry/ -run 'TestSpanLifecycle|TestMetrics' -v
```

## 最小接入示例（生产代码）

```go
ctx := context.Background()
logger := slog.Default()

p, err := opentelemetry.Init(ctx, opentelemetry.Options{
    ServiceName:    "order-service",
    ServiceVersion: "1.2.3",
    Environment:    "prod",
    Endpoint:       "otel-collector:4317", // OTLP gRPC
    Insecure:       false,                 // 内网明文则设 true
    SampleRatio:    0.1,                   // 10% 采样；ParentBased 会跟随上游决定
    MetricInterval: 30 * time.Second,
}, logger)
if err != nil {
    logger.Error("init telemetry", "error", err)
    return err
}

// 退出前必须 Shutdown，并给 ctx 加超时：
// collector 不可达时导出会阻塞，用 context.Background() 可能导致进程退不掉。
defer func() {
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _ = p.Shutdown(shutdownCtx)
}()

tracer := p.TracerProvider().Tracer("order-service")
ctx, span := tracer.Start(ctx, "CreateOrder")
defer span.End()
```

本地没有 collector 时，可以打开 `ConsoleOutput: true` 把 span 打到 stdout（见 `TestConsoleOutput`）。

## 本地查看链路（可选）

```bash
# 启动 Jaeger（OTLP gRPC 4317）
docker run --rm -p 16686:16686 -p 4317:4317 jaegertracing/all-in-one:latest

# 然后在代码里设置 Endpoint: "localhost:4317"、Insecure: true，跑起来后打开
# http://localhost:16686 搜索 service.name
```

## 六个常见坑（代码里都有对应注释）

1. **忘记 Shutdown**：BatchSpanProcessor 会丢最后一批 span；`Shutdown` 必须带超时 ctx。
2. **不设全局 Provider**：`otelgrpc` / `otelhttp` 等库默认从 `otel.GetTracerProvider()` 取，不注册全局会落到 noop，后端什么都看不到。
3. **采样器**：用 `ParentBased(TraceIDRatioBased(r))`；只写 `TraceIDRatioBased` 会在链路中间产生"半截 trace"。
4. **resource 缺 service.name**：Jaeger/Tempo 里会显示 `unknown_service`，无法按服务筛选。
5. **传播器**：至少注册 `TraceContext` + `Baggage`，否则跨服务断链、baggage 也传不过去。
6. **OTel 内部错误默认静默**：用 `otel.SetErrorHandler` 接到 slog，否则 collector 不可达都发现不了。

## 测试提示（一个容易踩的坑）

`tracetest.InMemoryExporter` 的 `Shutdown()` **会清空已导出的数据**：

```go
tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(tracetest.NewInMemoryExporter()))
// ... 产生 span ...
tp.Shutdown(ctx) // 内部会 flush 并调用 exporter.Shutdown() → InMemoryExporter 被 Reset
exp.GetSpans()   // 结果是空的！
```

本项目测试用 `helpers_test.go` 里的 `memExporter`（`Shutdown` 时保留数据）来断言；
另一种做法是改用 `sdktrace.WithSyncer(exp)`（SimpleSpanProcessor，`span.End()` 时同步导出），
在 `Shutdown` 之前读取。

## 相关概念

- **Trace**：一次请求的调用链，由 Span 组成；Span 有 TraceID/SpanID/Parent，靠 `traceparent` header 跨服务传播。
- **Metric**：聚合数值，按时间窗口上报；Counter/Histogram 是同步指标，Observable 指标在采集时回调取值。
- **Propagation**：Inject（写 header）/ Extract（读 header），让下游能接着上游的 trace。
- **Resource**：描述数据的来源（服务名、版本、环境、主机……），trace 和 metric 都会携带。
