# 命令速查表

快速查找常用命令。

## 🚀 启动命令

| 命令 | 说明 | 使用场景 |
|------|------|---------|
| `make dev` | 热重载开发 | 日常开发（推荐） |
| `make start` | 快速启动 | 无需热重载时 |
| `bash scripts/dev.sh` | 脚本启动 | 自动检查依赖 |
| `go run cmd/server/main.go` | 直接运行 | 简单测试 |

## 📦 环境准备

| 命令 | 说明 |
|------|------|
| `make check-env` | 检查开发环境 |
| `make install-tools` | 安装开发工具 |
| `make deps` | 下载依赖 |
| `make init` | 完整初始化（工具+依赖） |

## 🗄️ 数据库操作

| 命令 | 说明 |
|------|------|
| `make db-check` | 检查数据库连接 |
| `make db-init` | 初始化数据库 |
| `mysql -h 127.0.0.1 -u root -p go_sub < init.sql` | 手动导入 SQL |
| `mysql -h 127.0.0.1 -u root -p go_sub` | 连接数据库 |

## 🧪 测试命令

| 命令 | 说明 |
|------|------|
| `make test` | 运行测试 |
| `make test-race` | 运行测试（带竞争检测） |
| `make test-coverage` | 生成覆盖率报告 |
| `make benchmark` | 运行基准测试 |

## 🔍 代码质量

| 命令 | 说明 |
|------|------|
| `make check` | 完整检查（fmt+vet+lint） |
| `make fmt` | 格式化代码 |
| `make lint` | 代码检查 |
| `make vet` | 静态分析 |

## 🏗️ 构建命令

| 命令 | 说明 |
|------|------|
| `make build` | 构建当前平台 |
| `make build-linux` | 构建 Linux 版本 |
| `make build-windows` | 构建 Windows 版本 |
| `make build-all` | 构建所有平台 |

## 🐳 Docker 命令

| 命令 | 说明 |
|------|------|
| `make docker-build` | 构建 Docker 镜像 |
| `make docker-compose-up` | 启动所有服务 |
| `make docker-compose-down` | 停止所有服务 |
| `make docker-compose-logs` | 查看日志 |
| `docker-compose up -d mysql redis` | 只启动数据库 |

## 🧹 清理命令

| 命令 | 说明 |
|------|------|
| `make clean` | 清理构建文件 |
| `make clean-all` | 清理所有生成文件 |

## 📊 监控命令

| 命令 | 说明 |
|------|------|
| `make health` | 检查应用健康 |
| `make logs` | 查看应用日志 |
| `make profile-cpu` | CPU 性能分析 |
| `make profile-mem` | 内存性能分析 |

## 🔧 实用脚本

| 脚本 | 说明 |
|------|------|
| `bash scripts/dev.sh` | 开发环境启动 |
| `bash scripts/check-db.sh` | 数据库检查 |
| `bash scripts/init-local-db.sh` | 数据库初始化 |

## 📝 常用组合

### 首次启动
```bash
make install-tools && make deps
cp config.local.yaml config.yaml
# 编辑 config.yaml
make db-init
make dev
```

### 日常开发
```bash
make dev
# 修改代码...
make check
make test
```

### 提交前检查
```bash
make fmt
make lint
make test
git add .
git commit -m "your message"
```

### 完整测试
```bash
make check
make test
make test-coverage
make benchmark
```

### Docker 部署
```bash
make docker-build
docker-compose up -d
docker-compose logs -f
```

## 🆘 故障排查

### Air 找不到
```bash
make install-tools
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 数据库连接失败
```bash
make db-check
# 检查 MySQL 是否启动
brew services list | grep mysql
brew services start mysql
```

### 端口被占用
```bash
lsof -i :8080
kill -9 <PID>
# 或修改 config.yaml 中的端口
```

### 依赖问题
```bash
go mod tidy
go mod download
make deps
```

## 📚 更多信息

- 完整文档: [README.md](../README.md)
- 快速开始: [QUICKSTART.md](QUICKSTART.md)
- 本地开发: [LOCAL_DEVELOPMENT.md](LOCAL_DEVELOPMENT.md)
- 数据库迁移: [DATABASE_MIGRATION.md](DATABASE_MIGRATION.md)

## 💡 提示

- 使用 `make help` 查看所有可用命令
- 使用 `make check-env` 检查环境配置
- 开发时推荐使用 `make dev` 获得热重载体验
- 提交代码前运行 `make check` 确保代码质量
