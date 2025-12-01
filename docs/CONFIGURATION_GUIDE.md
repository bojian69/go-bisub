# 配置文件指南

## 📋 配置文件对比

### `.env.example` vs `config.local.yaml`

| 特性 | .env.example | config.local.yaml |
|------|--------------|-------------------|
| **用途** | Docker Compose 环境变量 | 应用程序配置文件 |
| **格式** | KEY=VALUE | YAML 结构化配置 |
| **作用域** | Docker 容器和服务 | Go 应用程序内部 |
| **优先级** | 低（被 config.yaml 覆盖） | 高（应用直接读取） |
| **适用场景** | Docker 部署 | 本地开发 |

### 配置层次结构

```
优先级（从高到低）：
1. 环境变量（ENV）
2. config.yaml（应用配置）
3. config.local.yaml（本地开发模板）
4. .env（Docker Compose 变量）
5. .env.example（模板）
```

## 🐳 Docker 环境配置

### 方案一：使用 .env + config.yaml（推荐）

**适用场景**：生产环境、Docker 部署

#### 步骤：

1. **创建 .env 文件**
```bash
cp .env.example .env
vim .env
```

2. **创建 config.docker.yaml**
```yaml
server:
  host: 0.0.0.0
  port: 8080
  timeout: 120s
  rate_limit: 1000

database:
  primary:
    host: mysql  # Docker 服务名
    port: 3306
    database: go_sub
    username: root
    password: password  # 从 .env 读取
    max_idle_conns: 10
    max_open_conns: 100
    conn_max_lifetime: 3600s
  
  data_sources:
    default:
      host: mysql
      port: 3306
      database: bi_data
      username: root
      password: password
      max_idle_conns: 10
      max_open_conns: 100
      conn_max_lifetime: 3600s

redis:
  host: redis  # Docker 服务名
  port: 6379
  password: ""
  db: 0

security:
  jwt_secret: "sk-mviKoV-IGNWNRxK0SX6MXyj"
  allowed_sql_types:
    - "SELECT"

logging:
  level: "info"
  format: "json"
  file_log_enabled: true
  file_log_dir: "/app/logs"
  log_request_body: false
  log_response_body: false
  output: "stdout"

web_ui:
  username: "admin"
  password: "admin123"

snowflake:
  node_id: 1
```

3. **更新 docker-compose.yml**
```yaml
services:
  go-bisub:
    volumes:
      - ./config.docker.yaml:/app/config.yaml:ro
      - ./logs:/app/logs
```

### 方案二：纯环境变量（简化）

**适用场景**：简单部署、容器化环境

修改应用代码支持环境变量覆盖：

```go
// 从环境变量读取配置
if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
    config.Database.Primary.Host = dbHost
}
```

## 📝 配置文件详解

### .env.example（Docker Compose 变量）

```bash
# ============================================
# Docker Compose 环境变量
# ============================================

# 1. 应用版本信息
VERSION=v1.0.0              # 应用版本号
BUILD_TIME=                 # 构建时间（自动生成）
GIT_COMMIT=                 # Git 提交哈希（自动生成）

# 2. MySQL 容器配置
MYSQL_ROOT_PASSWORD=password    # MySQL root 密码
MYSQL_USER=bisub               # 创建的用户名
MYSQL_PASSWORD=bisub123        # 用户密码
MYSQL_DATABASE=go_sub          # 默认数据库

# 3. Redis 容器配置
REDIS_PASSWORD=                # Redis 密码（可选）

# 4. 应用数据库连接（在容器内使用）
DB_HOST=mysql                  # 使用 Docker 服务名
DB_PORT=3306
DB_USER=root
DB_PASS=password
DB_NAME=go_sub

# 5. 应用 Redis 连接（在容器内使用）
REDIS_HOST=redis               # 使用 Docker 服务名
REDIS_PORT=6379
REDIS_PASS=

# 6. 应用配置
SERVER_PORT=8080
GIN_MODE=release               # release/debug
LOG_LEVEL=info                 # debug/info/warn/error
JWT_SECRET=your-secret-key     # 生产环境必须修改
SNOWFLAKE_NODE_ID=1
TZ=Asia/Shanghai
```

### config.local.yaml（应用配置）

