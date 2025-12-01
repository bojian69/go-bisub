# 本地开发环境配置指南

## 🎯 本地开发 vs Docker 环境

### 配置文件使用场景

| 场景 | 使用文件 | 数据库地址 | Redis 地址 |
|------|---------|-----------|-----------|
| **本地开发**（直接运行 Go） | `config.yaml` | `127.0.0.1` | `127.0.0.1` |
| **Docker 开发** | `config.docker.yaml` | `mysql` | `redis` |
| **Docker 生产** | `config.prod.yaml` | 实际地址 | 实际地址 |

## 📝 本地开发环境配置

### 方案一：使用本地 MySQL 和 Redis（推荐）

#### 1. 安装本地服务

```bash
# macOS
brew install mysql redis

# 启动服务
brew services start mysql
brew services start redis

# 验证服务
mysql -uroot -p
redis-cli ping
```

#### 2. 创建数据库

```bash
# 连接 MySQL
mysql -uroot -p

# 创建数据库
CREATE DATABASE go_sub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 创建用户（可选）
CREATE USER 'bisub'@'localhost' IDENTIFIED BY 'bisub123';
GRANT ALL PRIVILEGES ON go_sub.* TO 'bisub'@'localhost';
FLUSH PRIVILEGES;

# 导入初始化脚本
USE go_sub;
SOURCE init.sql;
SOURCE init_operation_logs.sql;
```

#### 3. 配置 config.yaml

```bash
# 复制模板
cp config.local.yaml config.yaml

# 编辑配置
vim config.yaml
```

**config.yaml 内容：**

```yaml
server:
  host: 0.0.0.0
  port: 8080
  timeout: 120s
  rate_limit: 1000

database:
  primary:
    host: 127.0.0.1  # ✅ 本地 MySQL
    port: 3306
    database: go_sub
    username: root
    password: ""  # 你的 MySQL root 密码
    max_idle_conns: 10
    max_open_conns: 100
    conn_max_lifetime: 3600s
  
  data_sources:
    default:
      host: 127.0.0.1  # ✅ 本地 MySQL
      port: 3306
      database: bi_data
      username: root
      password: ""
      max_idle_conns: 10
      max_open_conns: 100
      conn_max_lifetime: 3600s

redis:
  host: 127.0.0.1  # ✅ 本地 Redis
  port: 6379
  password: ""
  db: 0

security:
  jwt_secret: "sk-mviKoV-IGNWNRxK0SX6MXyj"
  allowed_sql_types:
    - "SELECT"

logging:
  level: "debug"  # 开发环境用 debug
  format: "json"
  file_log_enabled: true
  file_log_dir: "./logs"
  log_request_body: true
  log_response_body: true
  output: "stdout"

web_ui:
  username: "admin"
  password: "admin123"

snowflake:
  node_id: 1
```

#### 4. 启动应用

```bash
# 使用热重载
make dev

# 或直接运行
make run
```

### 方案二：使用远程数据库（当前配置）

如果你想连接远程数据库（如阿里云 RDS）：

#### config.yaml 配置：

```yaml
database:
  primary:
    host: testing-uhomes.rwlb.rds.aliyuncs.com  # 远程数据库
    port: 3306
    database: go_sub
    username: 
    password: 
    max_idle_conns: 10
    max_open_conns: 100
    conn_max_lifetime: 3600s
  
  data_sources:
    default:
      host: testing-uhomes.rwlb.rds.aliyuncs.com
      port: 3306
      database: uhomes
      username: 
      password: 
      max_idle_conns: 10
      max_open_conns: 100
      conn_max_lifetime: 3600s

redis:
  host: 127.0.0.1  # 本地 Redis
  port: 6379
  password: ""
  db: 0
```

### 方案三：混合模式（Docker 数据库 + 本地应用）

使用 Docker 运行数据库，本地运行应用：

#### 1. 启动 Docker 数据库

```bash
# 只启动数据库服务
docker-compose up -d mysql redis
```

#### 2. 配置 config.yaml

```yaml
database:
  primary:
    host: 127.0.0.1  # Docker 映射到本地
    port: 3306       # docker-compose.yml 中映射的端口
    database: go_sub
    username: root
    password: password  # .env 中的 MYSQL_ROOT_PASSWORD
    max_idle_conns: 10
    max_open_conns: 100
    conn_max_lifetime: 3600s

redis:
  host: 127.0.0.1  # Docker 映射到本地
  port: 6379       # docker-compose.yml 中映射的端口
  password: ""
  db: 0
```

#### 3. 启动应用

```bash
make dev
```

## 🔧 .env 文件说明

### ⚠️ 重要：.env 文件仅用于 Docker Compose

`.env` 文件**不会**被 Go 应用直接读取！它只用于配置 Docker Compose。

```bash
# ❌ 错误理解
# .env 中的 DB_HOST 不会被 Go 应用使用

# ✅ 正确理解
# .env 用于配置 Docker 容器
# Go 应用读取 config.yaml
```

### 本地开发不需要 .env

```bash
# 本地开发环境
├── config.yaml          # ✅ Go 应用读取这个
└── .env                 # ❌ 本地开发不需要

# Docker 环境
├── config.docker.yaml   # ✅ 挂载到容器内作为 config.yaml
└── .env                 # ✅ Docker Compose 读取这个
```

## 📊 配置对照表

### 数据库连接配置

| 场景 | 配置文件 | host | port | 说明 |
|------|---------|------|------|------|
| 本地 MySQL | config.yaml | `127.0.0.1` | `3306` | 本机安装的 MySQL |
| Docker MySQL（容器内） | config.docker.yaml | `mysql` | `3306` | Docker 服务名 |
| Docker MySQL（容器外） | config.yaml | `127.0.0.1` | `3306` | 端口映射到本机 |
| 远程 MySQL | config.yaml | `db.example.com` | `3306` | 实际数据库地址 |

