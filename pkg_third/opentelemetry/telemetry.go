// Package opentelemetry 演示 OpenTelemetry（trace + metric）在 Go 服务中的初始化与使用。
//
// 本目录是真实项目 internal/pkg/telemetry 的简化示例，重点演示：
//
//  1. resource：描述"遥测数据来自哪个服务"，service.name 必填，否则后端显示 unknown_service；
//  2. TracerProvider：Sampler（ParentBased + TraceIDRatioBased）、BatchSpanProcessor、SpanExporter；
//  3. MeterProvider：PeriodicReader 周期上报指标；
//  4. 传播（Propagation）：W3C tracecontext + baggage，保证跨服务链路不断；
//  5. Shutdown：必须带超时 ctx，把缓冲区里的 span/metric 冲刷出去，否则最后一批数据会丢。
//
// 为了让示例在没有 collector 的机器上也能跑测试，Options 提供了 SpanExporter / MetricReader
// 两个注入点：测试里用内存导出器（tracetest.NewInMemoryExporter）和手动 reader
// （sdkmetric.NewManualReader），生产代码不需要设置它们。
//
// 依赖版本与真实项目保持一致（OTel Go v1.38.0）。运行方式见 README.md。
package opentelemetry

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Options 是 Init 的配置项。
type Options struct {
	ServiceName    string        // 写进 resource，是 Jaeger/Tempo 里筛选服务的主键。
	ServiceVersion string        // 写进 resource，便于对比版本间的延迟变化。
	Environment    string        // 例如 local/dev/prod。
	Endpoint       string        // OTLP gRPC 端点(host:port)；为空表示不接 collector。
	Insecure       bool          // true 表示明文 gRPC，仅用于本地/已受信网络。
	SampleRatio    float64       // 采样比例 [0,1]；ParentBased 会跟随上游的采样决定。
	ExportTimeout  time.Duration // 单次导出超时；<=0 时用默认 10s。
	MetricInterval time.Duration // 指标周期上报间隔；<=0 时用默认 30s。
	ConsoleOutput  bool          // true 时额外把 span 打到 stdout，本地零依赖调试用。
	ConsoleWriter  io.Writer     // stdout 输出目标，默认 os.Stdout；测试可换成 bytes.Buffer。

	// 以下两个字段用于测试或接入自定义导出器；为 nil 时按 Endpoint / ConsoleOutput 自动创建。
	SpanExporter sdktrace.SpanExporter // 例如 tracetest.NewInMemoryExporter()
	MetricReader sdkmetric.Reader      // 例如 sdkmetric.NewManualReader()
}

// Provider 持有 OTel 的 Provider 与关闭逻辑。
type Provider struct {
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	shutdown       func(context.Context) error
	logger         *slog.Logger
}

// Init 初始化 OpenTelemetry 并注册为全局 Provider；返回的 *Provider 必须在退出前 Shutdown。
func Init(ctx context.Context, opts Options, logger *slog.Logger) (*Provider, error) {
	// ---- 参数校验与默认值 ----
	if opts.ServiceName == "" {
		return nil, errors.New("telemetry: ServiceName must not be empty")
	}
	if opts.SampleRatio < 0 || opts.SampleRatio > 1 {
		return nil, fmt.Errorf("telemetry: SampleRatio must be within [0,1], got %v", opts.SampleRatio)
	}
	if opts.ExportTimeout <= 0 {
		opts.ExportTimeout = 10 * time.Second
	}
	if opts.MetricInterval <= 0 {
		opts.MetricInterval = 30 * time.Second
	}
	if opts.ConsoleWriter == nil {
		opts.ConsoleWriter = os.Stdout
	}
	if logger == nil {
		logger = slog.Default()
	}

	// ---- resource：数据来自谁 ----
	res, err := buildResource(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("telemetry: build resource: %w", err)
	}

	// ---- trace ----
	tracerProvider, err := newTracerProvider(ctx, opts, res)
	if err != nil {
		return nil, err
	}

	// ---- metric ----
	meterProvider, err := newMeterProvider(ctx, opts, res)
	if err != nil {
		// trace 已建好但 metric 失败：先清理已分配资源再返回，避免泄漏。
		_ = tracerProvider.Shutdown(ctx)
		return nil, err
	}

	// ---- 注册全局（otelgrpc/otelhttp 等库默认从全局取 provider，不设置会落到 noop）----
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	// 同时支持 W3C tracecontext 与 baggage，保证跨服务不断链、baggage 也能透传。
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	// OTel 内部的导出错误默认是静默丢弃的，这里接到 slog，便于发现 collector 不可达等问题。
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.ErrorContext(ctx, "opentelemetry internal error", slog.String("error", err.Error()))
	}))

	logger.InfoContext(ctx, "telemetry initialized",
		slog.String("service", opts.ServiceName),
		slog.String("endpoint", opts.Endpoint),
		slog.Bool("insecure", opts.Insecure),
		slog.Float64("sample_ratio", opts.SampleRatio),
		slog.Bool("console_output", opts.ConsoleOutput),
	)

	p := &Provider{
		tracerProvider: tracerProvider,
		meterProvider:  meterProvider,
		logger:         logger,
	}
	p.shutdown = func(shutdownCtx context.Context) error {
		var errs []error
		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown tracer provider: %w", err))
		}
		if err := meterProvider.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown meter provider: %w", err))
		}
		return errors.Join(errs...)
	}
	return p, nil
}

