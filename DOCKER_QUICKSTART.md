# Docker 快速开始指南

## 🚀 5 分钟快速部署

### 1. 准备配置文件

```bash
# 复制环境变量文件
cp .env.example .env

# 编辑环境变量（可选，使用默认值也可以）
vim .env
```

### 2. 启动服务

```bash
# 一键启动所有服务
make docker-up

# 或使用 docker-compose
docker-compose up -d
```

### 3. 验证服务

```bash
# 查看服务状态
make docker-ps

# 查看日志
make docker-logs

# 健康检查
curl http://localhost:8080/health
```

### 4. 访问应用

- **Web UI**: http://localhost:8080/admin
  - 用户名: `admin`
  - 密码: `admin123`

- **API**: http://localhost:8080/v1/subscriptions

## 📋 配置文件说明

### 两个关键文件

#### 1. `.env` - Docker Compose 环境变量

```bash
# 用途：配置 Docker 容器和服务
# 位置：项目根目录
# 提交：❌ 不提交到 Git（包含密码）

# 主要配置：
MYSQL_ROOT_PASSWORD=password    # MySQL 密码
REDIS_PASSWORD=                 # Redis 密码（可选）
JWT_SECRET=your-secret-key      # JWT 密钥
```

#### 2. `config.docker.yaml` - 应用配置

```yaml
# 用途：Go 应用程序的配置
# 位置：项目根目录
# 提交：✅ 可以提交（不含敏感信息）

# 主要配置：
database:
  primary:
    host: mysql  # ⚠️ 使用 Docker 服务名，不是 localhost
    
redis:
  host: redis    # ⚠️ 使用 Docker 服务名，不是 localhost
```

## 🔄 配置对比

### 本地开发 vs Docker

| 配置项 | 本地开发 | Docker |
|--------|----------|--------|
| 配置文件 | `config.yaml` | `config.docker.yaml` |
| 数据库 Host | `127.0.0.1` | `mysql` |
| Redis Host | `127.0.0.1` | `redis` |
| 日志目录 | `./logs` | `/app/logs` |
| 日志级别 | `debug` | `info` |

### 为什么要用不同的 Host？

```yaml
# ❌ 本地开发（直接运行）
database:
  host: 127.0.0.1  # 连接本机 MySQL

# ✅ Docker 环境
database:
  host: mysql      # 连接 Docker 容器中的 MySQL
```

**原因**：
- Docker 容器有独立的网络
- 容器内的 `127.0.0.1` 指向容器自己，不是宿主机
- Docker Compose 创建了一个网络，服务通过服务名互相访问

## 🎯 常见场景

### 场景 1：首次部署

```bash
# 1. 克隆项目
git clone <repo>
cd go-bisub

# 2. 配置环境
cp .env.example .env

# 3. 启动服务
make docker-up

# 4. 访问
open http://localhost:8080/admin
```

### 场景 2：修改密码

```bash
# 1. 编辑 .env
vim .env
# 修改 MYSQL_ROOT_PASSWORD 和 JWT_SECRET

# 2. 重启服务
make docker-down
make docker-up
```

### 场景 3：查看日志

```bash
# 实时日志
make docker-logs

# 或指定服务
docker-compose logs -f go-bisub
docker-compose logs -f mysql
docker-compose logs -f redis
```

### 场景 4：进入容器调试

```bash
# 进入应用容器
make docker-shell

# 进入 MySQL 容器
docker-compose exec mysql mysql -uroot -ppassword

# 进入 Redis 容器
docker-compose exec redis redis-cli
```

### 场景 5：数据持久化

```bash
# 数据存储在 Docker 卷中
docker volume ls | grep go-bisub

# 备份数据
docker-compose exec mysql mysqldump -uroot -ppassword go_sub > backup.sql

# 恢复数据
docker-compose exec -T mysql mysql -uroot -ppassword go_sub < backup.sql
```

## ⚠️ 常见问题

### Q1: 为什么访问不了 localhost:8080？

**A**: 检查端口是否被占用

```bash
# 查看端口占用
lsof -i :8080

# 如果被占用，修改 docker-compose.yml
ports:
  - "8081:8080"  # 改用 8081 端口
```

### Q2: 数据库连接失败？

**A**: 检查配置文件中的 host

```yaml
# ❌ 错误
database:
  host: 127.0.0.1  # Docker 容器内无法访问

# ✅ 正确
database:
  host: mysql      # 使用 Docker 服务名
```

### Q3: 修改配置后不生效？

**A**: 需要重启服务

```bash
# 重启应用容器
docker-compose restart go-bisub

# 或重启所有服务
make docker-restart
```

### Q4: 如何清理所有数据？

**A**: 删除容器和数据卷

```bash
# 停止并删除容器和数据卷
make docker-clean

# 或手动执行
docker-compose down -v
```

## 📊 配置检查清单

部署前检查：

- [ ] 已复制 `.env.example` 为 `.env`
- [ ] 已修改 `MYSQL_ROOT_PASSWORD`（生产环境）
- [ ] 已修改 `JWT_SECRET`（生产环境）
- [ ] `config.docker.yaml` 中的 host 使用服务名（mysql, redis）
- [ ] 端口 8080、3306、6379 未被占用
- [ ] Docker 和 Docker Compose 已安装

## 🔧 高级配置

### 使用外部数据库

如果你有独立的 MySQL 服务器：

```yaml
# config.docker.yaml
database:
  primary:
    host: db.example.com  # 外部数据库地址
    port: 3306
    username: your_user
    password: your_password
```

```yaml
# docker-compose.yml
services:
  go-bisub:
    # 移除 depends_on 中的 mysql
    depends_on:
      - redis  # 只依赖 Redis
  
  # 注释掉 mysql 服务
  # mysql:
  #   ...
```

### 使用外部 Redis

```yaml
# config.docker.yaml
redis:
  host: redis.example.com  # 外部 Redis 地址
  port: 6379
  password: your_password
```

## 📚 下一步

- 查看完整文档：[docs/DOCKER_DEPLOYMENT.md](docs/DOCKER_DEPLOYMENT.md)
- 配置详解：[docs/CONFIGURATION_GUIDE.md](docs/CONFIGURATION_GUIDE.md)
- 本地开发：[docs/LOCAL_DEVELOPMENT.md](docs/LOCAL_DEVELOPMENT.md)

## 🆘 获取帮助

遇到问题？

1. 查看日志：`make docker-logs`
2. 检查状态：`make docker-ps`
3. 查看文档：`docs/` 目录
4. 提交 Issue
