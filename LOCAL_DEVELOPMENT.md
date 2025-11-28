# 本地开发指南

本指南专门针对使用本地 MySQL 和 Redis 进行开发的场景。

## 前置要求

### 必需
- Go 1.21+
- MySQL 8.0+ (本地安装)
- Redis 6.0+ (本地安装)

### 可选
- Docker (如果想使用容器化服务)

## 快速开始

### 1. 安装 MySQL 和 Redis

#### macOS
```bash
# 使用 Homebrew
brew install mysql redis

# 启动服务
brew services start mysql
brew services start redis
```

#### Linux (Ubuntu/Debian)
```bash
# 安装 MySQL
sudo apt update
sudo apt install mysql-server

# 安装 Redis
sudo apt install redis-server

# 启动服务
sudo systemctl start mysql
sudo systemctl start redis
```

#### Windows
- MySQL: 下载并安装 [MySQL Installer](https://dev.mysql.com/downloads/installer/)
- Redis: 下载并安装 [Redis for Windows](https://github.com/microsoftarchive/redis/releases)

### 2. 配置数据库

```bash
# 登录 MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE go_sub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 创建用户（可选）
CREATE USER 'sub_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON go_sub.* TO 'sub_user'@'localhost';
FLUSH PRIVILEGES;

# 退出
EXIT;
```

### 3. 初始化项目

```bash
# 克隆项目
cd go-bisub

# 安装开发工具
make install-tools

# 下载依赖
make deps

# 复制配置文件
cp config.local.yaml config.yaml
```

### 4. 修改配置

编辑 `config.yaml`，修改数据库连接信息：

```yaml
database:
  primary:
    host: 127.0.0.1
    port: 3306
    database: go_sub
    username: root
    password: "your_mysql_password"  # 修改为你的密码
```

### 5. 初始化数据库

```bash
# 方式 1: 使用 make 命令（推荐）
make db-init

# 方式 2: 手动导入
mysql -h 127.0.0.1 -u root -p go_sub < init.sql

# 方式 3: 使用脚本
DB_PASS=your_password bash scripts/init-local-db.sh
```

### 6. 验证连接

```bash
# 检查数据库和 Redis 连接
make db-check
```

### 7. 启动开发服务器

```bash
# 方式 1: 使用 make（推荐）
make dev

# 方式 2: 使用脚本（自动检查依赖）
bash scripts/dev.sh

# 方式 3: 快速启动（无热重载）
make start
```

## 常见问题

### MySQL 连接失败

**问题**: `ERROR 2002 (HY000): Can't connect to local MySQL server`

**解决方案**:
```bash
# 检查 MySQL 是否运行
# macOS
brew services list | grep mysql

# Linux
sudo systemctl status mysql

# 启动 MySQL
# macOS
brew services start mysql

# Linux
sudo systemctl start mysql
```

### MySQL 密码问题

**问题**: `ERROR 1045 (28000): Access denied for user 'root'@'localhost'`

**解决方案**:
```bash
# 重置 MySQL root 密码
# macOS
mysql.server stop
mysqld_safe --skip-grant-tables &
mysql -u root
# 在 MySQL 中执行:
ALTER USER 'root'@'localhost' IDENTIFIED BY 'new_password';
FLUSH PRIVILEGES;
EXIT;
# 重启 MySQL
mysql.server start

# Linux
sudo mysql
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'new_password';
FLUSH PRIVILEGES;
EXIT;
```

### Redis 连接失败

**问题**: Redis 连接超时或拒绝

**解决方案**:
```bash
# 检查 Redis 是否运行
redis-cli ping
# 应该返回 PONG

# 如果没有响应，启动 Redis
# macOS
brew services start redis

# Linux
sudo systemctl start redis

# 手动启动
redis-server
```

### 端口被占用

**问题**: `bind: address already in use`

**解决方案**:
```bash
# 查找占用端口的进程
lsof -i :8080

# 杀死进程
kill -9 <PID>

# 或者修改配置文件中的端口
vim config.yaml
# 修改 server.port 为其他端口，如 8081
```

### 数据库表不存在

**问题**: `Table 'go_sub.xxx' doesn't exist`

**解决方案**:
```bash
# 重新初始化数据库
make db-init

# 或手动导入
mysql -h 127.0.0.1 -u root -p go_sub < init.sql
```

## 开发工作流

### 日常开发

```bash
# 1. 启动开发服务器
make dev

# 2. 修改代码
# Air 会自动检测文件变化并重新编译

# 3. 测试 API
curl http://localhost:8080/health

# 4. 查看日志
# 日志会实时显示在终端
```

### 代码提交前

```bash
# 1. 格式化代码
make fmt

# 2. 运行代码检查
make lint

# 3. 运行测试
make test

# 4. 完整检查
make check
```

### 数据库操作

```bash
# 检查连接
make db-check

# 重新初始化（会删除所有数据！）
make db-init

# 连接到数据库
mysql -h 127.0.0.1 -u root -p go_sub

# 查看表结构
mysql -h 127.0.0.1 -u root -p go_sub -e "SHOW TABLES;"
```

## 性能优化建议

### MySQL 配置优化

编辑 MySQL 配置文件 (`/etc/my.cnf` 或 `/usr/local/etc/my.cnf`):

```ini
[mysqld]
# 连接数
max_connections = 200

# 缓冲池大小（根据内存调整）
innodb_buffer_pool_size = 1G

# 日志文件大小
innodb_log_file_size = 256M

# 查询缓存
query_cache_size = 64M
query_cache_type = 1
```

### Redis 配置优化

编辑 Redis 配置文件 (`/etc/redis/redis.conf` 或 `/usr/local/etc/redis.conf`):

```conf
# 最大内存
maxmemory 256mb

# 内存淘汰策略
maxmemory-policy allkeys-lru

# 持久化（开发环境可以关闭以提高性能）
save ""
appendonly no
```

## 调试技巧

### 使用 Delve 调试器

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug cmd/server/main.go

# 在代码中设置断点
(dlv) break main.main
(dlv) continue
```

### 查看实时日志

```bash
# 如果使用 make dev，日志会实时显示

# 如果日志输出到文件
tail -f logs/app.log

# 过滤特定日志
tail -f logs/app.log | grep ERROR
```

### 性能分析

```bash
# 启动服务后，访问 pprof
go tool pprof http://localhost:8080/debug/pprof/profile

# 或使用 make 命令
make profile-cpu
make profile-mem
```

## 环境变量

可以使用环境变量覆盖配置文件：

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件
vim .env

# 使用环境变量启动
export $(cat .env | xargs)
make dev
```

## 多环境配置

```bash
# 开发环境
cp config.local.yaml config.yaml

# 测试环境
cp config.test.yaml config.yaml

# 生产环境
cp config.prod.yaml config.yaml
```

## 有用的命令

```bash
# 查看所有可用命令
make help

# 检查环境
make check-env

# 检查数据库
make db-check

# 运行测试
make test

# 生成测试覆盖率
make test-coverage

# 构建二进制文件
make build

# 清理构建文件
make clean
```

## 推荐工具

### IDE 插件
- **VSCode**: Go extension, MySQL extension
- **GoLand**: 内置 Go 和数据库支持

### 数据库客户端
- **DBeaver**: 免费、跨平台
- **MySQL Workbench**: MySQL 官方工具
- **TablePlus**: macOS 推荐

### API 测试
- **Postman**: 功能强大
- **Insomnia**: 简洁易用
- **curl**: 命令行工具

### Redis 客户端
- **RedisInsight**: Redis 官方工具
- **Another Redis Desktop Manager**: 开源免费

## 获取帮助

- 查看 [README.md](README.md) 了解完整功能
- 查看 [QUICKSTART.md](QUICKSTART.md) 快速开始
- 查看 [docs/](docs/) 目录了解详细文档
- 提交 Issue 或 PR

## 下一步

- 阅读 API 文档了解接口详情
- 访问 Web 管理界面体验功能
- 开始开发你的第一个订阅服务
- 编写测试用例
