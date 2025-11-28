.PHONY: build run test clean docker-build docker-run dev

# 变量定义
APP_NAME=go-bisub
VERSION=latest
DOCKER_IMAGE=$(APP_NAME):$(VERSION)

# 构建应用
build:
	go build -o bin/$(APP_NAME) cmd/server/main.go

# 运行应用
run:
	go run cmd/server/main.go

# 运行测试
test:
	go test -v ./...

# 清理构建文件
clean:
	rm -rf bin/

# 安装依赖
deps:
	go mod tidy
	go mod download

# 代码格式化
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# Docker构建
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Docker运行
docker-run:
	docker run -d --name $(APP_NAME) -p 8080:8080 $(DOCKER_IMAGE)

# 开发环境启动
dev:
	docker-compose up -d mysql redis
	go run cmd/server/main.go

# 完整开发环境
dev-full:
	docker-compose up -d

# 停止开发环境
dev-stop:
	docker-compose down

# 查看日志
logs:
	docker-compose logs -f go-bisub

# 数据库迁移
migrate:
	docker-compose exec mysql mysql -uroot -ppassword go_sub < init.sql

# 生成JWT Token（用于测试）
token:
	@echo "生成测试JWT Token..."
	@go run -ldflags="-s -w" -o /tmp/token-gen - <<< 'package main; import ("fmt"; "time"; "github.com/golang-jwt/jwt/v5"); func main() { token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 1, "username": "admin", "exp": time.Now().Add(24*time.Hour).Unix()}); tokenString, _ := token.SignedString([]byte("your-jwt-secret-change-in-production")); fmt.Println("JWT Token:"); fmt.Println(tokenString) }' && /tmp/token-gen

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  build      - 构建应用"
	@echo "  run        - 运行应用"
	@echo "  test       - 运行测试"
	@echo "  clean      - 清理构建文件"
	@echo "  deps       - 安装依赖"
	@echo "  fmt        - 代码格式化"
	@echo "  lint       - 代码检查"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  docker-run - 运行Docker容器"
	@echo "  dev        - 启动开发环境（仅数据库）"
	@echo "  dev-full   - 启动完整开发环境"
	@echo "  dev-stop   - 停止开发环境"
	@echo "  logs       - 查看应用日志"
	@echo "  migrate    - 运行数据库迁移"
	@echo "  token      - 生成测试JWT Token"
	@echo "  help       - 显示帮助信息"