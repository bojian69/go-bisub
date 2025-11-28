#!/bin/bash

# BI订阅服务 - 开发环境启动脚本

set -e

echo "🚀 启动 BI 订阅服务开发环境..."

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go 未安装，请先安装 Go 1.21+${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Go 版本: $(go version)${NC}"

# 获取 GOPATH
GOPATH=$(go env GOPATH)
GOBIN="$GOPATH/bin"

# 检查并安装 air
if [ ! -f "$GOBIN/air" ]; then
    echo -e "${YELLOW}⚠️  Air 未安装，正在安装...${NC}"
    go install github.com/air-verse/air@latest
    echo -e "${GREEN}✓ Air 安装完成${NC}"
else
    echo -e "${GREEN}✓ Air 已安装${NC}"
fi

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    if [ -f "config.yaml.example" ]; then
        echo -e "${YELLOW}⚠️  配置文件不存在，从示例复制...${NC}"
        cp config.yaml.example config.yaml
        echo -e "${GREEN}✓ 配置文件已创建，请根据需要修改 config.yaml${NC}"
    else
        echo -e "${RED}❌ 配置文件不存在，请创建 config.yaml${NC}"
        exit 1
    fi
fi

# 检查依赖服务
echo ""
echo "📦 检查依赖服务..."

# 检查是否使用 Docker
USE_DOCKER=false
if command -v docker &> /dev/null && docker info &> /dev/null; then
    USE_DOCKER=true
fi

# 检查 MySQL
echo -n "检查 MySQL... "
if command -v mysql &> /dev/null; then
    # 尝试连接本地 MySQL
    if mysql -h 127.0.0.1 -u root -e "SELECT 1" &> /dev/null 2>&1; then
        echo -e "${GREEN}✓ 本地 MySQL 运行中${NC}"
    else
        echo -e "${YELLOW}⚠️  本地 MySQL 未运行或需要密码${NC}"
        echo -e "${BLUE}ℹ️  请确保 MySQL 已启动并配置正确的连接信息${NC}"
    fi
elif [ "$USE_DOCKER" = true ]; then
    if docker ps | grep -q mysql; then
        echo -e "${GREEN}✓ Docker MySQL 运行中${NC}"
    else
        echo -e "${YELLOW}⚠️  MySQL 未运行，正在启动 Docker 容器...${NC}"
        docker-compose up -d mysql
        echo "⏳ 等待 MySQL 启动..."
        sleep 10
        echo -e "${GREEN}✓ MySQL 已启动${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  无法检测 MySQL 状态${NC}"
    echo -e "${BLUE}ℹ️  请确保 MySQL 已启动${NC}"
fi

# 检查 Redis
echo -n "检查 Redis... "
if command -v redis-cli &> /dev/null; then
    # 尝试连接本地 Redis
    if redis-cli -h 127.0.0.1 ping &> /dev/null; then
        echo -e "${GREEN}✓ 本地 Redis 运行中${NC}"
    else
        echo -e "${YELLOW}⚠️  本地 Redis 未运行${NC}"
        if [ "$USE_DOCKER" = true ]; then
            echo -e "${YELLOW}正在启动 Docker Redis...${NC}"
            docker-compose up -d redis
            sleep 3
            echo -e "${GREEN}✓ Redis 已启动${NC}"
        else
            echo -e "${BLUE}ℹ️  请手动启动 Redis: redis-server${NC}"
        fi
    fi
elif [ "$USE_DOCKER" = true ]; then
    if docker ps | grep -q redis; then
        echo -e "${GREEN}✓ Docker Redis 运行中${NC}"
    else
        echo -e "${YELLOW}⚠️  Redis 未运行，正在启动 Docker 容器...${NC}"
        docker-compose up -d redis
        sleep 3
        echo -e "${GREEN}✓ Redis 已启动${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  无法检测 Redis 状态${NC}"
    echo -e "${BLUE}ℹ️  请确保 Redis 已启动${NC}"
fi

# 下载依赖
echo ""
echo "📥 下载 Go 依赖..."
go mod download
echo -e "${GREEN}✓ 依赖下载完成${NC}"

# 启动服务
echo ""
echo "🎯 启动开发服务器（热重载）..."
echo ""
echo "访问地址："
echo "  - API: http://localhost:8080"
echo "  - 管理界面: http://localhost:8080/admin"
echo "  - 健康检查: http://localhost:8080/health"
echo ""
echo "按 Ctrl+C 停止服务"
echo ""

# 使用 air 启动
if [ -f "$GOBIN/air" ]; then
    "$GOBIN/air"
else
    # 降级到直接运行
    echo -e "${YELLOW}⚠️  使用直接运行模式（无热重载）${NC}"
    go run cmd/server/main.go
fi
