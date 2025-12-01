# Docker 部署指南

## 📦 快速开始

### 1. 准备环境

确保已安装：
- Docker 20.10+
- Docker Compose 2.0+

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑环境变量
vim .env
```

### 3. 启动服务

```bash
# 使用 docker-compose 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f go-bisub
```

### 4. 访问服务

- Web UI: http://localhost:8080/admin
- API: http://localhost:8080/v1
- 健康检查: http://localhost:8080/health

默认账号：
- 用户名: `admin`
- 密码: `admin123`

## 🔧 高级配置

### 自定义构建

```bash
# 使用构建脚本
./scripts/docker-build.sh v1.0.0

# 或手动构建
docker build \
  --build-arg VERSION=v1.0.0 \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  -t go-bisub:v1.0.0 \
  .
```

### 自定义配置文件

```bash
# 挂载自定义配置
docker-compose up -d \
  -v $(pwd)/config.prod.yaml:/app/config.yaml:ro
```

### 环境变量说明

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `VERSION` | 应用版本 | `dev` |
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码 | `password` |
| `MYSQL_USER` | MySQL 用户名 | `bisub` |
| `MYSQL_PASSWORD` | MySQL 密码 | `bisub123` |
| `REDIS_PASSWORD` | Redis 密码 | 空 |
| `JWT_SECRET` | JWT 密钥 | 需修改 |
| `GIN_MODE` | Gin 模式 | `release` |
| `LOG_LEVEL` | 日志级别 | `info` |

## 📊 服务管理

### 查看服务状态

```bash
# 查看所有服务
docker-compose ps

# 查看特定服务
docker-compose ps go-bisub
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f go-bisub

# 查看最近 100 行日志
docker-compose logs --tail=100 go-bisub
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart go-bisub
```

### 停止服务

```bash
# 停止所有服务
docker-compose stop

# 停止并删除容器
docker-compose down

# 停止并删除容器和数据卷
docker-compose down -v
```

## 🔍 健康检查

### 应用健康检查

```bash
# 检查应用健康状态
curl http://localhost:8080/health

# 预期响应
{
  "status": "ok",
  "timestamp": "2024-12-01T10:00:00Z"
}
```

### 数据库健康检查

```bash
# 进入 MySQL 容器
docker-compose exec mysql mysql -uroot -ppassword

# 检查数据库
SHOW DATABASES;
USE go_sub;
SHOW TABLES;
```

### Redis 健康检查

```bash
# 进入 Redis 容器
docker-compose exec redis redis-cli

# 检查连接
PING
# 预期响应: PONG
```

## 🚀 生产部署

### 1. 安全配置

```bash
# 修改默认密码
MYSQL_ROOT_PASSWORD=<strong-password>
MYSQL_PASSWORD=<strong-password>
REDIS_PASSWORD=<strong-password>
JWT_SECRET=<random-secret-key>
```

### 2. 性能优化

```yaml
# docker-compose.yml
services:
  go-bisub:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

### 3. 数据备份

```bash
# 备份 MySQL 数据
docker-compose exec mysql mysqldump -uroot -ppassword go_sub > backup.sql

# 备份 Redis 数据
docker-compose exec redis redis-cli SAVE
docker cp go-bisub-redis:/data/dump.rdb ./backup/
```

### 4. 日志管理

```yaml
# docker-compose.yml
services:
  go-bisub:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 🐛 故障排查

### 容器无法启动

```bash
# 查看容器日志
docker-compose logs go-bisub

# 查看容器详细信息
docker inspect go-bisub-app

# 检查端口占用
lsof -i :8080
```

### 数据库连接失败

```bash
# 检查 MySQL 是否启动
docker-compose ps mysql

# 检查 MySQL 日志
docker-compose logs mysql

# 测试数据库连接
docker-compose exec mysql mysql -uroot -ppassword -e "SELECT 1"
```

### Redis 连接失败

```bash
# 检查 Redis 是否启动
docker-compose ps redis

# 检查 Redis 日志
docker-compose logs redis

# 测试 Redis 连接
docker-compose exec redis redis-cli ping
```

## 📈 监控和指标

### Prometheus 集成

```yaml
# docker-compose.yml
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
```

### Grafana 集成

```yaml
# docker-compose.yml
services:
  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

## 🔄 更新和回滚

### 更新应用

```bash
# 拉取最新代码
git pull

# 重新构建镜像
docker-compose build

# 重启服务
docker-compose up -d
```

### 回滚版本

```bash
# 使用特定版本
docker-compose down
docker-compose up -d go-bisub:v1.0.0
```

## 📝 最佳实践

1. **使用环境变量**：不要在代码中硬编码敏感信息
2. **定期备份**：定期备份数据库和 Redis 数据
3. **监控日志**：使用日志聚合工具（如 ELK）
4. **资源限制**：设置合理的 CPU 和内存限制
5. **健康检查**：配置合适的健康检查参数
6. **网络隔离**：使用 Docker 网络隔离服务
7. **数据持久化**：使用 Docker 卷持久化数据
8. **安全更新**：定期更新基础镜像和依赖

## 🆘 获取帮助

- 查看日志：`docker-compose logs -f`
- 进入容器：`docker-compose exec go-bisub sh`
- 查看配置：`docker-compose config`
- 官方文档：https://docs.docker.com/
