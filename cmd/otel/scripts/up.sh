#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

ARCH=arm64
[ "$(uname -m)" = "x86_64" ] && ARCH=amd64

mkdir -p bin
GOOS=linux GOARCH="$ARCH" go build -o bin/order-service ./service-a
GOOS=linux GOARCH="$ARCH" go build -o bin/inventory-service ./service-b

docker build -f Dockerfile --build-arg SERVICE=order-service -t otel-demo/order-service:dev .
docker build -f Dockerfile --build-arg SERVICE=inventory-service -t otel-demo/inventory-service:dev .

docker pull grafana/otel-lgtm:latest
docker pull grafana/alloy:latest
docker pull redis:7-alpine
docker pull apache/kafka:3.9.0

if ! kind get clusters 2>/dev/null | grep -q '^otel-demo$'; then
  kind create cluster --name otel-demo --config deploy/kind-config.yaml
fi

load_image() {
  local img="$1" tar
  tar="$(mktemp /tmp/kind-image-XXXXXX.tar)"
  docker save "$img" -o "$tar"
  for node in $(kind get nodes --name otel-demo); do
    docker exec -i "$node" ctr -n k8s.io images import --digests --snapshotter=overlayfs - < "$tar" >/dev/null
  done
  rm -f "$tar"
}

load_image grafana/otel-lgtm:latest
load_image grafana/alloy:latest
load_image redis:7-alpine
load_image apache/kafka:3.9.0
load_image otel-demo/order-service:dev
load_image otel-demo/inventory-service:dev

kubectl apply -f deploy/00-namespaces.yaml
kubectl apply -f deploy/lgtm.yaml
kubectl apply -f deploy/redis.yaml
kubectl apply -f deploy/kafka.yaml
kubectl -n demo rollout status deploy/redis --timeout=120s
kubectl -n demo rollout status deploy/kafka --timeout=240s
kubectl -n demo wait --for=condition=complete job/kafka-init --timeout=180s
kubectl apply -f deploy/services.yaml
kubectl apply -f deploy/alloy.yaml
kubectl -n observability rollout status deploy/otel-lgtm --timeout=240s
kubectl -n observability rollout status ds/alloy --timeout=120s
kubectl -n demo rollout status deploy/order-service --timeout=120s
kubectl -n demo rollout status deploy/inventory-service --timeout=120s
