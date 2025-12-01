#!/bin/bash

# Docker 镜像推送脚本
# 用于将本地构建的镜像推送到远程仓库

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认值
APP_NAME="go-bisub"
VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "latest")}

echo -e "${GREEN}=== Docker 镜像推送工具 ===${NC}"
echo ""

# 检查本地镜像是否存在
if ! docker images | grep -q "^${APP_NAME}"; then
    echo -e "${RED}错误: 本地镜像 ${APP_NAME} 不存在${NC}"
    echo "请先运行: make docker-build"
    exit 1
fi

echo "本地镜像: ${APP_NAME}:latest"
echo "版本标签: ${VERSION}"
echo ""

# 选择镜像仓库类型
echo "请选择镜像仓库类型:"
echo "1) Docker Hub (docker.io)"
echo "2) 阿里云容器镜像服务 (registry.cn-hangzhou.aliyuncs.com)"
echo "3) 腾讯云容器镜像服务 (ccr.ccs.tencentyun.com)"
echo "4) 自定义仓库"
read -p "请输入选项 [1-4]: " registry_choice

case $registry_choice in
    1)
        read -p "请输入 Docker Hub 用户名: " username
        REGISTRY="docker.io"
        FULL_REGISTRY="${username}"
        ;;
    2)
        read -p "请输入阿里云命名空间: " namespace
        REGISTRY="registry.cn-hangzhou.aliyuncs.com"
        FULL_REGISTRY="${REGISTRY}/${namespace}"
        ;;
    3)
        read -p "请输入腾讯云命名空间: " namespace
        REGISTRY="ccr.ccs.tencentyun.com"
        FULL_REGISTRY="${REGISTRY}/${namespace}"
        ;;
    4)
        read -p "请输入完整的仓库地址 (例如: registry.example.com/namespace): " custom_registry
        FULL_REGISTRY="${custom_registry}"
        ;;
    *)
        echo -e "${RED}无效的选项${NC}"
        exit 1
        ;;
esac

# 确认镜像名称
read -p "镜像名称 [${APP_NAME}]: " image_name
image_name=${image_name:-${APP_NAME}}

# 构建完整的镜像标签
IMAGE_TAG_VERSION="${FULL_REGISTRY}/${image_name}:${VERSION}"
IMAGE_TAG_LATEST="${FULL_REGISTRY}/${image_name}:latest"

echo ""
echo -e "${YELLOW}准备推送以下镜像:${NC}"
echo "  - ${IMAGE_TAG_VERSION}"
echo "  - ${IMAGE_TAG_LATEST}"
echo ""

read -p "是否继续? [y/N]: " confirm
if [[ ! $confirm =~ ^[Yy]$ ]]; then
    echo "已取消"
    exit 0
fi

# 登录到镜像仓库
echo ""
echo -e "${GREEN}登录到镜像仓库...${NC}"
docker login ${REGISTRY}

if [ $? -ne 0 ]; then
    echo -e "${RED}登录失败${NC}"
    exit 1
fi

# 打标签
echo ""
echo -e "${GREEN}为镜像打标签...${NC}"
docker tag ${APP_NAME}:latest ${IMAGE_TAG_VERSION}
docker tag ${APP_NAME}:latest ${IMAGE_TAG_LATEST}

# 推送镜像
echo ""
echo -e "${GREEN}推送镜像 ${IMAGE_TAG_VERSION}...${NC}"
docker push ${IMAGE_TAG_VERSION}

echo ""
echo -e "${GREEN}推送镜像 ${IMAGE_TAG_LATEST}...${NC}"
docker push ${IMAGE_TAG_LATEST}

# 完成
echo ""
echo -e "${GREEN}=== 推送完成 ===${NC}"
echo ""
echo "镜像已成功推送到远程仓库:"
echo "  ${IMAGE_TAG_VERSION}"
echo "  ${IMAGE_TAG_LATEST}"
echo ""
echo "在其他服务器上拉取镜像:"
echo "  docker pull ${IMAGE_TAG_VERSION}"
echo ""
echo "或在 docker-compose.yml 中使用:"
echo "  image: ${IMAGE_TAG_VERSION}"