// newTracerProvider 组装 trace 导出器、采样器与 resource。
func newTracerProvider(ctx context.Context, opts Options, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	exporters := make([]sdktrace.SpanExporter, 0, 3)

	if opts.SpanExporter != nil {
		exporters = append(exporters, opts.SpanExporter)
	}
	if opts.Endpoint != "" {
		traceOpts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(opts.Endpoint),
			otlptracegrpc.WithTimeout(opts.ExportTimeout),
		}
		if opts.Insecure {
			traceOpts = append(traceOpts, otlptracegrpc.WithInsecure())
		}
		otlpTrace, err := otlptracegrpc.New(ctx, traceOpts...)
		if err != nil {
			return nil, fmt.Errorf("telemetry: create otlp trace exporter: %w", err)
		}
		exporters = append(exporters, otlpTrace)
	}
	if opts.ConsoleOutput {
		console, err := stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(), // 输出稳定，便于对比/贴日志
			stdouttrace.WithWriter(opts.ConsoleWriter),
		)
		if err != nil {
			return nil, fmt.Errorf("telemetry: create console trace exporter: %w", err)
		}
		exporters = append(exporters, console)
	}

	providerOpts := []sdktrace.TracerProviderOption{
		// ParentBased 跟随上游的采样决定；纯 TraceIDRatioBased 会在链路中间产生半截 trace。
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(opts.SampleRatio))),
		sdktrace.WithResource(res),
	}
	for _, exp := range exporters {
		// 每个 exporter 一个 BatchSpanProcessor：攒批导出，减少网络开销。
		providerOpts = append(providerOpts, sdktrace.WithBatcher(exp, sdktrace.WithExportTimeout(opts.ExportTimeout)))
	}
	return sdktrace.NewTracerProvider(providerOpts...), nil
}

// newMeterProvider 组装 metric reader。
func newMeterProvider(ctx context.Context, opts Options, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	var reader sdkmetric.Reader
	switch {
	case opts.MetricReader != nil:
		reader = opts.MetricReader
	case opts.Endpoint != "":
		metricOpts := []otlpmetricgrpc.Option{
			otlpmetricgrpc.WithEndpoint(opts.Endpoint),
			otlpmetricgrpc.WithTimeout(opts.ExportTimeout),
		}
		if opts.Insecure {
			metricOpts = append(metricOpts, otlpmetricgrpc.WithInsecure())
		}
		otlpMetric, err := otlpmetricgrpc.New(ctx, metricOpts...)
		if err != nil {
			return nil, fmt.Errorf("telemetry: create otlp metric exporter: %w", err)
		}
		reader = sdkmetric.NewPeriodicReader(otlpMetric,
			sdkmetric.WithInterval(opts.MetricInterval),
			sdkmetric.WithTimeout(opts.ExportTimeout),
		)
	}

	if reader == nil {
		// 没有 endpoint 也没有注入 reader：指标只记录不导出（no-op），演示环境下可接受。
		return sdkmetric.NewMeterProvider(sdkmetric.WithResource(res)), nil
	}
	return sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	), nil
}

// Shutdown 冲刷并关闭所有 exporter。
// 必须传入带超时的 ctx：collector 不可达时会阻塞，用 Background 会让进程退不掉。
func (p *Provider) Shutdown(ctx context.Context) error {
	if p == nil || p.shutdown == nil {
		return nil
	}
	if err := p.shutdown(ctx); err != nil {
		return err
	}
	p.logger.InfoContext(ctx, "telemetry shut down")
	return nil
}

// TracerProvider 暴露底层 provider，供需要显式构造 tracer 的场景使用。
func (p *Provider) TracerProvider() *sdktrace.TracerProvider {
	if p == nil {
		return nil
	}
	return p.tracerProvider
}

// MeterProvider 暴露底层 provider，供需要显式构造 meter 的场景使用。
func (p *Provider) MeterProvider() *sdkmetric.MeterProvider {
	if p == nil {
		return nil
	}
	return p.meterProvider
}

// buildResource 描述遥测数据来自谁；缺少 service.name 时 Jaeger 会显示 unknown_service。
func buildResource(ctx context.Context, opts Options) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithFromEnv(),      // OTEL_RESOURCE_ATTRIBUTES / OTEL_SERVICE_NAME
		resource.WithProcess(),      // 进程信息
		resource.WithOS(),           // 操作系统信息
		resource.WithContainer(),    // 容器信息（非容器环境自动跳过）
		resource.WithHost(),         // 主机名
		resource.WithTelemetrySDK(), // telemetry.sdk.* 属性
		resource.WithAttributes(
			semconv.ServiceName(opts.ServiceName),
			semconv.ServiceVersion(opts.ServiceVersion),
			attribute.String("deployment.environment", opts.Environment),
		),
	)
}

// Hostname 返回当前主机名，便于在自定义 resource 属性里补充实例标识。
func Hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}
