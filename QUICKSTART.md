# 快速启动指南

## 问题解决：make dev 找不到 air

如果遇到 `make dev` 提示 `air: No such file or directory`，有以下几种解决方案：

### 方案 1：使用开发脚本（推荐）

```bash
# 自动检查和安装依赖，然后启动
bash scripts/dev.sh

# 或使用 make 命令
make dev-script
```

### 方案 2：手动安装工具

```bash
# 安装所有开发工具
make install-tools

# 将 Go bin 目录添加到 PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# 然后启动开发服务器
make dev
```

### 方案 3：快速启动（无热重载）

```bash
# 直接运行，不需要 air
make start

# 或者
go run cmd/server/main.go
```

## 完整开发流程

### 1. 首次设置

#### 使用本地数据库（推荐）

```bash
# 1. 检查环境
make check-env

# 2. 安装开发工具
make install-tools

# 3. 下载依赖
make deps

# 4. 复制配置文件
cp config.local.yaml config.yaml
# 然后编辑 config.yaml，修改数据库密码等配置

# 5. 检查数据库连接
make db-check

# 6. 初始化数据库
make db-init
# 或手动执行：
# mysql -h 127.0.0.1 -u root -p < init.sql
```

#### 使用 Docker（可选）

```bash
# 1-3 步骤同上

# 4. 启动 Docker 服务
docker-compose up -d mysql redis

# 5. 等待服务启动
sleep 10

# 6. 初始化数据库
docker-compose exec mysql mysql -u root -ppassword < init.sql
```

### 2. 启动开发服务器

选择以下任一方式：

```bash
# 方式 1: 热重载开发（推荐）
make dev

# 方式 2: 使用脚本（自动检查依赖）
bash scripts/dev.sh

# 方式 3: 快速启动（无热重载）
make start
```

### 3. 访问服务

- **API 端点**: http://localhost:8080
- **管理界面**: http://localhost:8080/admin
  - 用户名: `admin`
  - 密码: `admin123`
- **健康检查**: http://localhost:8080/health

### 4. 开发工作流

```bash
# 修改代码后，air 会自动重新编译和重启服务

# 代码格式化
make fmt

# 代码检查
make lint

# 运行测试
make test

# 完整检查
make check
```

## 常用命令

### 开发相关

```bash
make help           # 查看所有可用命令
make check-env      # 检查开发环境
make dev            # 启动热重载开发服务器
make start          # 快速启动（无热重载）
make run            # 直接运行
```

### 代码质量

```bash
make fmt            # 格式化代码
make lint           # 代码检查
make vet            # 静态分析
make check          # 完整检查（fmt + vet + lint）
```

### 测试相关

```bash
make test           # 运行测试
make test-coverage  # 生成覆盖率报告
make benchmark      # 运行基准测试
```

### 构建部署

```bash
make build          # 构建当前平台
make build-all      # 构建所有平台
make docker-build   # 构建 Docker 镜像
```

### Docker 相关

```bash
make docker-compose-up      # 启动所有服务
make docker-compose-down    # 停止所有服务
make docker-compose-logs    # 查看日志
```

## 故障排查

### Air 找不到

**问题**: `make dev` 提示 `air: No such file or directory`

**解决**:
```bash
# 检查 air 是否安装
ls -la $(go env GOPATH)/bin/air

# 如果不存在，安装
go install github.com/air-verse/air@latest

# 添加到 PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# 或者使用完整路径
$(go env GOPATH)/bin/air
```

### 数据库连接失败

**问题**: 启动时提示数据库连接失败

**解决**:
```bash
# 检查 MySQL 是否运行
docker ps | grep mysql

# 如果没有运行，启动
docker-compose up -d mysql

# 等待启动完成
sleep 10

# 测试连接
mysql -h 127.0.0.1 -u root -ppassword -e "SELECT 1"
```

### Redis 连接失败

**问题**: 启动时提示 Redis 连接失败

**解决**:
```bash
# 检查 Redis 是否运行
docker ps | grep redis

# 如果没有运行，启动
docker-compose up -d redis

# 测试连接
redis-cli -h 127.0.0.1 ping
```

### 端口被占用

**问题**: 启动时提示端口 8080 被占用

**解决**:
```bash
# 查找占用端口的进程
lsof -i :8080

# 杀死进程
kill -9 <PID>

# 或者修改配置文件中的端口
vim config.yaml
# 修改 server.port 为其他端口
```

### 配置文件不存在

**问题**: 启动时提示找不到配置文件

**解决**:
```bash
# 从示例复制
cp config.yaml.example config.yaml

# 根据需要修改配置
vim config.yaml
```

## 开发技巧

### 1. 使用 Air 配置

Air 的配置文件是 `.air.toml`，可以自定义：

```toml
# 修改监听的文件类型
include_ext = ["go", "tpl", "tmpl", "html", "yaml"]

# 排除不需要监听的目录
exclude_dir = ["assets", "tmp", "vendor", "web/static"]

# 修改构建延迟（毫秒）
delay = 1000
```

### 2. 查看实时日志

```bash
# 使用 air 时，日志会实时显示在终端

# 如果使用 Docker，查看日志
docker-compose logs -f go-bisub

# 查看应用日志文件
tail -f logs/app.log
```

### 3. 调试技巧

```bash
# 使用 delve 调试器
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug cmd/server/main.go

# 或者在代码中添加日志
import "log"
log.Printf("Debug: %+v", variable)
```

### 4. 性能分析

```bash
# 启动服务后，访问 pprof 端点
# CPU 分析
go tool pprof http://localhost:8080/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:8080/debug/pprof/heap

# 或使用 make 命令
make profile-cpu
make profile-mem
```

## 推荐开发工具

- **IDE**: VSCode / GoLand
- **API 测试**: Postman / Insomnia / curl
- **数据库客户端**: DBeaver / MySQL Workbench
- **Redis 客户端**: RedisInsight / Another Redis Desktop Manager
- **Docker 管理**: Docker Desktop / Portainer

## 下一步

- 阅读 [README.md](README.md) 了解完整功能
- 查看 [API 文档](docs/) 了解接口详情
- 访问 Web 管理界面体验功能
- 开始开发你的第一个订阅服务

## 获取帮助

- 查看所有 Make 命令: `make help`
- 检查环境配置: `make check-env`
- 查看项目文档: `docs/` 目录
- 提交 Issue 或 PR
