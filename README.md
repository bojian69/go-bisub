# go-sub

GO BI Subscription 订阅BI数据服务

> 🎯 **新手？** 从 [这里开始](docs/START_HERE.md) | 📚 [命令速查表](docs/COMMANDS.md)

## 🚀 快速启动

```bash
# 1. 安装依赖
make install-tools && make deps

# 2. 配置数据库
cp config.local.yaml config.yaml
# 编辑 config.yaml 修改数据库密码

# 3. 初始化数据库
make db-init

# 4. 启动服务
make dev
```

访问：
- **API**: http://localhost:8080
- **管理界面**: http://localhost:8080/admin (admin/admin123)
- **健康检查**: http://localhost:8080/health

详细说明请查看 [快速开始](#快速开始) 章节。

> 💡 **提示**: 查看 [命令速查表](docs/COMMANDS.md) 快速查找常用命令

---

## 📋 目录

- [快速启动](#-快速启动)
- [功能特性](#功能特性)
- [快速开始](#快速开始)
  - [使用 Docker Compose](#使用docker-compose推荐)
  - [本地开发](#本地开发)
- [常用开发命令](#常用开发命令)
- [API文档](#api文档)
- [项目启动流程](#项目启动流程)
- [故障排查](#故障排查)
- [配置说明](#配置说明)
- [数据库表结构](#数据库表结构)
- [部署](#部署)
- [开发](#开发)
- [文档](#文档)

## 功能特性

- ✅ 订阅管理：创建、查询、更新订阅服务
- ✅ 版本控制：支持多版本订阅，自动版本选择
- ✅ SQL执行：安全的SQL执行引擎，支持变量替换
- ✅ 多数据源：支持配置多个数据库连接
- ✅ 异步统计：不影响API响应的统计数据收集
- ✅ 限流保护：基于Redis的分布式限流
- ✅ 认证授权：JWT认证和基础认证
- ✅ Web管理界面：简单的订阅管理界面
- ✅ 操作日志：完整的操作审计日志记录
- ✅ 容器化部署：Docker支持

## 快速开始

### 环境要求
- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0+
- Redis 6.0+

### 使用Docker Compose（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd go-bisub

# 启动所有服务
make docker-compose-up
# 或者
docker-compose up -d

# 查看日志
make docker-compose-logs
# 或者
docker-compose logs -f go-bisub
```

### 本地开发

#### 方式 1: 使用本地 MySQL 和 Redis（推荐）

```bash
# 1. 检查开发环境
make check-env

# 2. 安装开发工具和依赖
make install-tools
make deps

# 3. 复制并配置
cp config.local.yaml config.yaml
# 编辑 config.yaml，修改数据库密码等配置

# 4. 检查数据库连接
make db-check

# 5. 初始化数据库
make db-init

# 6. 启动开发服务器
make dev                # 热重载（推荐）
# 或
bash scripts/dev.sh     # 使用脚本（自动检查依赖）
```

#### 方式 2: 使用 Docker

```bash
# 1-2 步骤同上

# 3. 启动 Docker 服务
docker-compose up -d mysql redis

# 4. 初始化数据库
mysql -h 127.0.0.1 -u root -ppassword < init.sql

# 5. 启动开发服务器
make dev
```

**其他启动方式：**
```bash
make start              # 快速启动（无热重载）
go run cmd/server/main.go  # 直接运行
```

**注意事项：**
- 如果 `make dev` 提示找不到 `air`，请先运行 `make install-tools`
- 开发工具会安装到 `$GOPATH/bin`，确保该目录在 PATH 中
- 或者将 `$(go env GOPATH)/bin` 添加到 PATH：
  ```bash
  export PATH="$PATH:$(go env GOPATH)/bin"
  ```

### 常用开发命令

#### 启动相关
```bash
make dev                # 启动开发服务器（热重载，推荐）
make start              # 快速启动（无热重载）
bash scripts/dev.sh     # 使用脚本启动（自动检查依赖）
go run cmd/server/main.go  # 直接运行
```

#### 数据库相关
```bash
make db-check           # 检查数据库连接
make db-init            # 初始化数据库
```

#### 开发工具
```bash
make help               # 查看所有可用命令
make check-env          # 检查开发环境
make install-tools      # 安装开发工具
make deps               # 下载依赖
```

#### 代码质量
```bash
make check              # 完整检查（格式化+检查+测试）
make fmt                # 代码格式化
make lint               # 代码检查
make test               # 运行测试
make test-coverage      # 生成测试覆盖率报告
```

#### 构建部署
```bash
make build              # 构建应用
make build-all          # 构建所有平台版本
make docker-build       # 构建 Docker 镜像
make docker-compose-up  # 启动 Docker 服务
```

#### 其他
```bash
make health             # 查看应用健康状态
make logs               # 查看应用日志
make clean              # 清理构建文件
```

## API文档

### 认证

所有API请求需要在Header中包含JWT Token：

```
Authorization: Bearer <your-jwt-token>
```

### 订阅管理

#### 创建订阅

```bash
POST /v1/subscriptions
Content-Type: application/json

{
  "type": "A",
  "sub_key": "house_report",
  "version": 1,
  "title": "房源报表",
  "abstract": "获取房源基本信息",
  "status": "B",
  "extra_config": {
    "sql_content": "SELECT * FROM houses WHERE id = house_id_replace",
    "sql_replace": {"house_id_replace": "房源ID"},
    "example": "{\"house_id_replace\": 1}"
  }
}
```

#### 获取订阅列表

```bash
GET /v1/subscriptions?limit=20&offset=0
```

#### 获取订阅详情

```bash
GET /v1/subscriptions/{key}
GET /v1/subscriptions/{key}/versions/{version}
```

### 订阅执行

#### 执行订阅（默认版本）

```bash
POST /v1/subscriptions/house_report:execute
Content-Type: application/json

{
  "variables": {
    "house_id_replace": 1
  },
  "timeout": 30000,
  "data_source": "default"
}
```

#### 执行特定版本

```bash
POST /v1/subscriptions/house_report/versions/1:execute
Content-Type: application/json

{
  "variables": {
    "house_id_replace": 1
  }
}
```

### 统计查询

```bash
GET /v1/subscriptions/stats?start_time=2025-01-01&end_time=2025-01-31&limit=20&offset=0
```

### 操作日志

#### 获取操作日志

```bash
GET /v1/operation-logs?start_time=2025-01-01&end_time=2025-01-31&operation=CREATE&limit=20&offset=0
```

支持的查询参数：
- `start_time`: 开始时间 (YYYY-MM-DD)
- `end_time`: 结束时间 (YYYY-MM-DD)
- `user_id`: 用户ID
- `username`: 用户名（模糊匹配）
- `operation`: 操作类型 (CREATE/UPDATE/DELETE/EXECUTE/QUERY)
- `resource`: 资源类型（模糊匹配）
- `status`: 操作状态 (SUCCESS/FAILED)
- `client_ip`: 客户端IP
- `limit`: 每页数量 (默认20，最大100)
- `offset`: 偏移量 (默认0)

## Web管理界面

访问 `http://localhost:8080/admin` 使用Web界面管理订阅。

默认账号密码：
- 用户名：admin
- 密码：admin123

## 配置说明

主要配置项在 `config.yaml` 中：

```yaml
server:
  port: 8080              # 服务端口
  timeout: 120s           # 请求超时
  rate_limit: 1000        # 限流QPS

database:
  primary:                # 主数据库（存储订阅信息）
    host: localhost
    port: 3306
    database: go_sub
    username: root
    password: password
  
  data_sources:           # 数据源配置
    default:
      host: localhost
      port: 3306
      database: bi_data
      username: readonly
      password: password

security:
  jwt_secret: your-jwt-secret
  allowed_sql_types: ["SELECT"]  # 允许的SQL类型

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

web_ui:
  username: admin
  password: admin123
```

## 数据库表结构

### 订阅表 (sub_subscription_theme)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT UNSIGNED | 主键ID |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| type | CHAR(1) | 订阅类型 A:分析数据 |
| sub_key | VARCHAR(120) | 订阅key |
| version | TINYINT UNSIGNED | 版本号 |
| title | VARCHAR(240) | 订阅标题 |
| abstract | TINYTEXT | 订阅简介 |
| status | CHAR(1) | 状态 A:待生效 B:生效中 C:生效中-强制兼容低版本 D:已失效 |
| created_by | BIGINT UNSIGNED | 创建人ID |
| extra_config | JSON | 扩展配置(sql_content,sql_replace,example) |

### 统计表 (sub_logs_bidata_response)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT UNSIGNED | 主键ID |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| sub_key | VARCHAR(120) | 订阅key |
| version | TINYINT UNSIGNED | 订阅版本号 |
| execution_duration | MEDIUMINT UNSIGNED | 执行耗时(毫秒) |
| request_url | VARCHAR(1000) | 请求链接 |
| request_response | JSON | 请求详情 |
| instance_source | VARCHAR(120) | 数据实例标识 |

### 操作日志表 (sub_logs_operation)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT UNSIGNED | 主键ID |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| user_id | BIGINT UNSIGNED | 操作用户ID |
| username | VARCHAR(120) | 操作用户名 |
| operation | VARCHAR(50) | 操作类型 |
| resource | VARCHAR(200) | 操作资源 |
| resource_id | VARCHAR(120) | 资源ID |
| status | VARCHAR(20) | 操作状态 |
| client_ip | VARCHAR(45) | 客户端IP |
| user_agent | VARCHAR(500) | 用户代理 |
| request_url | VARCHAR(1000) | 请求URL |
| method | VARCHAR(10) | HTTP方法 |
| duration | INT UNSIGNED | 执行耗时(毫秒) |
| error_msg | TEXT | 错误信息 |
| request_data | JSON | 请求数据 |
| response_data | JSON | 响应数据 |

## 部署

### Docker部署

```bash
# 构建镜像
docker build -t go-bisub:latest .

# 运行容器
docker run -d \
  --name go-bisub \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/root/config.yaml \
  go-bisub:latest
```

### Kubernetes部署

参考 `k8s/` 目录下的YAML文件。

## 开发

### 项目结构

```
go-bisub/
├── cmd/server/          # 主程序入口
├── internal/            # 私有应用代码
│   ├── config/         # 配置管理
│   ├── models/         # 数据模型（支持分布式ID）
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   ├── handler/        # HTTP处理器
│   ├── middleware/     # 中间件（认证、限流、日志）
│   └── utils/          # 工具函数（分布式ID生成）
├── web/                # Web管理界面
│   ├── static/         # 静态资源
│   └── templates/      # HTML模板
├── docs/               # 项目文档
├── config.yaml         # 配置文件
├── docker-compose.yml  # Docker Compose配置
├── Dockerfile          # Docker镜像构建
├── Makefile           # 开发命令集合
├── .golangci.yml      # 代码质量检查配置
├── .air.toml          # 热重载配置
└── init.sql           # 数据库初始化脚本
```

### 添加新功能

1. **定义数据模型** (`internal/models/`)
   ```bash
   # 使用分布式ID的基础模型
   type NewModel struct {
       BaseModel  # 自动包含分布式ID
       // 其他字段
   }
   ```

2. **实现数据访问** (`internal/repository/`)
   ```bash
   # 实现Repository接口
   type NewRepository interface {
       Create(ctx context.Context, model *NewModel) error
       // 其他方法
   }
   ```

3. **实现业务逻辑** (`internal/service/`)
   ```bash
   # 实现Service接口
   type NewService interface {
       CreateNew(ctx context.Context, req *CreateRequest) (*Response, error)
       // 其他方法
   }
   ```

4. **实现HTTP接口** (`internal/handler/`)
   ```bash
   # 实现HTTP处理器
   func (h *NewHandler) Create(c *gin.Context) {
       // HTTP处理逻辑
   }
   ```

5. **注册路由** (`cmd/server/main.go`)
   ```bash
   # 注册API路由
   v1.POST("/new", handler.Create)
   ```

### 开发工具和命令

#### 项目初始化
```bash
# 安装所有开发工具
make install-tools

# 初始化项目环境（包含依赖下载和工具安装）
make init
```

#### 代码质量检查
```bash
# 运行完整的代码检查流程
make check

# 单独运行各项检查
make fmt      # 代码格式化
make vet      # 静态分析
make lint     # 代码质量检查
```

#### 测试相关
```bash
# 运行所有测试
make test

# 运行测试（带竞争检测）
make test-race

# 生成测试覆盖率报告
make test-coverage

# 运行基准测试
make benchmark
```

#### 构建和部署
```bash
# 构建当前平台版本
make build

# 构建所有平台版本
make build-all

# 构建Docker镜像
make docker-build
```

#### 开发调试
```bash
# 启动热重载开发服务器
make dev

# 查看应用日志
make logs

# 检查应用健康状态
make health

# CPU性能分析
make profile-cpu

# 内存性能分析
make profile-mem
```

## 监控和日志

- 健康检查：`GET /health`
- 日志格式：JSON结构化日志
- 指标收集：支持Prometheus格式指标
- 操作审计：完整的用户操作日志记录
- 性能监控：API响应时间和执行统计

## 代码质量

### 开发前检查
```bash
# 运行完整的代码检查流程
make check

# 运行测试确保功能正常
make test
```

### 提交前检查清单
```bash
# 1. 代码格式化和质量检查
make check

# 2. 运行所有测试
make test

# 3. 确保构建成功
make build

# 4. 安全检查（可选）
make security
```

### 持续集成
项目配置了完整的代码质量检查工具：
- **golangci-lint**: 30+种代码检查规则
- **测试覆盖率**: 自动生成覆盖率报告
- **安全扫描**: gosec安全漏洞检测
- **性能分析**: 内置pprof性能分析

## 安全考虑

- SQL注入防护：严格的SQL验证和变量替换
- 认证授权：JWT Token认证
- 限流保护：防止API滥用
- 输入验证：所有输入参数验证
- 操作审计：完整的操作日志记录和追踪
- 数据安全：敏感数据脱敏和加密存储

## 技术特性

### 分布式ID生成
- **Snowflake算法**: 高性能分布式ID生成
- **并发安全**: 支持高并发场景下的ID唯一性
- **时间排序**: ID包含时间戳信息，天然排序
- **故障降级**: 自动降级到UUID v7

### 性能优化
- **连接池**: 数据库连接池优化
- **Redis缓存**: 分布式缓存支持
- **批量操作**: 支持批量数据处理
- **异步统计**: 不阻塞主流程的统计收集

### 安全特性
- **JWT认证**: 标准JWT Token认证
- **SQL注入防护**: 参数化查询和SQL验证
- **限流保护**: Redis分布式限流
- **操作审计**: 完整的操作日志记录

### 监控运维
- **健康检查**: `/health` 端点
- **指标收集**: Prometheus格式指标
- **结构化日志**: JSON格式日志输出
- **性能分析**: 内置pprof支持

## 项目启动流程

### 首次启动（本地开发）

```
┌─────────────────────────────────────────────────────────────┐
│ 1. 环境准备                                                  │
│    make check-env                                           │
│    make install-tools                                       │
│    make deps                                                │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. 配置文件                                                  │
│    cp config.local.yaml config.yaml                         │
│    编辑 config.yaml（修改数据库密码等）                      │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. 数据库初始化                                              │
│    make db-check    # 检查连接                              │
│    make db-init     # 初始化数据库                          │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. 启动服务                                                  │
│    make dev         # 热重载开发                            │
│    或 bash scripts/dev.sh                                   │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. 访问服务                                                  │
│    http://localhost:8080        - API                      │
│    http://localhost:8080/admin  - 管理界面                  │
│    http://localhost:8080/health - 健康检查                  │
└─────────────────────────────────────────────────────────────┘
```

### 日常开发

```bash
# 1. 启动服务
make dev

# 2. 修改代码（Air 会自动重新编译）

# 3. 测试 API
curl http://localhost:8080/health

# 4. 提交前检查
make check
make test
```

### 使用 Docker

```bash
# 1. 启动所有服务（包括数据库）
docker-compose up -d

# 2. 查看日志
docker-compose logs -f go-bisub

# 3. 停止服务
docker-compose down
```

## 故障排查

### 常见问题

| 问题 | 解决方案 |
|------|---------|
| `air: command not found` | 运行 `make install-tools` |
| `MySQL connection failed` | 检查 MySQL 是否启动，运行 `make db-check` |
| `Redis connection failed` | 检查 Redis 是否启动 |
| `Database not found` | 运行 `make db-init` 初始化数据库 |
| `Port 8080 already in use` | 修改 `config.yaml` 中的端口或杀死占用进程 |

详细故障排查请查看 [快速启动指南](docs/QUICKSTART.md) 和 [本地开发指南](docs/LOCAL_DEVELOPMENT.md)

## 文档

### 快速参考
- [新手入门](docs/START_HERE.md) - 从这里开始 🎯
- [命令速查表](docs/COMMANDS.md) - 常用命令快速查找 ⭐
- [快速启动指南](docs/QUICKSTART.md) - 详细的启动步骤和故障排查
- [本地开发指南](docs/LOCAL_DEVELOPMENT.md) - 本地开发环境配置
- [数据库迁移指南](docs/DATABASE_MIGRATION.md) - 数据库变更说明

### 技术文档
- [更新日志](docs/CHANGELOG.md)
- [操作日志实现](docs/OPERATION_LOGS_IMPLEMENTATION.md)

## 许可证

MIT License