```yaml
# ============================================
# 应用程序配置文件
# ============================================

# 1. 服务器配置
server:
  host: 0.0.0.0              # 监听地址
  port: 8080                 # 监听端口
  timeout: 120s              # 请求超时
  rate_limit: 1000           # 速率限制

# 2. 数据库配置
database:
  primary:                   # 主数据库
    host: 127.0.0.1         # 本地开发用 localhost
    port: 3306
    database: go_sub
    username: root
    password: ""
    max_idle_conns: 10
    max_open_conns: 100
    conn_max_lifetime: 3600s
  
  data_sources:              # 数据源配置
    default:
      host: 127.0.0.1
      port: 3306
      database: bi_data
      username: root
      password: ""
      max_idle_conns: 10
      max_open_conns: 100
      conn_max_lifetime: 3600s

# 3. Redis 配置
redis:
  host: 127.0.0.1           # 本地开发用 localhost
  port: 6379
  password: ""
  db: 0

# 4. 安全配置
security:
  jwt_secret: "your-secret-key"
  allowed_sql_types:
    - "SELECT"

# 5. 日志配置
logging:
  level: "debug"            # 开发环境用 debug
  format: "json"
  file_log_enabled: true
  file_log_dir: "./logs"
  log_request_body: true
  log_response_body: true
  output: "stdout"

# 6. Web UI 配置
web_ui:
  username: "admin"
  password: "admin123"

# 7. Snowflake ID 配置
snowflake:
  node_id: 1
```

## 🔧 不同环境的配置

### 本地开发环境

```bash
# 1. 使用 config.local.yaml
cp config.local.yaml config.yaml

# 2. 修改数据库连接为本地
vim config.yaml
# host: 127.0.0.1

# 3. 启动应用
make dev
```

### Docker 开发环境

```bash
# 1. 创建 .env
cp .env.example .env

# 2. 创建 config.docker.yaml
# 使用 Docker 服务名（mysql, redis）

# 3. 启动服务
make docker-up
```

### 生产环境

```bash
# 1. 创建 .env（生产配置）
cp .env.example .env
vim .env
# 修改所有密码和密钥

# 2. 创建 config.prod.yaml
# 使用生产数据库地址
# 关闭调试日志
# 启用安全配置

# 3. 部署
docker-compose -f docker-compose.prod.yml up -d
```

## 📊 配置对照表

### 数据库连接

| 环境 | Host | 说明 |
|------|------|------|
| 本地开发 | `127.0.0.1` 或 `localhost` | 本机 MySQL |
| Docker | `mysql` | Docker Compose 服务名 |
| 生产环境 | `db.example.com` | 实际数据库地址 |

### Redis 连接

| 环境 | Host | 说明 |
|------|------|------|
| 本地开发 | `127.0.0.1` 或 `localhost` | 本机 Redis |
| Docker | `redis` | Docker Compose 服务名 |
| 生产环境 | `redis.example.com` | 实际 Redis 地址 |

### 日志级别

| 环境 | Level | 说明 |
|------|-------|------|
| 本地开发 | `debug` | 详细日志 |
| Docker 开发 | `info` | 一般日志 |
| 生产环境 | `warn` 或 `error` | 只记录警告和错误 |

## 🔒 安全最佳实践

### 1. 密码管理

```bash
# ❌ 不要在代码中硬编码密码
password: "123456"

# ✅ 使用环境变量
password: ${DB_PASSWORD}

# ✅ 使用密钥管理服务
# AWS Secrets Manager
# HashiCorp Vault
# Kubernetes Secrets
```

### 2. JWT Secret

```bash
# ❌ 不要使用默认值
jwt_secret: "your-secret-key"

# ✅ 生成强随机密钥
jwt_secret: "$(openssl rand -base64 32)"
```

### 3. 配置文件权限

```bash
# 设置只读权限
chmod 600 config.yaml
chmod 600 .env

# 不要提交到 Git
echo "config.yaml" >> .gitignore
echo ".env" >> .gitignore
```

## 🎯 推荐配置方案

### Docker 部署（推荐）

```
项目结构：
├── .env                    # Docker Compose 变量（不提交）
├── .env.example            # 环境变量模板（提交）
├── config.docker.yaml      # Docker 配置（提交）
├── config.local.yaml       # 本地开发模板（提交）
└── docker-compose.yml      # Docker 编排（提交）

使用方式：
1. cp .env.example .env
2. 修改 .env 中的密码
3. docker-compose up -d
4. 应用读取 config.docker.yaml
```

### 本地开发

```
项目结构：
├── config.yaml             # 实际配置（不提交）
├── config.local.yaml       # 配置模板（提交）
└── .env.example            # 环境变量模板（提交）

使用方式：
1. cp config.local.yaml config.yaml
2. 修改 config.yaml 中的数据库连接
3. make dev
```

## 🔄 配置迁移

### 从本地开发迁移到 Docker

```bash
# 1. 创建 Docker 配置
cp config.local.yaml config.docker.yaml

# 2. 修改主机地址
sed -i 's/127.0.0.1/mysql/g' config.docker.yaml
sed -i 's/127.0.0.1/redis/g' config.docker.yaml

# 3. 创建环境变量
cp .env.example .env

# 4. 启动 Docker
make docker-up
```

## 📚 相关文档

- [Docker 部署指南](./DOCKER_DEPLOYMENT.md)
- [本地开发指南](./LOCAL_DEVELOPMENT.md)
- [环境变量参考](./ENVIRONMENT_VARIABLES.md)
