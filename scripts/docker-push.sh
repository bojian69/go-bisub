#!/bin/bash

# Docker 镜像推送脚本
# 用于构建并推送镜像到 Docker 仓库
# 
# 使用方法:
#   ./scripts/docker-push.sh <namespace> <version>
#   例如: ./scripts/docker-push.sh myusername v1.0.0

set -e

# 显示使用说明
show_usage() {
  echo "使用方法:"
  echo "  $0 <namespace> <version>"
  echo ""
  echo "参数说明:"
  echo "  namespace  - Docker Hub 用户名或组织名"
  echo "  version    - 镜像版本号（如 v1.0.0）"
  echo ""
  echo "示例:"
  echo "  $0 myusername v1.0.0"
  echo "  $0 mycompany latest"
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

# 完整镜像名称
FULL_IMAGE_NAME="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"

echo "=========================================="
echo "Docker 镜像构建和推送"
echo "=========================================="
echo "镜像名称: ${FULL_IMAGE_NAME}"
echo "=========================================="

# 1. 构建镜像
echo "步骤 1: 构建 Docker 镜像..."
docker build \
  --build-arg VERSION="${VERSION}" \
  --build-arg BUILD_TIME="$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
  --build-arg GIT_COMMIT="$(git rev-parse HEAD 2>/dev/null || echo 'unknown')" \
  -t "${FULL_IMAGE_NAME}" \
  -f Dockerfile \
  .

echo "✓ 镜像构建完成"

# 2. 标记 latest 版本
if [ "${VERSION}" != "latest" ]; then
  LATEST_IMAGE="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:latest"
  echo "步骤 2: 标记 latest 版本..."
  docker tag "${FULL_IMAGE_NAME}" "${LATEST_IMAGE}"
  echo "✓ 已标记: ${LATEST_IMAGE}"
fi

# 3. 登录 Docker 仓库（如果需要）
echo "步骤 3: 登录 Docker 仓库..."
if [ -n "${DOCKER_USERNAME}" ] && [ -n "${DOCKER_PASSWORD}" ]; then
  echo "${DOCKER_PASSWORD}" | docker login "${REGISTRY}" -u "${DOCKER_USERNAME}" --password-stdin
  echo "✓ 登录成功"
else
  echo "提示: 如果未登录，请手动执行: docker login ${REGISTRY}"
fi

# 4. 推送镜像
echo "步骤 4: 推送镜像到仓库..."
docker push "${FULL_IMAGE_NAME}"
echo "✓ 推送完成: ${FULL_IMAGE_NAME}"

if [ "${VERSION}" != "latest" ]; then
  docker push "${LATEST_IMAGE}"
  echo "✓ 推送完成: ${LATEST_IMAGE}"
fi

echo "=========================================="
echo "✓ 所有操作完成！"
echo "=========================================="
echo "镜像已推送到: ${FULL_IMAGE_NAME}"
echo ""
echo "在其他机器上拉取镜像："
echo "  docker pull ${FULL_IMAGE_NAME}"
echo ""
echo "运行容器："
echo "  docker run -d -p 8080:8080 --env-file .env ${FULL_IMAGE_NAME}"
echo "=========================================="
