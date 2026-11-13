package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cmd/otel/pkg/telemetry"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type orderEvent struct {
	ID        string `json:"id"`
	Stock     int    `json:"stock"`
	CacheHit  bool   `json:"cache_hit"`
	CreatedAt string `json:"created_at"`
}

type app struct {
	tracer       trace.Tracer
	logger       *slog.Logger
	redis        *redis.Client
	writer       *kafka.Writer
	reader       *kafka.Reader
	httpClient   *http.Client
	inventoryURL string
	topic        string
	orders       metric.Int64Counter
	cacheHits    metric.Int64Counter
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serviceName := telemetry.Env("OTEL_SERVICE_NAME", "order-service")
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
	defer a.close()

	telemetry.StartPprof(telemetry.Env("PPROF_ADDR", ":6060"))
	go a.consume(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /order/{id}", a.handleOrder)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    ":" + telemetry.Env("PORT", "8080"),
		Handler: otelhttp.NewHandler(mux, "order-service"),
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
	orders, err := meter.Int64Counter("orders_total")
	if err != nil {
		return nil, err
	}
	cacheHits, err := meter.Int64Counter("cache_hits_total")
	if err != nil {
		return nil, err
	}
	brokers := strings.Split(telemetry.Env("KAFKA_BROKERS", "kafka:9092"), ",")
	topic := telemetry.Env("KAFKA_TOPIC", "orders")
	return &app{
		tracer: otel.Tracer(serviceName),
		logger: otelslog.NewLogger(serviceName),
		redis:  redis.NewClient(&redis.Options{Addr: telemetry.Env("REDIS_ADDR", "redis:6379")}),
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			RequiredAcks:           kafka.RequireOne,
			AllowAutoTopicCreation: true,
			BatchTimeout:           5 * time.Millisecond,
		},
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			GroupID:     "order-consumer",
			Topic:       topic,
			StartOffset: kafka.FirstOffset,
			MinBytes:    1,
			MaxBytes:    10e6,
		}),
		httpClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   5 * time.Second,
		},
		inventoryURL: telemetry.Env("INVENTORY_URL", "http://inventory-service:8080"),
		topic:        topic,
		orders:       orders,
		cacheHits:    cacheHits,
	}, nil
}

func (a *app) close() {
	_ = a.writer.Close()
	_ = a.reader.Close()
	_ = a.redis.Close()
}

func (a *app) handleOrder(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.tracer.Start(r.Context(), "handle_order")
	defer span.End()

	id := r.PathValue("id")
	if id == "" {
		id = "unknown"
	}

	ctx, getSpan := a.tracer.Start(ctx, "redis.get", trace.WithAttributes(
		attribute.String("db.system", "redis"),
		attribute.String("db.operation.name", "GET"),
		attribute.String("db.redis.key", "order:"+id),
	))
	cached, getErr := a.redis.Get(ctx, "order:"+id).Result()
	getSpan.End()
	if getErr == nil && cached != "" {
		var event orderEvent
		if err := json.Unmarshal([]byte(cached), &event); err == nil {
			event.CacheHit = true
			span.SetAttributes(attribute.Bool("cache.hit", true))
			a.orders.Add(ctx, 1, metric.WithAttributes(attribute.String("cache", "hit")))
			a.cacheHits.Add(ctx, 1)
			a.logger.InfoContext(ctx, "order served from cache", "order_id", id)
			writeJSON(w, event)
			return
		}
	}
	if getErr != nil && getErr != redis.Nil {
		span.RecordError(getErr)
		span.SetStatus(codes.Error, getErr.Error())
		a.logger.ErrorContext(ctx, "redis get failed", "error", getErr)
	}

	stock, err := a.fetchInventory(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	event := orderEvent{ID: id, Stock: stock, CacheHit: false, CreatedAt: time.Now().UTC().Format(time.RFC3339)}

	ctx, setSpan := a.tracer.Start(ctx, "redis.set", trace.WithAttributes(
		attribute.String("db.system", "redis"),
		attribute.String("db.operation.name", "SET"),
		attribute.String("db.redis.key", "order:"+id),
	))
	payload, _ := json.Marshal(event)
	if err := a.redis.Set(ctx, "order:"+id, payload, 5*time.Minute).Err(); err != nil {
		setSpan.RecordError(err)
		setSpan.SetStatus(codes.Error, err.Error())
		a.logger.ErrorContext(ctx, "redis set failed", "error", err)
	}
	setSpan.End()

	if err := a.publish(ctx, event); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.SetAttributes(attribute.Bool("cache.hit", false), attribute.Int("inventory.stock", stock))
	a.orders.Add(ctx, 1, metric.WithAttributes(attribute.String("cache", "miss")))
	a.logger.InfoContext(ctx, "order created", "order_id", id, "stock", stock)
	writeJSON(w, event)
}

func (a *app) fetchInventory(ctx context.Context, id string) (int, error) {
	ctx, span := a.tracer.Start(ctx, "inventory.request")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.inventoryURL+"/inventory/"+id, nil)
	if err != nil {
		return 0, err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, errInventoryStatus(resp.StatusCode)
	}
	var out struct {
		Stock int `json:"stock"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return 0, err
	}
	span.SetAttributes(attribute.Int("inventory.stock", out.Stock))
	return out.Stock, nil
}

func (a *app) publish(ctx context.Context, event orderEvent) error {
	ctx, span := a.tracer.Start(ctx, "kafka.produce", trace.WithSpanKind(trace.SpanKindProducer), trace.WithAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination.name", a.topic),
		attribute.String("messaging.operation.name", "publish"),
	))
	defer span.End()

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	headers := make([]kafka.Header, 0, len(carrier))
	for k, v := range carrier {
		headers = append(headers, kafka.Header{Key: k, Value: []byte(v)})
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := a.writer.WriteMessages(ctx, kafka.Message{
		Key:     []byte(event.ID),
		Value:   body,
		Headers: headers,
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (a *app) consume(ctx context.Context) {
	for {
		msg, err := a.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			time.Sleep(time.Second)
			continue
		}
		carrier := propagation.MapCarrier{}
		for _, h := range msg.Headers {
			carrier[h.Key] = string(h.Value)
		}
		msgCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)
		msgCtx, span := a.tracer.Start(msgCtx, "kafka.consume", trace.WithSpanKind(trace.SpanKindConsumer), trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination.name", a.topic),
			attribute.String("messaging.operation.name", "process"),
		))
		a.logger.InfoContext(msgCtx, "order event consumed", "order_id", string(msg.Key))
		_ = a.redis.Set(msgCtx, "consumed:"+string(msg.Key), "1", 10*time.Minute).Err()
		span.End()
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

type inventoryStatusError int

func (e inventoryStatusError) Error() string {
	return "inventory status " + http.StatusText(int(e))
}

func errInventoryStatus(code int) error {
	return inventoryStatusError(code)
}
