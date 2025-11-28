# Makefile for go-bisub project

# 变量定义
APP_NAME := go-bisub
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')
GIT_COMMIT := $(shell git rev-parse HEAD)

# 构建标志
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# 开发相关命令
.PHONY: dev
dev: ## 启动开发环境（热重载）
	@echo "Starting development server with hot reload..."
	air

.PHONY: run
run: ## 运行应用
	@echo "Running $(APP_NAME)..."
	go run $(LDFLAGS) cmd/server/main.go

.PHONY: build
build: ## 构建应用
	@echo "Building $(APP_NAME)..."
	go build $(LDFLAGS) -o bin/$(APP_NAME) cmd/server/main.go

.PHONY: build-linux
build-linux: ## 构建Linux版本
	@echo "Building $(APP_NAME) for Linux..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-linux cmd/server/main.go

.PHONY: build-windows
build-windows: ## 构建Windows版本
	@echo "Building $(APP_NAME) for Windows..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-windows.exe cmd/server/main.go

.PHONY: build-all
build-all: build build-linux build-windows ## 构建所有平台版本

# 代码质量检查
.PHONY: fmt
fmt: ## 格式化代码
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

.PHONY: lint
lint: ## 运行代码检查
	@echo "Running linters..."
	golangci-lint run

.PHONY: vet
vet: ## 运行go vet
	@echo "Running go vet..."
	go vet ./...

.PHONY: check
check: fmt vet lint ## 运行所有代码检查

# 测试相关命令
.PHONY: test
test: ## 运行测试
	@echo "Running tests..."
	go test -v -race -cover ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-integration
test-integration: ## 运行集成测试
	@echo "Running integration tests..."
	go test -v -tags=integration ./test/integration/...

.PHONY: benchmark
benchmark: ## 运行基准测试
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# 依赖管理
.PHONY: deps
deps: ## 下载依赖
	@echo "Downloading dependencies..."
	go mod download

.PHONY: deps-update
deps-update: ## 更新依赖
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

.PHONY: deps-verify
deps-verify: ## 验证依赖
	@echo "Verifying dependencies..."
	go mod verify

.PHONY: deps-clean
deps-clean: ## 清理依赖
	@echo "Cleaning dependencies..."
	go mod tidy

# Docker相关命令
.PHONY: docker-build
docker-build: ## 构建Docker镜像
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

.PHONY: docker-run
docker-run: ## 运行Docker容器
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

.PHONY: docker-compose-up
docker-compose-up: ## 启动Docker Compose服务
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## 停止Docker Compose服务
	@echo "Stopping services with Docker Compose..."
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## 查看Docker Compose日志
	@echo "Showing Docker Compose logs..."
	docker-compose logs -f

# 数据库相关命令
.PHONY: db-migrate-up
db-migrate-up: ## 执行数据库迁移
	@echo "Running database migrations..."
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/go_bisub" up

.PHONY: db-migrate-down
db-migrate-down: ## 回滚数据库迁移
	@echo "Rolling back database migrations..."
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/go_bisub" down

.PHONY: db-migrate-create
db-migrate-create: ## 创建新的迁移文件
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations $$name

# 清理命令
.PHONY: clean
clean: ## 清理构建文件
	@echo "Cleaning build files..."
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html

.PHONY: clean-all
clean-all: clean ## 清理所有生成文件
	@echo "Cleaning all generated files..."
	go clean -cache -modcache -testcache

# 安装工具
.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 生成相关命令
.PHONY: generate
generate: ## 运行go generate
	@echo "Running go generate..."
	go generate ./...

.PHONY: mock
mock: ## 生成mock文件
	@echo "Generating mock files..."
	mockgen -source=internal/repository/subscription.go -destination=internal/repository/mocks/subscription_mock.go
	mockgen -source=internal/service/subscription.go -destination=internal/service/mocks/subscription_mock.go

# 部署相关命令
.PHONY: deploy-staging
deploy-staging: ## 部署到测试环境
	@echo "Deploying to staging..."
	# 添加部署脚本

.PHONY: deploy-prod
deploy-prod: ## 部署到生产环境
	@echo "Deploying to production..."
	# 添加部署脚本

# 监控和日志
.PHONY: logs
logs: ## 查看应用日志
	@echo "Showing application logs..."
	tail -f logs/app.log

.PHONY: health
health: ## 检查应用健康状态
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || echo "Application is not healthy"

# 性能分析
.PHONY: profile-cpu
profile-cpu: ## CPU性能分析
	@echo "Starting CPU profiling..."
	go tool pprof http://localhost:8080/debug/pprof/profile

.PHONY: profile-mem
profile-mem: ## 内存性能分析
	@echo "Starting memory profiling..."
	go tool pprof http://localhost:8080/debug/pprof/heap

# 安全检查
.PHONY: security
security: ## 运行安全检查
	@echo "Running security checks..."
	gosec ./...

# 文档生成
.PHONY: docs
docs: ## 生成文档
	@echo "Generating documentation..."
	godoc -http=:6060
	@echo "Documentation server started at http://localhost:6060"

# 版本信息
.PHONY: version
version: ## 显示版本信息
	@echo "App Name: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"

# 初始化项目
.PHONY: init
init: deps install-tools ## 初始化项目环境
	@echo "Project initialized successfully!"
	@echo "Run 'make dev' to start development server"