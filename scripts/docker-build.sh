#!/bin/bash

# Docker 构建脚本
# 用法: ./scripts/docker-build.sh [version]

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取版本号
VERSION=${1:-dev}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${GREEN}=== Docker 构建开始 ===${NC}"
echo -e "${YELLOW}版本: ${VERSION}${NC}"
echo -e "${YELLOW}构建时间: ${BUILD_TIME}${NC}"
echo -e "${YELLOW}Git Commit: ${GIT_COMMIT}${NC}"

# 构建镜像
echo -e "${GREEN}正在构建 Docker 镜像...${NC}"
docker build \
  --build-arg VERSION="${VERSION}" \
  --build-arg BUILD_TIME="${BUILD_TIME}" \
  --build-arg GIT_COMMIT="${GIT_COMMIT}" \
  -t go-bisub:${VERSION} \
  -t go-bisub:latest \
  .

echo -e "${GREEN}=== 构建完成 ===${NC}"
echo -e "${YELLOW}镜像标签:${NC}"
echo "  - go-bisub:${VERSION}"
echo "  - go-bisub:latest"

# 显示镜像信息
echo -e "${GREEN}=== 镜像信息 ===${NC}"
docker images | grep go-bisub | head -2

echo -e "${GREEN}=== 使用方法 ===${NC}"
echo "启动容器:"
echo "  docker-compose up -d"
echo ""
echo "查看日志:"
echo "  docker-compose logs -f go-bisub"
echo ""
echo "停止容器:"
echo "  docker-compose down"