### Redis 连接配置

| 场景 | 配置文件 | host | port | 说明 |
|------|---------|------|------|------|
| 本地 Redis | config.yaml | `127.0.0.1` | `6379` | 本机安装的 Redis |
| Docker Redis（容器内） | config.docker.yaml | `redis` | `6379` | Docker 服务名 |
| Docker Redis（容器外） | config.yaml | `127.0.0.1` | `6379` | 端口映射到本机 |
| 远程 Redis | config.yaml | `redis.example.com` | `6379` | 实际 Redis 地址 |

## 🎯 快速配置命令

### 本地 MySQL + Redis

```bash
# 1. 安装服务
brew install mysql redis

# 2. 启动服务
brew services start mysql
brew services start redis

# 3. 创建数据库
mysql -uroot -p -e "CREATE DATABASE go_sub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 4. 导入数据
mysql -uroot -p go_sub < init.sql
mysql -uroot -p go_sub < init_operation_logs.sql

# 5. 配置应用
cp config.local.yaml config.yaml
# 修改 config.yaml 中的密码

# 6. 启动应用
make dev
```

### Docker 数据库 + 本地应用

```bash
# 1. 启动数据库
docker-compose up -d mysql redis

# 2. 等待数据库就绪
sleep 10

# 3. 配置应用
cp config.local.yaml config.yaml
# 修改 host 为 127.0.0.1
# 修改 password 为 .env 中的 MYSQL_ROOT_PASSWORD

# 4. 启动应用
make dev
```

## 🔍 验证配置

### 测试数据库连接

```bash
# MySQL
mysql -h 127.0.0.1 -P 3306 -uroot -p

# 或使用 Docker
docker-compose exec mysql mysql -uroot -ppassword
```

### 测试 Redis 连接

```bash
# Redis
redis-cli -h 127.0.0.1 -p 6379 ping

# 或使用 Docker
docker-compose exec redis redis-cli ping
```

### 测试应用连接

```bash
# 启动应用
make dev

# 查看日志，确认连接成功
# 应该看到类似输出：
# INFO  Database connected successfully
# INFO  Redis connected successfully
```

## ⚠️ 常见问题

### Q1: 本地开发需要配置 .env 吗？

**A**: 不需要！`.env` 只用于 Docker Compose。

```bash
# 本地开发
只需要 config.yaml ✅

# Docker 部署
需要 .env + config.docker.yaml ✅
```

### Q2: 如何知道使用哪个配置文件？

**A**: 看你的运行方式

```bash
# 直接运行 Go 程序
make dev          # 读取 config.yaml

# Docker 运行
make docker-up    # 读取 config.docker.yaml（挂载为 config.yaml）
```

### Q3: 数据库连接失败？

**A**: 检查 host 配置

```yaml
# ❌ 本地开发错误配置
database:
  host: mysql  # 这是 Docker 服务名，本地无法解析

# ✅ 本地开发正确配置
database:
  host: 127.0.0.1  # 本地地址
```

### Q4: 如何在本地使用 Docker 的数据库？

**A**: 使用端口映射

```bash
# 1. 启动 Docker 数据库
docker-compose up -d mysql redis

# 2. 配置 config.yaml
database:
  host: 127.0.0.1  # 通过端口映射访问
  port: 3306       # Docker 映射到本地的端口

# 3. 启动应用
make dev
```

## 📝 配置模板

### 本地开发模板（config.yaml）

```yaml
# 本地开发配置
server:
  host: 0.0.0.0
  port: 8080
  timeout: 120s
  rate_limit: 1000

database:
  primary:
    host: 127.0.0.1      # 本地 MySQL
    port: 3306
    database: go_sub
    username: root
    password: ""         # 你的密码
    max_idle_conns: 10
    max_open_conns: 100
    conn_max_lifetime: 3600s

redis:
  host: 127.0.0.1        # 本地 Redis
  port: 6379
  password: ""
  db: 0

security:
  jwt_secret: "sk-mviKoV-IGNWNRxK0SX6MXyj"
  allowed_sql_types:
    - "SELECT"

logging:
  level: "debug"
  format: "json"
  file_log_enabled: true
  file_log_dir: "./logs"
  log_request_body: true
  log_response_body: true
  output: "stdout"

web_ui:
  username: "admin"
  password: "admin123"

snowflake:
  node_id: 1
```

## 🚀 推荐配置

### 新手推荐：Docker 数据库 + 本地应用

**优点**：
- ✅ 不需要安装 MySQL 和 Redis
- ✅ 数据隔离，不影响本机
- ✅ 一键启动数据库
- ✅ 应用可以热重载

**步骤**：
```bash
# 1. 启动数据库
docker-compose up -d mysql redis

# 2. 配置应用
cp config.local.yaml config.yaml
# host: 127.0.0.1
# password: password (来自 .env)

# 3. 启动应用
make dev
```

### 进阶推荐：本地 MySQL + Redis

**优点**：
- ✅ 性能更好
- ✅ 不依赖 Docker
- ✅ 可以使用 GUI 工具

**步骤**：
```bash
# 1. 安装服务
brew install mysql redis

# 2. 启动服务
brew services start mysql redis

# 3. 配置应用
cp config.local.yaml config.yaml

# 4. 启动应用
make dev
```

## 📚 相关文档

- [配置文件指南](./CONFIGURATION_GUIDE.md)
- [Docker 快速开始](../DOCKER_QUICKSTART.md)
- [Docker 部署指南](./DOCKER_DEPLOYMENT.md)
