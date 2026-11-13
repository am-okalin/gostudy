package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cmd/otel/pkg/telemetry"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type inventory struct {
	ID       string `json:"id"`
	Stock    int    `json:"stock"`
	CacheHit bool   `json:"cache_hit"`
}

type app struct {
	tracer    trace.Tracer
	logger    *slog.Logger
	redis     *redis.Client
	requests  metric.Int64Counter
	cacheHits metric.Int64Counter
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serviceName := telemetry.Env("OTEL_SERVICE_NAME", "inventory-service")
	shutdown, err := telemetry.Init(ctx, telemetry.Config{
		ServiceName:    serviceName,
		ServiceVersion: telemetry.Env("SERVICE_VERSION", "dev"),
		Environment:    telemetry.Env("DEPLOY_ENV", "dev"),
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(closeCtx)
	}()

	a, err := newApp(serviceName)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = a.redis.Close() }()

	telemetry.StartPprof(telemetry.Env("PPROF_ADDR", ":6060"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /inventory/{id}", a.handleInventory)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    ":" + telemetry.Env("PORT", "8080"),
		Handler: otelhttp.NewHandler(mux, "inventory-service"),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	<-ctx.Done()
	closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(closeCtx)
}

func newApp(serviceName string) (*app, error) {
	meter := otel.Meter(serviceName)
	requests, err := meter.Int64Counter("inventory_requests_total")
	if err != nil {
		return nil, err
	}
	cacheHits, err := meter.Int64Counter("inventory_cache_hits_total")
	if err != nil {
		return nil, err
	}
	return &app{
		tracer:    otel.Tracer(serviceName),
		logger:    otelslog.NewLogger(serviceName),
		redis:     redis.NewClient(&redis.Options{Addr: telemetry.Env("REDIS_ADDR", "redis:6379")}),
		requests:  requests,
		cacheHits: cacheHits,
	}, nil
}

func (a *app) handleInventory(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "handle_inventory")
	defer span.End()

	id := r.PathValue("id")
	if id == "" {
		id = "unknown"
	}

	ctx, getSpan := a.tracer.Start(ctx, "redis.get", trace.WithAttributes(
		attribute.String("db.system", "redis"),
		attribute.String("db.operation.name", "GET"),
		attribute.String("db.redis.key", "stock:"+id),
	))
	cached, getErr := a.redis.Get(ctx, "stock:"+id).Result()
	getSpan.End()
	if getErr == nil && cached != "" {
		var inv inventory
		if err := json.Unmarshal([]byte(cached), &inv); err == nil {
			inv.CacheHit = true
			span.SetAttributes(attribute.Bool("cache.hit", true))
			a.requests.Add(ctx, 1, metric.WithAttributes(attribute.String("cache", "hit")))
			a.cacheHits.Add(ctx, 1)
			a.logger.InfoContext(ctx, "inventory served from cache", "item_id", id)
			writeJSON(w, inv)
			return
		}
	}
	if getErr != nil && getErr != redis.Nil {
		span.RecordError(getErr)
		span.SetStatus(codes.Error, getErr.Error())
		a.logger.ErrorContext(ctx, "redis get failed", "error", getErr)
	}

	inv := inventory{ID: id, Stock: rand.IntN(100), CacheHit: false}
	payload, _ := json.Marshal(inv)
	ctx, setSpan := a.tracer.Start(ctx, "redis.set", trace.WithAttributes(
		attribute.String("db.system", "redis"),
		attribute.String("db.operation.name", "SET"),
		attribute.String("db.redis.key", "stock:"+id),
	))
	if err := a.redis.Set(ctx, "stock:"+id, payload, 5*time.Minute).Err(); err != nil {
		setSpan.RecordError(err)
		setSpan.SetStatus(codes.Error, err.Error())
		a.logger.ErrorContext(ctx, "redis set failed", "error", err)
	}
	setSpan.End()

	span.SetAttributes(attribute.Bool("cache.hit", false), attribute.Int("inventory.stock", inv.Stock))
	a.requests.Add(ctx, 1, metric.WithAttributes(attribute.String("cache", "miss")))
	a.logger.InfoContext(ctx, "inventory computed", "item_id", id, "stock", inv.Stock)
	writeJSON(w, inv)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
