# Docker 镜像仓库部署指南

本指南介绍如何将 go-bisub 应用打包成 Docker 镜像，推送到镜像仓库，并在其他机器上部署。

## 目录

- [准备工作](#准备工作)
- [构建和推送镜像](#构建和推送镜像)
- [在其他机器上部署](#在其他机器上部署)
- [使用不同的镜像仓库](#使用不同的镜像仓库)

---

## 准备工作

### 1. 选择镜像仓库

你可以使用以下任一镜像仓库：

- **Docker Hub**（公共/私有）：https://hub.docker.com
- **阿里云容器镜像服务**：https://cr.console.aliyun.com
- **腾讯云容器镜像服务**：https://console.cloud.tencent.com/tcr
- **Harbor**（私有部署）
- **GitHub Container Registry**：https://ghcr.io

### 2. 创建仓库账号

以 Docker Hub 为例：
1. 访问 https://hub.docker.com 注册账号
2. 创建一个新的仓库（Repository）
3. 记录你的用户名和仓库名

### 3. 本地登录

```bash
# Docker Hub
docker login

# 阿里云（示例）
docker login --username=your-username registry.cn-hangzhou.aliyuncs.com

# 腾讯云（示例）
docker login --username=your-username ccr.ccs.tencentyun.com
```

---

## 构建和推送镜像

### 方法 1：使用自动化脚本（推荐）

#### 1. 执行推送脚本

```bash
# 赋予执行权限
chmod +x scripts/docker-push.sh

# 执行推送（传入用户名和版本号）
./scripts/docker-push.sh your-username v1.0.0

# 或者推送为 latest 版本
./scripts/docker-push.sh your-username latest
```

**可选：配置自动登录**

如果需要脚本自动登录，可以设置环境变量：

```bash
export DOCKER_USERNAME="your-username"
export DOCKER_PASSWORD="your-password"
./scripts/docker-push.sh your-username v1.0.0
```

### 方法 2：手动构建和推送

#### 1. 构建镜像

```bash
# 基本构建
docker build -t your-username/go-bisub:v1.0.0 .

# 带版本信息的构建
docker build \
  --build-arg VERSION="v1.0.0" \
  --build-arg BUILD_TIME="$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
  --build-arg GIT_COMMIT="$(git rev-parse HEAD)" \
  -t your-username/go-bisub:v1.0.0 \
  .
```

#### 2. 标记镜像

```bash
# 标记为 latest
docker tag your-username/go-bisub:v1.0.0 your-username/go-bisub:latest

# 如果使用其他仓库，需要添加仓库地址
docker tag your-username/go-bisub:v1.0.0 registry.cn-hangzhou.aliyuncs.com/your-namespace/go-bisub:v1.0.0
```

#### 3. 推送镜像

```bash
# 推送到 Docker Hub
docker push your-username/go-bisub:v1.0.0
docker push your-username/go-bisub:latest

# 推送到阿里云
docker push registry.cn-hangzhou.aliyuncs.com/your-namespace/go-bisub:v1.0.0
```

---

## 在其他机器上部署

### 方法 1：使用自动化脚本（推荐）

#### 1. 准备部署机器

```bash
# 在目标机器上创建工作目录
mkdir -p ~/go-bisub
cd ~/go-bisub

# 复制必要文件
# - scripts/docker-deploy.sh
# - .env.example
```

#### 2. 配置环境

```bash
# 复制并编辑环境变量
cp .env.example .env
vim .env
```

#### 3. 执行部署

```bash
# 赋予执行权限
chmod +x scripts/docker-deploy.sh

# 执行部署（传入用户名和版本号）
./scripts/docker-deploy.sh your-username v1.0.0

# 或者部署 latest 版本
./scripts/docker-deploy.sh your-username latest
```

### 方法 2：手动部署

#### 1. 拉取镜像

```bash
# 从 Docker Hub 拉取
docker pull your-username/go-bisub:v1.0.0

# 从阿里云拉取
docker pull registry.cn-hangzhou.aliyuncs.com/your-namespace/go-bisub:v1.0.0
```

#### 2. 准备配置文件

```bash
# 创建 .env 文件
cat > .env << 'EOF'
# 数据库配置
DB_HOST=your-db-host
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=go_sub

# Redis 配置
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0

# 其他配置...
EOF
```

#### 3. 启动容器

```bash
# 基本启动
docker run -d \
  --name go-bisub-app \
  --restart unless-stopped \
  -p 8080:8080 \
  --env-file .env \
  -v $(pwd)/logs:/app/logs \
  your-username/go-bisub:v1.0.0
```

#### 4. 验证部署

```bash
# 查看容器状态
docker ps

# 查看日志
docker logs -f go-bisub-app

# 测试服务
curl http://localhost:8080/health
```

### 方法 3：使用 docker-compose

#### 1. 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  go-bisub:
    image: your-username/go-bisub:v1.0.0
    container_name: go-bisub-app
    ports:
      - "8080:8080"
    env_file:
      - .env
    volumes:
      - ./logs:/app/logs
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s
```

#### 2. 启动服务

```bash
# 拉取并启动
docker-compose pull
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

---

## 使用不同的镜像仓库

### Docker Hub

```bash
# 构建
docker build -t your-username/go-bisub:v1.0.0 .

# 推送
docker push your-username/go-bisub:v1.0.0

# 拉取
docker pull your-username/go-bisub:v1.0.0
```

### 阿里云容器镜像服务

```bash
# 登录
docker login --username=your-username registry.cn-hangzhou.aliyuncs.com

# 构建和标记
docker build -t go-bisub:v1.0.0 .
docker tag go-bisub:v1.0.0 registry.cn-hangzhou.aliyuncs.com/your-namespace/go-bisub:v1.0.0

# 推送
docker push registry.cn-hangzhou.aliyuncs.com/your-namespace/go-bisub:v1.0.0

# 拉取
docker pull registry.cn-hangzhou.aliyuncs.com/your-namespace/go-bisub:v1.0.0
```

### 腾讯云容器镜像服务

```bash
# 登录
docker login --username=your-username ccr.ccs.tencentyun.com

# 构建和标记
docker build -t go-bisub:v1.0.0 .
docker tag go-bisub:v1.0.0 ccr.ccs.tencentyun.com/your-namespace/go-bisub:v1.0.0

# 推送
docker push ccr.ccs.tencentyun.com/your-namespace/go-bisub:v1.0.0

# 拉取
docker pull ccr.ccs.tencentyun.com/your-namespace/go-bisub:v1.0.0
```

### GitHub Container Registry

```bash
# 登录（使用 Personal Access Token）
echo $GITHUB_TOKEN | docker login ghcr.io -u your-username --password-stdin

# 构建和标记
docker build -t go-bisub:v1.0.0 .
docker tag go-bisub:v1.0.0 ghcr.io/your-username/go-bisub:v1.0.0

# 推送
docker push ghcr.io/your-username/go-bisub:v1.0.0

# 拉取
docker pull ghcr.io/your-username/go-bisub:v1.0.0
```

---

## 常用命令

### 镜像管理

```bash
# 查看本地镜像
docker images

# 删除镜像
docker rmi your-username/go-bisub:v1.0.0

# 清理未使用的镜像
docker image prune -a
```

### 容器管理

```bash
# 查看运行中的容器
docker ps

# 查看所有容器
docker ps -a

# 停止容器
docker stop go-bisub-app

# 启动容器
docker start go-bisub-app

# 重启容器
docker restart go-bisub-app

# 删除容器
docker rm go-bisub-app

# 查看容器日志
docker logs -f go-bisub-app

# 进入容器
docker exec -it go-bisub-app sh
```

### 更新部署

```bash
# 拉取最新镜像
docker pull your-username/go-bisub:latest

# 停止并删除旧容器
docker stop go-bisub-app
docker rm go-bisub-app

# 启动新容器
docker run -d \
  --name go-bisub-app \
  --restart unless-stopped \
  -p 8080:8080 \
  --env-file .env \
  -v $(pwd)/logs:/app/logs \
  your-username/go-bisub:latest
```

---

## 故障排查

### 1. 镜像推送失败

```bash
# 检查登录状态
docker login

# 检查镜像标记是否正确
docker images | grep go-bisub

# 检查网络连接
ping hub.docker.com
```

### 2. 容器启动失败

```bash
# 查看容器日志
docker logs go-bisub-app

# 检查环境变量
docker exec go-bisub-app env

# 检查配置文件
docker exec go-bisub-app cat /app/config.yaml
```

### 3. 无法访问服务

```bash
# 检查容器是否运行
docker ps | grep go-bisub-app

# 检查端口映射
docker port go-bisub-app

# 检查防火墙
sudo firewall-cmd --list-ports

# 测试容器内服务
docker exec go-bisub-app wget -O- http://localhost:8080/health
```

---

## 最佳实践

1. **使用版本标签**：不要只使用 `latest`，为每个版本打上明确的标签
2. **环境变量管理**：敏感信息使用 `.env` 文件，不要硬编码
3. **日志管理**：使用 volume 挂载日志目录，便于查看和备份
4. **健康检查**：配置 healthcheck 确保服务正常运行
5. **资源限制**：在生产环境中设置 CPU 和内存限制
6. **自动重启**：使用 `--restart unless-stopped` 确保服务自动恢复
7. **镜像优化**：定期清理未使用的镜像和容器

---

## 相关文档

- [Docker 快速开始](./DOCKER_QUICKSTART.md)
- [Docker 部署指南](./DOCKER_DEPLOYMENT.md)
- [配置指南](./CONFIGURATION_GUIDE.md)
