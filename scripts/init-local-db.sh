#!/bin/bash

# 本地数据库初始化脚本

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "🗄️  初始化本地数据库..."
echo ""

# 数据库配置
DB_HOST=${DB_HOST:-"127.0.0.1"}
DB_PORT=${DB_PORT:-"3306"}
DB_USER=${DB_USER:-"root"}
DB_PASS=${DB_PASS:-""}
DB_NAME=${DB_NAME:-"go_sub"}

# 检查 MySQL 命令
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}❌ mysql 命令未找到${NC}"
    echo -e "${BLUE}ℹ️  请安装 MySQL 客户端${NC}"
    exit 1
fi

# 检查 init.sql 文件
if [ ! -f "init.sql" ]; then
    echo -e "${RED}❌ init.sql 文件不存在${NC}"
    exit 1
fi

# 构建 MySQL 命令
MYSQL_CMD="mysql -h $DB_HOST -P $DB_PORT -u $DB_USER"
if [ -n "$DB_PASS" ]; then
    MYSQL_CMD="$MYSQL_CMD -p$DB_PASS"
fi

echo "连接信息："
echo "  主机: $DB_HOST:$DB_PORT"
echo "  用户: $DB_USER"
echo "  数据库: $DB_NAME"
echo ""

# 测试连接
echo "测试 MySQL 连接..."
if ! $MYSQL_CMD -e "SELECT 1" &> /dev/null; then
    echo -e "${RED}❌ MySQL 连接失败${NC}"
    echo ""
    echo "请检查："
    echo "  1. MySQL 是否已启动"
    echo "  2. 用户名密码是否正确"
    echo ""
    echo "如果需要密码，请设置环境变量："
    echo "  export DB_PASS='your-password'"
    echo "  bash scripts/init-local-db.sh"
    exit 1
fi
echo -e "${GREEN}✓ MySQL 连接成功${NC}"

# 检查数据库是否存在
echo ""
echo "检查数据库 '$DB_NAME'..."
if $MYSQL_CMD -e "USE $DB_NAME" &> /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  数据库 '$DB_NAME' 已存在${NC}"
    read -p "是否要重新初始化？这将删除所有数据！(y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "取消初始化"
        exit 0
    fi
    echo "删除现有数据库..."
    $MYSQL_CMD -e "DROP DATABASE IF EXISTS $DB_NAME"
fi

# 创建数据库
echo "创建数据库 '$DB_NAME'..."
if $MYSQL_CMD -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"; then
    echo -e "${GREEN}✓ 数据库创建成功${NC}"
else
    echo -e "${RED}❌ 数据库创建失败${NC}"
    exit 1
fi

# 导入 SQL
echo ""
echo "导入数据表结构..."
if $MYSQL_CMD $DB_NAME < init.sql; then
    echo -e "${GREEN}✓ 数据表导入成功${NC}"
else
    echo -e "${RED}❌ 数据表导入失败${NC}"
    exit 1
fi

# 验证
echo ""
echo "验证数据库..."
TABLE_COUNT=$($MYSQL_CMD -N -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='$DB_NAME'")
echo -e "${GREEN}✓ 成功创建 $TABLE_COUNT 个表${NC}"

# 显示表列表
echo ""
echo "数据表列表："
$MYSQL_CMD -e "USE $DB_NAME; SHOW TABLES;"

echo ""
echo -e "${GREEN}✅ 数据库初始化完成！${NC}"
echo ""
echo "下一步："
echo "  1. 修改 config.yaml 中的数据库配置"
echo "  2. 运行: make dev 或 bash scripts/dev.sh"
