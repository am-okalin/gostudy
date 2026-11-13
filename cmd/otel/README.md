# otel：Grafana LGTM 可观测性实践

一个可运行的 OpenTelemetry 端到端示例：两个 Go 服务（HTTP + Redis + Kafka）用 OTLP 上报
trace / metric / log，在 Grafana 里统一查询。目的是跑通「一次请求 → 三类信号 → 可查询」，
不是生产部署模板。

## 架构

```text
curl /order/{id} ──► order-service ──HTTP──► inventory-service
                         │ Redis / Kafka(orders)
                         │ OTLP :4317
                         ▼
                    otel-lgtm（单容器：Grafana + Tempo + Prometheus + Loki）
```

**没有 collector**：2 服务、6 实例的规模直连后端就够；collector 的价值在按 traceID 分片
（tail sampling / spanmetrics）、多后端扇出、节点级采集，规模上来再加。

## 快速开始

```bash
cd cmd/otel
docker compose up -d --build      # 首次构建约 1~2 分钟
curl http://localhost:18080/order/1   # 连打两次，第二次命中 Redis 缓存
```

- Grafana：<http://localhost:13000>，`admin` / `admin`
- 停止：`docker compose down`；连数据一起清：`docker compose down -v`
- 端口被占用（例如之前手动 `docker run` 过一个 otel-lgtm）：先 `docker rm -f <容器名>`，否则 13000/24317/24318 会冲突
- 想模拟多副本：`docker compose up -d --scale inventory-service=3`（order-service 有宿主端口映射，扩多个需先去掉 `ports`）

## 在 Grafana 里查什么

| 信号 | 数据源 | 查询 |
|---|---|---|
| trace | Tempo | `{resource.service.name="order-service" && name="GET /order/{id}"}` |
| log | Loki | `{service_name=~"order-service|inventory-service"}` |
| metric | Prometheus | `orders_total`、`cache_hits_total`、`inventory_requests_total` |

一次 cache miss 的完整链路：`GET /order/{id}` → `handle_order` → `redis.get` → `inventory.request` →
`inventory-service GET /inventory/{id}` → `redis.get/set` → `redis.set` → `kafka.produce` → `kafka.consume`
（消费端从消息头恢复 traceparent）。

## 关键设计

- **不装 collector**：这个规模直连后端最省；collector 不是 OTel 协议的必需品，是规模化/治理的需求。
- **换后端不改代码**：OTLP 出口只由 `OTEL_EXPORTER_OTLP_ENDPOINT` 决定，把值改成另一个 collector
  或厂商 endpoint 即可。
- **profile 不在本 demo 范围**：应用仍暴露 `:6060/debug/pprof`，可以 `go tool pprof` 手动分析；
  要自动采集需要额外引入 agent（Alloy / OTel Collector）或 `pyroscope-go` SDK，Pyroscope 自己不会拉取。

## 已知限制

- `otel-lgtm` 是单容器单副本：无 HA、无对象存储。数据写 `lgtm-data` volume，`docker compose down -v` 会清空。
- Kafka 是 KRaft 单节点，`orders` topic 由 `kafka-init` 一次性创建。
- 镜像版本差异：`grafana/otel-lgtm:latest` 实测 v0.40.0，metrics 后端是 **Prometheus**（容器内 :9090）；
  更早的镜像是 Mimir。Grafana 里统一显示为 `Prometheus` 数据源，没有名为 Mimir 的数据源。
- 历史路径已移除，需要时从 git 历史取：kind + K8s 清单在 `5287f3c` 的 tree 里
  （如 `git show 5287f3c:cmd/otel/deploy/alloy.yaml`），Alloy collector 的 compose override
  在 `9d1be18` 的 tree 里。

## 文件

```text
cmd/otel/
├── docker-compose.yml            # lgtm + redis + kafka + 两个服务
├── Dockerfile                    # 多阶段构建，compose 直接 build
├── README.md
├── pkg/telemetry/telemetry.go    # OTel 初始化（两个服务共用）
├── service-a/main.go             # order-service
└── service-b/main.go             # inventory-service
```

| 文件 | 内容 |
|---|---|
| `docker-compose.yml` | lgtm + redis + kafka + 两个服务（应用直连 lgtm）；`kafka-init` 一次性建 topic |
| `Dockerfile` | 多阶段构建（golang:1.27.3 → alpine），compose 直接 build |
| `service-a/main.go` | order-service：HTTP、Redis、Kafka 生产/消费、自定义 counter |
| `service-b/main.go` | inventory-service：HTTP、Redis |
| `pkg/telemetry/telemetry.go` | OTel 初始化：Resource、三套 Provider、OTLP exporter、W3C 传播 |

## 排障

```bash
docker compose ps                                  # 看各容器状态/健康
docker compose logs -f order-service               # 导出失败会每 2s 重试并打错误
curl -s localhost:18080/healthz                    # ok
curl -s localhost:13000/api/health                 # Grafana
# 确认 trace 到了 Tempo：刚写入的 span 在 ingester 里还没 flush，搜索接口会返回空，
# 所以看摄入计数更可靠
docker compose exec lgtm curl -s http://127.0.0.1:3200/metrics | grep tempo_distributor_spans_received_total
```

## 端口

| 服务 | 宿主机端口 |
|---|---|
| order-service HTTP | 18080 |
| Grafana | 13000 |
| OTLP gRPC / HTTP（本机调试用） | 24317 / 24318 |
