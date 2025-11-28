# go-bisub

GO BI Subscription 订阅BI数据服务

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

### 使用Docker Compose（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd go-bisub

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f go-bisub
```

### 本地开发

```bash
# 安装依赖
go mod tidy

# 启动MySQL和Redis
docker-compose up -d mysql redis

# 运行应用
go run cmd/server/main.go
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
    database: go_bisub
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
| status | CHAR(1) | 状态 A:待生效 B:生效中 C:强制兼容 D:已失效 |
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
├── internal/
│   ├── config/         # 配置管理
│   ├── models/         # 数据模型
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   ├── handler/        # HTTP处理器
│   └── middleware/     # 中间件
├── web/                # Web界面资源
├── config.yaml         # 配置文件
├── docker-compose.yml  # Docker Compose配置
├── Dockerfile          # Docker镜像构建
└── init.sql           # 数据库初始化脚本
```

### 添加新功能

1. 在 `internal/models/` 中定义数据模型
2. 在 `internal/repository/` 中实现数据访问
3. 在 `internal/service/` 中实现业务逻辑
4. 在 `internal/handler/` 中实现HTTP接口
5. 在 `cmd/server/main.go` 中注册路由

### 代码格式化和审查

#### 基础命令
```bash
# 格式化代码
go fmt ./...

# 清理依赖
go mod tidy

# 检查代码问题
go vet ./...

# 编译检查
go build ./...
```

#### 高级工具
```bash
# 安装goimports（自动管理导入包）
go install golang.org/x/tools/cmd/goimports@latest

# 使用goimports格式化代码
goimports -w .

# 安装golangci-lint（综合代码检查）
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 运行代码检查
golangci-lint run
```

#### 命令说明
- **go fmt**: 格式化Go代码，统一代码风格
- **go mod tidy**: 清理未使用的依赖，添加缺失的依赖
- **go vet**: 检查代码中的常见错误和问题
- **goimports**: 自动管理import语句，移除未使用的导入
- **golangci-lint**: 运行多种代码检查工具，发现潜在问题

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

# 或者手动执行
go fmt ./...
go mod tidy
go vet ./...
go build ./...
```

### 提交前检查
```bash
# 确保代码质量
goimports -w .
golangci-lint run
go test ./...
```

## 安全考虑

- SQL注入防护：严格的SQL验证和变量替换
- 认证授权：JWT Token认证
- 限流保护：防止API滥用
- 输入验证：所有输入参数验证
- 操作审计：完整的操作日志记录和追踪
- 数据安全：敏感数据脱敏和加密存储

## 许可证

MIT License