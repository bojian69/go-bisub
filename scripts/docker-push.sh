#!/bin/bash

# Docker 镜像推送脚本（支持多架构）
# 用于构建并推送镜像到 Docker 仓库
# 支持 linux/amd64 和 linux/arm64 架构
# 
# 使用方法:
#   ./scripts/docker-push.sh <namespace> <version> [--multi-arch]
#   例如: ./scripts/docker-push.sh myusername v1.0.0
#        ./scripts/docker-push.sh myusername v1.0.0 --multi-arch

set -e

# 显示使用说明
show_usage() {
  echo "使用方法:"
  echo "  $0 <namespace> <version> [--multi-arch]"
  echo ""
  echo "参数说明:"
  echo "  namespace    - Docker Hub 用户名或组织名"
  echo "  version      - 镜像版本号（如 v1.0.0）"
  echo "  --multi-arch - 构建多架构镜像（amd64 + arm64）"
  echo ""
  echo "示例:"
  echo "  $0 myusername v1.0.0              # 单架构（当前平台）"
  echo "  $0 myusername v1.0.0 --multi-arch # 多架构（amd64 + arm64）"
  exit 1
}

# 检查参数
if [ $# -lt 2 ]; then
  echo "错误: 缺少必需参数"
  echo ""
  show_usage
fi

# 配置变量
IMAGE_NAME="go-bisub"
REGISTRY="${DOCKER_REGISTRY:-docker.io}"  # 默认 Docker Hub，可通过环境变量修改
NAMESPACE="$1"  # 从命令行参数获取
VERSION="$2"    # 从命令行参数获取
MULTI_ARCH=false

# 检查是否启用多架构
if [ "$3" = "--multi-arch" ]; then
  MULTI_ARCH=true
fi

# 完整镜像名称
FULL_IMAGE_NAME="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"
LATEST_IMAGE="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:latest"

# 构建参数
BUILD_ARGS="--build-arg VERSION=${VERSION} \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo 'unknown')"

echo "=========================================="
echo "Docker 镜像构建和推送"
echo "=========================================="
echo "镜像名称: ${FULL_IMAGE_NAME}"
if [ "$MULTI_ARCH" = true ]; then
  echo "架构支持: linux/amd64, linux/arm64"
else
  echo "架构支持: 当前平台"
fi
echo "=========================================="

# 登录 Docker 仓库
echo "步骤 1: 登录 Docker 仓库..."
if [ -n "${DOCKER_USERNAME}" ] && [ -n "${DOCKER_PASSWORD}" ]; then
  echo "${DOCKER_PASSWORD}" | docker login "${REGISTRY}" -u "${DOCKER_USERNAME}" --password-stdin
  echo "✓ 登录成功"
else
  echo "提示: 如果未登录，请手动执行: docker login ${REGISTRY}"
fi

if [ "$MULTI_ARCH" = true ]; then
  # 多架构构建
  echo ""
  echo "步骤 2: 创建/使用 buildx builder..."
  
  # 检查 builder 是否存在
  if ! docker buildx inspect multiarch-builder > /dev/null 2>&1; then
    echo "创建新的 builder: multiarch-builder"
    docker buildx create --name multiarch-builder --use
  else
    echo "使用现有 builder: multiarch-builder"
    docker buildx use multiarch-builder
  fi
  
  # 启动 builder
  docker buildx inspect --bootstrap
  echo "✓ Builder 准备完成"
  
  echo ""
  echo "步骤 3: 构建并推送多架构镜像..."
  echo "这可能需要几分钟时间，请耐心等待..."
  
  # 构建并推送版本镜像
  docker buildx build \
    --platform linux/amd64,linux/arm64 \
    ${BUILD_ARGS} \
    -t "${FULL_IMAGE_NAME}" \
    -f Dockerfile \
    --push \
    .
  
  echo "✓ 推送完成: ${FULL_IMAGE_NAME}"
  
  # 构建并推送 latest 镜像
  if [ "${VERSION}" != "latest" ]; then
    echo ""
    echo "步骤 4: 推送 latest 标签..."
    docker buildx build \
      --platform linux/amd64,linux/arm64 \
      ${BUILD_ARGS} \
      -t "${LATEST_IMAGE}" \
      -f Dockerfile \
      --push \
      .
    echo "✓ 推送完成: ${LATEST_IMAGE}"
  fi
  
else
  # 单架构构建
  echo ""
  echo "步骤 2: 构建 Docker 镜像..."
  docker build \
    ${BUILD_ARGS} \
    -t "${FULL_IMAGE_NAME}" \
    -f Dockerfile \
    .
  
  echo "✓ 镜像构建完成"
  
  # 标记 latest 版本
  if [ "${VERSION}" != "latest" ]; then
    echo ""
    echo "步骤 3: 标记 latest 版本..."
    docker tag "${FULL_IMAGE_NAME}" "${LATEST_IMAGE}"
    echo "✓ 已标记: ${LATEST_IMAGE}"
  fi
  
  # 推送镜像
  echo ""
  echo "步骤 4: 推送镜像到仓库..."
  docker push "${FULL_IMAGE_NAME}"
  echo "✓ 推送完成: ${FULL_IMAGE_NAME}"
  
  if [ "${VERSION}" != "latest" ]; then
    docker push "${LATEST_IMAGE}"
    echo "✓ 推送完成: ${LATEST_IMAGE}"
  fi
fi

echo ""
echo "=========================================="
echo "✓ 所有操作完成！"
echo "=========================================="
echo "镜像已推送到: ${FULL_IMAGE_NAME}"
if [ "$MULTI_ARCH" = true ]; then
  echo "支持架构: linux/amd64, linux/arm64"
fi
echo ""
echo "在其他机器上拉取镜像："
echo "  docker pull ${FULL_IMAGE_NAME}"
echo ""
echo "运行容器："
echo "  docker run -d -p 8080:8080 --env-file .env ${FULL_IMAGE_NAME}"
echo ""
echo "查看镜像信息："
echo "  docker buildx imagetools inspect ${FULL_IMAGE_NAME}"
echo "=========================================="
