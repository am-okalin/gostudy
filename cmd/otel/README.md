# OTel 可观测性实践：kind + Alloy + LGTM

## 架构

```text
curl /order/{id}
      │
      ▼
order-service ──HTTP──► inventory-service
      │  │                    │
      │  └──► Redis          └──► Redis
      │
      └──► Kafka(topic=orders) ──► order-service consumer
      │
      │ OTLP(trace/metric/log)
      ▼
Alloy DaemonSet（每节点一个）──► otel-lgtm
                                   ├─ Tempo
                                   ├─ Loki
                                   ├─ Mimir
                                   └─ Pyroscope
Alloy Profiler Deployment（单副本）──► Pyroscope
      ▲
      └── 定时 scrape /debug/pprof（order-service / inventory-service）
```

## 组件

| 组件 | 说明 |
|---|---|
| order-service | HTTP 入口，读写 Redis，调用 inventory-service，生产/消费 Kafka |
| inventory-service | 库存服务，读写 Redis |
| redis | 缓存 |
| kafka | Apache Kafka 3.9 KRaft 单节点 |
| alloy DaemonSet | 接收 OTLP（trace/metric/log），批量转发到 otel-lgtm |
| alloy-profiler Deployment | 单副本，`pyroscope.scrape` 抓取两个服务的 pprof 并写入 Pyroscope |
| otel-lgtm | Grafana + Loki + Tempo + Mimir + Pyroscope |

## 前置

```bash
brew install kind kubectl
```

Docker Desktop 建议至少 4GB 内存。本次实践使用 kind v1.37、apache/kafka:3.9.0。

## 启动

```bash
cd cmd/otel
./scripts/up.sh
kubectl port-forward -n demo svc/order-service 18080:8080
```

另开终端：

```bash
curl http://localhost:18080/order/1
curl http://localhost:18080/order/1
```

Grafana：<http://localhost:13000>（NodePort 30300），账号 `admin` / `admin`。

如果 NodePort 不可达，用 port-forward 代替：

```bash
kubectl port-forward -n observability svc/otel-lgtm 13001:3000
```

然后访问 <http://localhost:13001>。

## 查看数据

| 信号 | 数据源 | 查询 |
|---|---|---|
| trace | Tempo | `{resource.service.name="order-service" && name="GET /order/{id}"}` |
| log | Loki | `{service_name=~"order-service|inventory-service"}` |
| metric | Mimir | `orders_total`、`inventory_requests_total` |
| profile | Pyroscope | `{service_name="order-service"}`，profile type `process_cpu:cpu:nanoseconds:cpu:nanoseconds` |

一次 cache miss 的 trace 包含：`GET /order/{id}` → `handle_order` → `redis.get` →
`inventory.request` → `inventory-service GET /inventory/{id}` → `redis.get/set` →
`redis.set` → `kafka.produce`，Kafka 消费端通过消息头恢复 traceparent，产生
`kafka.consume` 子 span。

## 为什么 Alloy 分成 DaemonSet 和 Deployment

- DaemonSet 负责 OTLP 网关：每个节点一个，应用发到 `alloy.observability.svc:4317`。
- pprof 抓取是单例：两个 DaemonSet 实例同时 scrape 同一个 Go 进程会触发
  `Could not enable CPU profiling: cpu profiling already in use`，所以 profile
  抓取放到单副本 Deployment，避免重复抓取。

## 关于 eBPF

本次在 macOS + Docker Desktop 的 kind 环境里没有启用 eBPF，profile 走
Alloy `pyroscope.scrape` 拉取 `/debug/pprof`。在 Linux 集群上，可在 DaemonSet 里
加 `pyroscope.ebpf`（需要 hostPID、privileged 与内核支持），实现零代码采集。

## 注意

- otel-lgtm 在 kind 里没有挂持久卷，重启 pod 后 Grafana/Loki/Tempo/Mimir/Pyroscope 数据会清空。
- kind 的节点容器默认无法从 Docker Hub 拉镜像（容器内代理配置指向 127.0.0.1），
  `scripts/up.sh` 用 `docker save` + `ctr images import` 手动导入镜像。
- macOS/Docker Desktop 下未启用 eBPF；profile 通过 Alloy scrape `/debug/pprof` 获取。

## 销毁

```bash
./scripts/down.sh
```
