# otel：Grafana LGTM 可观测性实践

一个可运行的 OpenTelemetry 端到端示例：两个 Go 服务（HTTP + Redis + Kafka）通过 OTLP 把 trace/metric/log 发给 Alloy，再进 LGTM 栈；profile 由 Alloy 抓 `/debug/pprof`。
目的只是跑通「一次请求 → 四类信号 → Grafana 查得到」，不是生产部署模板。

## 架构

```text
curl /order/{id} ──► order-service (service-a :8080) ──HTTP──► inventory-service (service-b :8080)
                          │ Redis / Kafka(orders)              │ Redis
                          │ OTLP :4317(gRPC) / :4318(HTTP)
                          ▼
                    Alloy DaemonSet ──OTLP──► otel-lgtm
                                               ├─ Grafana    :3000
                                               ├─ Loki       :3100
                                               ├─ Tempo      :3200
                                               ├─ Prometheus :9090
                                               └─ Pyroscope  :4040
                    Alloy Deployment（单副本）── pyroscope.scrape ──► Pyroscope
```

## 快速开始（kind 全链路）

```bash
brew install kind kubectl      # 需要 Docker Desktop，建议 ≥4GB 内存

cd cmd/otel
./scripts/up.sh                # 建集群、构建两个镜像、导入、apply、等就绪
kubectl port-forward -n demo svc/order-service 18080:8080
```

```bash
curl http://localhost:18080/order/1   # 连打两次，第二次命中 Redis 缓存
```

- Grafana：<http://localhost:13000>（NodePort 30300 映射到宿主机 13000），`admin` / `admin`
- NodePort 不通时：`kubectl port-forward -n observability svc/otel-lgtm 13001:3000`
- 销毁：`./scripts/down.sh`（删掉 kind 集群 `otel-demo`）

## 只起后端（不建 kind 集群）

```bash
docker run -d --name otel-lgtm \
  -p 13000:3000 -p 14040:4040 -p 24317:4317 -p 24318:4318 \
  grafana/otel-lgtm:latest
```

应用侧设置 `OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:24317` 即可（service-a/b 还需要可达的 Redis、Kafka）。

## 在 Grafana 里查什么

| 信号 | 数据源 | 查询 |
|---|---|---|
| trace | Tempo | `{resource.service.name="order-service" && name="GET /order/{id}"}` |
| log | Loki | `{service_name=~"order-service|inventory-service"}` |
| metric | Prometheus | `orders_total`、`cache_hits_total`、`inventory_requests_total` |
| profile | Pyroscope | `{service_name="order-service"}` |

一次 cache miss 的完整链路：`GET /order/{id}` → `handle_order` → `redis.get` → `inventory.request` → `inventory-service GET /inventory/{id}` → `redis.get/set` → `redis.set` → `kafka.produce` → `kafka.consume`（消费端从消息头恢复 traceparent）。

## 关键设计

- **Alloy 拆成 DaemonSet + Deployment**：DaemonSet 每节点一个，作为 OTLP 网关；pprof 抓取必须单副本，两个实例同时抓同一个进程会报 `cpu profiling already in use`，所以 profile 抓取放在单副本 Deployment。
- **profile 走 Alloy `pyroscope.scrape`**，不是应用内上报；应用只负责暴露 `:6060/debug/pprof`。

## 已知限制 / 坑

- otel-lgtm 没挂持久卷，重启容器后 Grafana / Loki / Tempo / Prometheus / Pyroscope 数据全清空。
- kind 节点容器默认拉不到 Docker Hub（容器内代理指向 127.0.0.1），`up.sh` 用 `docker save` + `ctr images import` 手动导入镜像。
- macOS + Docker Desktop 没有 eBPF：profile 只能走 pprof 拉取。镜像内置的 OBI（eBPF 自动埋点）需 `ENABLE_OBI=true --pid=host --privileged` 且 Linux 内核才可用，本示例未启用。
- **镜像版本差异**：`grafana/otel-lgtm:latest` 实测为 v0.40.0，metrics 后端是 **Prometheus**（容器内 :9090，Grafana 里以 `Prometheus` 数据源出现）；更早的镜像是 Mimir。两者在 Grafana 里都表现为 Prometheus 数据源，不存在名为 Mimir 的数据源。

## 文件

| 文件 | 内容 |
|---|---|
| `service-a/main.go` | order-service：HTTP、Redis、Kafka 生产/消费、自定义 counter |
| `service-b/main.go` | inventory-service：HTTP、Redis |
| `internal/telemetry/telemetry.go` | OTel 初始化：Trace/Metric/Log provider + OTLP gRPC exporter（读标准环境变量） |
| `deploy/lgtm.yaml` | `grafana/otel-lgtm`（NodePort 30300 → Grafana） |
| `deploy/alloy.yaml` | Alloy DaemonSet（OTLP 网关）+ 单副本 Deployment（pprof 抓取） |
| `deploy/services.yaml` | 两个服务（:8080、pprof :6060）+ `OTEL_EXPORTER_OTLP_ENDPOINT` |
| `deploy/{redis,kafka}.yaml` | Redis、Kafka 3.9 KRaft 单节点 |
| `scripts/up.sh` / `down.sh` | 一键起 / 删 |
| `Dockerfile` | linux 静态二进制 → alpine |

## 跑通验证

```bash
kubectl -n demo get pods -o wide            # 全 Running
kubectl -n observability get pods           # otel-lgtm / alloy / alloy-profiler
curl -s localhost:18080/healthz             # ok
```
