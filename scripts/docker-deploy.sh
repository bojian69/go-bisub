#!/bin/bash

# Docker 部署脚本
# 用于在目标机器上拉取并运行镜像
#
# 使用方法:
#   ./scripts/docker-deploy.sh <namespace> <version>
#   例如: ./scripts/docker-deploy.sh myusername v1.0.0

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
REGISTRY="${DOCKER_REGISTRY:-docker.io}"
NAMESPACE="$1"  # 从命令行参数获取
VERSION="$2"    # 从命令行参数获取
CONTAINER_NAME="go-bisub-app"

# 完整镜像名称
FULL_IMAGE_NAME="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"

echo "=========================================="
echo "Docker 容器部署"
echo "=========================================="
echo "镜像: ${FULL_IMAGE_NAME}"
echo "容器名称: ${CONTAINER_NAME}"
echo "=========================================="

# 1. 停止并删除旧容器（如果存在）
echo "步骤 1: 清理旧容器..."
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
  echo "发现旧容器，正在停止..."
  docker stop "${CONTAINER_NAME}" || true
  echo "正在删除旧容器..."
  docker rm "${CONTAINER_NAME}" || true
  echo "✓ 旧容器已清理"
else
  echo "✓ 没有发现旧容器"
fi

# 2. 拉取最新镜像
echo "步骤 2: 拉取最新镜像..."
docker pull "${FULL_IMAGE_NAME}"
echo "✓ 镜像拉取完成"

# 3. 检查 .env 文件
if [ ! -f ".env" ]; then
  echo "警告: 未找到 .env 文件"
  echo "请创建 .env 文件或从 .env.example 复制："
  echo "  cp .env.example .env"
  echo "然后编辑 .env 文件配置环境变量"
  exit 1
fi

# 4. 启动新容器
echo "步骤 3: 启动新容器..."
docker run -d \
  --name "${CONTAINER_NAME}" \
  --restart unless-stopped \
  -p 8081:8080 \
  --env-file .env \
  -v "$(pwd)/logs:/app/logs" \
  "${FULL_IMAGE_NAME}"

echo "✓ 容器启动成功"

# 5. 等待健康检查
echo "步骤 4: 等待服务启动..."
sleep 5

# 6. 检查容器状态
if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
  echo "✓ 容器运行正常"
  echo ""
  echo "=========================================="
  echo "部署成功！"
  echo "=========================================="
  echo "容器名称: ${CONTAINER_NAME}"
  echo "访问地址: http://localhost:8081"
  echo ""
  echo "查看日志："
  echo "  docker logs -f ${CONTAINER_NAME}"
  echo ""
  echo "停止容器："
  echo "  docker stop ${CONTAINER_NAME}"
  echo ""
  echo "重启容器："
  echo "  docker restart ${CONTAINER_NAME}"
  echo "=========================================="
else
  echo "✗ 容器启动失败"
  echo "查看日志："
  echo "  docker logs ${CONTAINER_NAME}"
  exit 1
fi
