# 多阶段构建 - 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata make

# 设置 Go 环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖（利用 Docker 缓存）
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建应用（添加版本信息）
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT
RUN go build -a -installsuffix cgo \
    -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o main cmd/server/main.go

# 运行阶段 - 使用更小的基础镜像
FROM alpine:3.19

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata wget curl && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 创建应用目录
WORKDIR /app

# 创建非 root 用户
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser && \
    chown -R appuser:appgroup /app

# 从构建阶段复制二进制文件
COPY --from=builder --chown=appuser:appgroup /app/main .

# 复制 web 静态文件和模板
COPY --from=builder --chown=appuser:appgroup /app/web ./web

# 复制配置文件（可以被 volume 覆盖）
COPY --from=builder --chown=appuser:appgroup /app/config.yaml .

# 创建日志目录
RUN mkdir -p /app/logs && chown -R appuser:appgroup /app/logs

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 设置环境变量
ENV GIN_MODE=release

# 运行应用
CMD ["./main"]