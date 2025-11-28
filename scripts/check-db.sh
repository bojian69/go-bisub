#!/bin/bash

# 数据库连接检查脚本

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "🔍 检查数据库连接..."
echo ""

# 从配置文件读取数据库信息（简化版本，实际应该解析 YAML）
DB_HOST=${DB_HOST:-"127.0.0.1"}
DB_PORT=${DB_PORT:-"3306"}
DB_USER=${DB_USER:-"root"}
DB_PASS=${DB_PASS:-""}
DB_NAME=${DB_NAME:-"go_sub"}

# 检查 MySQL 命令是否存在
if ! command -v mysql &> /dev/null; then
    echo -e "${YELLOW}⚠️  mysql 命令未找到${NC}"
    echo -e "${BLUE}ℹ️  请安装 MySQL 客户端或使用 Docker${NC}"
    exit 1
fi

# 测试连接
echo "测试连接: mysql -h $DB_HOST -P $DB_PORT -u $DB_USER"
echo ""

if [ -z "$DB_PASS" ]; then
    # 无密码连接
    if mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -e "SELECT 1" &> /dev/null; then
        echo -e "${GREEN}✓ MySQL 连接成功${NC}"
        
        # 检查数据库是否存在
        if mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -e "USE $DB_NAME" &> /dev/null; then
            echo -e "${GREEN}✓ 数据库 '$DB_NAME' 存在${NC}"
            
            # 检查表
            TABLE_COUNT=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -D "$DB_NAME" -N -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='$DB_NAME'")
            echo -e "${GREEN}✓ 数据库包含 $TABLE_COUNT 个表${NC}"
            
            if [ "$TABLE_COUNT" -eq 0 ]; then
                echo -e "${YELLOW}⚠️  数据库为空，需要初始化${NC}"
                echo -e "${BLUE}ℹ️  运行: mysql -h $DB_HOST -u $DB_USER < init.sql${NC}"
            fi
        else
            echo -e "${YELLOW}⚠️  数据库 '$DB_NAME' 不存在${NC}"
            echo -e "${BLUE}ℹ️  创建数据库: mysql -h $DB_HOST -u $DB_USER -e \"CREATE DATABASE $DB_NAME\"${NC}"
            echo -e "${BLUE}ℹ️  初始化表: mysql -h $DB_HOST -u $DB_USER < init.sql${NC}"
        fi
    else
        echo -e "${RED}❌ MySQL 连接失败${NC}"
        echo -e "${BLUE}ℹ️  请检查：${NC}"
        echo "  1. MySQL 是否已启动"
        echo "  2. 连接信息是否正确"
        echo "  3. 是否需要密码（设置 DB_PASS 环境变量）"
        exit 1
    fi
else
    # 有密码连接
    if mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" -e "SELECT 1" &> /dev/null; then
        echo -e "${GREEN}✓ MySQL 连接成功${NC}"
    else
        echo -e "${RED}❌ MySQL 连接失败${NC}"
        exit 1
    fi
fi

echo ""
echo "检查 Redis..."

# 检查 Redis
if command -v redis-cli &> /dev/null; then
    if redis-cli -h 127.0.0.1 ping &> /dev/null; then
        echo -e "${GREEN}✓ Redis 连接成功${NC}"
    else
        echo -e "${YELLOW}⚠️  Redis 连接失败${NC}"
        echo -e "${BLUE}ℹ️  启动 Redis: redis-server${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  redis-cli 命令未找到${NC}"
    echo -e "${BLUE}ℹ️  请安装 Redis 或使用 Docker${NC}"
fi

echo ""
echo "✅ 检查完成"
