# Docker 云服务快速开始

使用云数据库和云 Redis 部署 go-bisub 应用的快速指南。

## 前提条件

- 已有云 MySQL 数据库（阿里云 RDS、腾讯云 CDB 等）
- 已有云 Redis 实例
- 已安装 Docker 和 Docker Compose

## 5 分钟快速部署

### 1. 初始化云数据库

在云 MySQL 中执行初始化脚本：

```bash
# 方式 A: 使用 mysql 命令行
mysql -h your-cloud-mysql.com -u your-user -p < init.sql
mysql -h your-cloud-mysql.com -u your-user -p < init_operation_logs.sql

# 方式 B: 使用云服务控制台的 SQL 执行工具
# 复制 init.sql 和 init_operation_logs.sql 的内容到控制台执行
```

### 2. 配置环境变量

```bash
# 复制配置模板
cp .env.example .env

# 编辑 .env 文件
vim .env
```

最少需要配置以下内容：

```bash
# 云数据库配置
DB_HOST=your-cloud-mysql.com
DB_PORT=3306
DB_USER=your-user
DB_PASSWORD=your-password
DB_NAME=go_sub

# 云 Redis 配置
REDIS_HOST=your-cloud-redis.com
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password

# JWT 密钥（请修改）
JWT_SECRET=your-random-secret-key

# 数据源配置（执行订阅 SQL 的目标库）
DATASOURCE_DEFAULT_HOST=your-datasource.com
DATASOURCE_DEFAULT_USER=your-datasource-user
DATASOURCE_DEFAULT_PASS=your-datasource-password
```

### 3. 构建和启动

```bash
# 构建镜像
make docker-build

# 启动应用
make docker-up

# 查看日志
make docker-logs
```

### 4. 验证部署

```bash
# 检查容器状态
docker ps

# 测试健康检查
curl http://localhost:8080/health

# 访问应用
open http://localhost:8080
```

## 推送到远程仓库

如果需要在其他服务器部署，可以推送镜像到远程仓库：

```bash
# 使用交互式脚本
./scripts/docker-push.sh

# 或使用 Makefile
make docker-push
```

## 在生产服务器部署

### 方式 A: 使用远程镜像

```bash
# 1. 在生产服务器上创建配置
vim .env

# 2. 创建 docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'
services:
  go-bisub:
    image: your-registry/go-bisub:latest
    container_name: go-bisub-app
    ports:
      - "8080:8080"
    env_file:
      - .env
    volumes:
      - ./logs:/app/logs
    restart: unless-stopped
EOF

# 3. 启动服务
docker-compose up -d
```

### 方式 B: 本地构建

```bash
# 1. 克隆代码
git clone <your-repo>
cd go-bisub

# 2. 配置环境变量
cp .env.example .env
vim .env

# 3. 构建并启动
make docker-build
make docker-up
```

## 常用命令

```bash
# 查看日志
make docker-logs

# 重启服务
make docker-restart

# 停止服务
make docker-down

# 进入容器
make docker-shell

# 查看状态
docker ps
```

## 网络配置检查清单

- [ ] 云数据库安全组允许应用服务器 IP
- [ ] 云 Redis 安全组允许应用服务器 IP
- [ ] 云数据库白名单包含应用服务器 IP
- [ ] 云 Redis 白名单包含应用服务器 IP
- [ ] 应用服务器可以访问云服务（telnet 测试）

## 故障排查

### 无法连接数据库

```bash
# 1. 检查网络连通性
telnet your-cloud-mysql.com 3306

# 2. 查看应用日志
docker logs go-bisub-app

# 3. 进入容器测试
docker exec -it go-bisub-app sh
wget -O- http://localhost:8080/health
```

### 应用无法启动

```bash
# 查看详细日志
docker logs go-bisub-app

# 检查环境变量
docker exec go-bisub-app env | grep DB_

# 重新创建容器
docker-compose down
docker-compose up -d
```

## 环境变量说明

| 变量名 | 说明 | 示例 |
|--------|------|------|
| DB_HOST | 云数据库地址 | rm-xxx.mysql.rds.aliyuncs.com |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户 | root |
| DB_PASSWORD | 数据库密码 | your-password |
| DB_NAME | 数据库名称 | go_sub |
| REDIS_HOST | Redis 地址 | r-xxx.redis.rds.aliyuncs.com |
| REDIS_PORT | Redis 端口 | 6379 |
| REDIS_PASSWORD | Redis 密码 | your-redis-password |
| JWT_SECRET | JWT 密钥 | 随机字符串 |
| APP_PORT | 应用端口 | 8080 |

## 下一步

- 配置 HTTPS（使用 Nginx 反向代理）
- 设置监控和告警
- 配置日志收集
- 实现自动备份

## 相关文档

- [完整云服务部署指南](./CLOUD_DEPLOYMENT.md)
- [Docker 配置更新](./DOCKER_UPDATES.md)
- [本地开发指南](./LOCAL_DEVELOPMENT.md)
