# 云服务部署指南

本指南介绍如何使用云服务（云数据库和云 Redis）部署 go-bisub 应用。

## 前置条件

1. 已准备好云服务：
   - 云 MySQL 数据库（如阿里云 RDS、腾讯云 CDB 等）
   - 云 Redis 实例
2. 已安装 Docker 和 Docker Compose
3. 服务器可以访问云服务（网络连通性、安全组配置）

## 部署步骤

### 1. 准备云服务

#### MySQL 数据库初始化

在云 MySQL 中执行初始化脚本：

```bash
# 连接到云数据库
mysql -h your-cloud-mysql-host.com -u your-user -p

# 创建数据库
CREATE DATABASE IF NOT EXISTS go_sub DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 使用数据库
USE go_sub;

# 执行初始化脚本
source init.sql;
source init_operation_logs.sql;
```

或者使用云服务提供的数据库管理工具导入 SQL 文件。

#### Redis 配置

确保云 Redis 实例已创建并获取连接信息：
- 主机地址
- 端口（默认 6379）
- 密码
- 数据库编号（默认 0）

### 2. 配置环境变量

复制并编辑环境变量文件：

```bash
cp .env.example .env
```

编辑 `.env` 文件，填入实际的云服务配置：

```bash
# 应用端口
APP_PORT=8080

# 云数据库配置
DB_HOST=your-cloud-mysql-host.com
DB_PORT=3306
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=go_sub

# 云 Redis 配置
REDIS_HOST=your-cloud-redis-host.com
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0

# 数据源配置（执行订阅 SQL 的目标库）
DATASOURCE_DEFAULT_HOST=your-datasource-host.com
DATASOURCE_DEFAULT_PORT=3306
DATASOURCE_DEFAULT_USER=your-datasource-user
DATASOURCE_DEFAULT_PASS=your-datasource-password
DATASOURCE_DEFAULT_NAME=uhomes

# JWT 密钥（请修改为随机字符串）
JWT_SECRET=your-secure-random-jwt-secret-key

# 日志级别
LOG_LEVEL=info
```

### 3. 构建和启动应用

```bash
# 构建镜像
docker-compose build

# 启动应用
docker-compose up -d

# 查看日志
docker-compose logs -f go-bisub
```

### 4. 验证部署

```bash
# 检查容器状态
docker-compose ps

# 测试健康检查接口
curl http://localhost:8080/health

# 访问应用
curl http://localhost:8080/
```

## 网络配置

### 安全组配置

确保云服务的安全组允许应用服务器访问：

**MySQL 安全组规则：**
- 入站规则：允许应用服务器 IP 访问 3306 端口

**Redis 安全组规则：**
- 入站规则：允许应用服务器 IP 访问 6379 端口

### 白名单配置

在云服务控制台添加应用服务器 IP 到白名单：
- MySQL 白名单
- Redis 白名单

## 常见问题

### 1. 无法连接到云数据库

**检查项：**
- 网络连通性：`telnet your-cloud-mysql-host.com 3306`
- 安全组配置是否正确
- 白名单是否包含应用服务器 IP
- 数据库用户权限是否正确

**解决方法：**
```bash
# 测试数据库连接
docker exec -it go-bisub-app sh
wget -O- "http://localhost:8080/health"
```

### 2. Redis 连接失败

**检查项：**
- Redis 密码是否正确
- 端口是否正确（某些云服务使用非标准端口）
- 安全组和白名单配置

### 3. 应用启动失败

查看详细日志：
```bash
docker-compose logs go-bisub
```

检查配置文件：
```bash
docker exec -it go-bisub-app cat /app/config.yaml
```

## 生产环境建议

### 1. 使用私有镜像仓库

```bash
# 构建镜像
docker build -t your-registry.com/go-bisub:v1.0.0 .

# 推送到私有仓库
docker push your-registry.com/go-bisub:v1.0.0

# 在生产环境拉取
docker pull your-registry.com/go-bisub:v1.0.0
```

### 2. 使用 Docker Secrets 管理敏感信息

```yaml
# docker-compose.yml
services:
  go-bisub:
    secrets:
      - db_password
      - redis_password
      - jwt_secret

secrets:
  db_password:
    external: true
  redis_password:
    external: true
  jwt_secret:
    external: true
```

### 3. 配置日志收集

将日志输出到云日志服务：
```yaml
services:
  go-bisub:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### 4. 监控和告警

- 配置云监控服务监控容器状态
- 设置数据库连接数告警
- 设置 Redis 内存使用告警
- 配置应用健康检查告警

### 5. 备份策略

- 定期备份云数据库
- 配置自动备份策略
- 测试备份恢复流程

## 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建镜像
docker-compose build

# 重启服务（零停机更新）
docker-compose up -d

# 清理旧镜像
docker image prune -f
```

## 回滚

```bash
# 查看镜像历史
docker images go-bisub

# 使用旧版本镜像
docker tag go-bisub:old go-bisub:latest
docker-compose up -d
```

## 性能优化

### 1. 数据库连接池

在配置文件中调整连接池参数：
```yaml
database:
  primary:
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: 3600
```

### 2. Redis 连接池

```yaml
redis:
  pool_size: 10
  min_idle_conns: 5
```

### 3. 应用资源限制

```yaml
services:
  go-bisub:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

## 相关文档

- [Docker 快速开始](./DOCKER_QUICKSTART.md)
- [配置指南](./CONFIGURATION_GUIDE.md)
- [本地开发环境](./LOCAL_DEVELOPMENT.md)
