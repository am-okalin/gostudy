# 多阶段构建：容器内编译 linux 静态二进制，`docker compose build` 直接可用。
# 注意 build context 必须是 cmd/（cmd/otel 属于 cmd 这个 Go module），
# 所以 compose 里写的是 `context: ..` + `dockerfile: otel/Dockerfile`。
FROM golang:1.27.3 AS build

ARG SERVICE_DIR
ARG GOPROXY=https://goproxy.cn,direct

ENV CGO_ENABLED=0 GOOS=linux GOWORK=off GOPROXY=${GOPROXY}

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /out/app ./otel/${SERVICE_DIR}

FROM alpine:latest
RUN adduser -D -u 10001 app
USER app
COPY --from=build /out/app /app
ENTRYPOINT ["/app"]
