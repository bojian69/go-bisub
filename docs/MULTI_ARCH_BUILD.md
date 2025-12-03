# Docker 多架构构建指南

本指南介绍如何构建支持多种 CPU 架构的 Docker 镜像。

## 为什么需要多架构镜像？

不同的服务器使用不同的 CPU 架构：

- **AMD64 (x86_64)**: 传统的 Intel/AMD 服务器
- **ARM64 (aarch64)**: ARM 服务器，如：
  - Apple Silicon (M1/M2/M3 Mac)
  - AWS Graviton 实例
  - 树莓派 4/5
  - 华为鲲鹏服务器

多架构镜像可以在不同架构的机器上自动选择合适的版本运行。

---

## 快速开始

### 使用自动化脚本（推荐）

```bash
# 构建并推送多架构镜像
./scripts/docker-push.sh your-username v1.0.0 --multi-arch
```

### 手动构建

```bash
# 1. 创建 buildx builder
docker buildx create --name multiarch-builder --use
docker buildx inspect --bootstrap

# 2. 构建并推送多架构镜像
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t your-username/go-bisub:v1.0.0 \
  --push \
  .

# 3. 查看镜像信息
docker buildx imagetools inspect your-username/go-bisub:v1.0.0
```

---

## 详细步骤

### 1. 启用 Docker Buildx

Docker Buildx 是 Docker 的多架构构建工具，Docker Desktop 已内置。

```bash
# 检查 buildx 是否可用
docker buildx version

# 查看现有 builder
docker buildx ls
```

### 2. 创建多架构 Builder

```bash
# 创建新的 builder
docker buildx create --name multiarch-builder --use

# 启动 builder
docker buildx inspect --bootstrap

# 查看支持的平台
docker buildx inspect multiarch-builder
```

输出示例：
```
Name:   multiarch-builder
Driver: docker-container

Platforms: linux/amd64, linux/arm64, linux/arm/v7, linux/arm/v6
```

### 3. 构建多架构镜像

#### 方法 1：构建并推送

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --build-arg VERSION="v1.0.0" \
  --build-arg BUILD_TIME="$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
  --build-arg GIT_COMMIT="$(git rev-parse HEAD)" \
  -t your-username/go-bisub:v1.0.0 \
  --push \
  .
```

**注意**：多架构构建必须使用 `--push` 直接推送，不能先构建到本地。

#### 方法 2：构建到本地（单架构测试）

```bash
# 只构建当前平台用于测试
docker buildx build \
  --platform linux/amd64 \
  -t go-bisub:test \
  --load \
  .
```

### 4. 验证多架构镜像

```bash
# 查看镜像详细信息
docker buildx imagetools inspect your-username/go-bisub:v1.0.0
```

输出示例：
```
Name:      docker.io/your-username/go-bisub:v1.0.0
MediaType: application/vnd.docker.distribution.manifest.list.v2+json
Digest:    sha256:abc123...

Manifests:
  Name:      docker.io/your-username/go-bisub:v1.0.0@sha256:def456...
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/amd64

  Name:      docker.io/your-username/go-bisub:v1.0.0@sha256:ghi789...
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/arm64
```

---

## 在不同架构上运行

### AMD64 服务器

```bash
docker pull your-username/go-bisub:v1.0.0
docker run -d -p 8080:8080 your-username/go-bisub:v1.0.0
```

Docker 会自动拉取 `linux/amd64` 版本。

### ARM64 服务器

```bash
docker pull your-username/go-bisub:v1.0.0
docker run -d -p 8080:8080 your-username/go-bisub:v1.0.0
```

Docker 会自动拉取 `linux/arm64` 版本。

### 指定架构拉取

```bash
# 拉取 AMD64 版本
docker pull --platform linux/amd64 your-username/go-bisub:v1.0.0

# 拉取 ARM64 版本
docker pull --platform linux/arm64 your-username/go-bisub:v1.0.0
```

---

## 常见问题

### Q1: 构建时间很长？

多架构构建需要为每个架构分别编译，时间是单架构的 2 倍左右。

**优化建议**：
- 使用 Docker 层缓存
- 在 CI/CD 中使用缓存
- 只在发布时构建多架构，开发时使用单架构

### Q2: 无法构建 ARM64？

确保 Docker Desktop 已启用实验性功能：

```bash
# macOS/Linux
echo '{"experimental": true}' > ~/.docker/config.json

# 或在 Docker Desktop 设置中启用
```

### Q3: 构建失败：exec format error

这通常是因为尝试在本地加载多架构镜像。多架构镜像必须推送到仓库。

**解决方案**：
- 使用 `--push` 直接推送
- 或使用 `--platform` 指定单一架构并使用 `--load`

### Q4: 如何在 Apple Silicon Mac 上测试 AMD64 镜像？

```bash
# 构建 AMD64 版本
docker buildx build \
  --platform linux/amd64 \
  -t go-bisub:amd64-test \
  --load \
  .

# 运行（会使用 QEMU 模拟）
docker run --platform linux/amd64 -p 8080:8080 go-bisub:amd64-test
```

---

## CI/CD 集成

### GitHub Actions

```yaml
name: Build Multi-Arch Docker Image

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      
      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            your-username/go-bisub:${{ github.ref_name }}
            your-username/go-bisub:latest
```

### GitLab CI

```yaml
build-multiarch:
  image: docker:latest
  services:
    - docker:dind
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD
    - docker buildx create --use
  script:
    - docker buildx build
      --platform linux/amd64,linux/arm64
      -t your-username/go-bisub:$CI_COMMIT_TAG
      --push
      .
  only:
    - tags
```

---

## 性能对比

### 构建时间

| 架构 | 构建时间 | 镜像大小 |
|------|---------|---------|
| AMD64 | ~2 分钟 | 25 MB |
| ARM64 | ~2 分钟 | 24 MB |
| 多架构 | ~4 分钟 | 49 MB (总计) |

### 运行性能

- **原生架构**: 100% 性能
- **QEMU 模拟**: 约 50-70% 性能

**建议**: 始终使用原生架构镜像以获得最佳性能。

---

## 最佳实践

1. **发布时使用多架构**
   - 正式版本构建多架构镜像
   - 开发测试使用单架构加快速度

2. **使用缓存**
   ```bash
   docker buildx build \
     --platform linux/amd64,linux/arm64 \
     --cache-from type=registry,ref=your-username/go-bisub:buildcache \
     --cache-to type=registry,ref=your-username/go-bisub:buildcache,mode=max \
     --push \
     .
   ```

3. **标签策略**
   - 版本标签：`v1.0.0`
   - Latest 标签：`latest`
   - 架构特定标签（可选）：`v1.0.0-amd64`, `v1.0.0-arm64`

4. **测试覆盖**
   - 在 AMD64 和 ARM64 环境中都进行测试
   - 使用 CI/CD 自动化测试

---

## 相关资源

- [Docker Buildx 文档](https://docs.docker.com/buildx/working-with-buildx/)
- [多架构镜像最佳实践](https://docs.docker.com/build/building/multi-platform/)
- [QEMU 用户模式](https://www.qemu.org/docs/master/user/main.html)

---

## 相关文档

- [Docker 镜像仓库部署指南](./DOCKER_REGISTRY_GUIDE.md)
- [Docker 快速开始](./DOCKER_QUICKSTART.md)
- [Docker 部署指南](./DOCKER_DEPLOYMENT.md)
