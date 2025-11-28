# 数据库名称变更说明

## 变更内容

数据库名称已从 `go_bisub` 更改为 `go_sub`

## 影响的文件

已更新以下文件中的数据库名称：

### 配置文件
- ✅ `config.yaml` - 主配置文件
- ✅ `config.local.yaml` - 本地开发配置模板
- ✅ `.env.example` - 环境变量示例
- ✅ `docker-compose.yml` - Docker 配置

### 脚本文件
- ✅ `scripts/check-db.sh` - 数据库检查脚本
- ✅ `scripts/init-local-db.sh` - 数据库初始化脚本
- ✅ `Makefile` - 构建脚本

### 文档文件
- ✅ `README.md` - 项目说明
- ✅ `LOCAL_DEVELOPMENT.md` - 本地开发指南
- ✅ `init.sql` - 数据库初始化 SQL

## 迁移步骤

### 如果你已经有旧数据库 `go_bisub`

#### 方式 1: 重命名数据库（保留数据）

```bash
# 1. 导出旧数据库
mysqldump -u root -p go_bisub > go_bisub_backup.sql

# 2. 创建新数据库
mysql -u root -p -e "CREATE DATABASE go_sub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"

# 3. 导入数据
mysql -u root -p go_sub < go_bisub_backup.sql

# 4. 验证数据
mysql -u root -p go_sub -e "SHOW TABLES;"

# 5. 删除旧数据库（可选）
mysql -u root -p -e "DROP DATABASE go_bisub"
```

#### 方式 2: 全新初始化（不保留数据）

```bash
# 1. 删除旧数据库（如果存在）
mysql -u root -p -e "DROP DATABASE IF EXISTS go_bisub"

# 2. 初始化新数据库
make db-init

# 或手动执行
mysql -u root -p < init.sql
```

### 如果是新项目

直接初始化即可：

```bash
make db-init
```

## 验证

检查数据库是否正确创建：

```bash
# 使用 make 命令
make db-check

# 或手动检查
mysql -u root -p -e "SHOW DATABASES LIKE 'go_sub'"
mysql -u root -p go_sub -e "SHOW TABLES"
```

## 配置更新

确保你的 `config.yaml` 中的数据库名称已更新：

```yaml
database:
  primary:
    host: 127.0.0.1
    port: 3306
    database: go_sub  # 确保是 go_sub
    username: root
    password: your_password
```

## 常见问题

### Q: 启动时提示数据库不存在

**A**: 运行 `make db-init` 初始化数据库

### Q: 如何保留旧数据？

**A**: 参考上面的"方式 1: 重命名数据库"步骤

### Q: Docker 环境如何更新？

**A**: 
```bash
# 停止并删除旧容器
docker-compose down -v

# 重新启动（会自动创建 go_sub 数据库）
docker-compose up -d
```

## 回滚

如果需要回滚到旧的数据库名称：

```bash
# 1. 备份当前数据
mysqldump -u root -p go_sub > go_sub_backup.sql

# 2. 创建旧数据库
mysql -u root -p -e "CREATE DATABASE go_bisub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"

# 3. 导入数据
mysql -u root -p go_bisub < go_sub_backup.sql

# 4. 修改配置文件中的数据库名称
# 将所有 go_sub 改回 go_bisub
```

## 注意事项

1. **备份数据**: 在进行任何数据库操作前，请先备份数据
2. **更新配置**: 确保所有配置文件都已更新
3. **重启服务**: 修改配置后需要重启应用服务
4. **测试验证**: 迁移后进行完整的功能测试

## 相关命令

```bash
# 检查数据库连接
make db-check

# 初始化数据库
make db-init

# 查看数据库信息
mysql -u root -p go_sub -e "SELECT DATABASE(); SHOW TABLES;"

# 查看表结构
mysql -u root -p go_sub -e "DESCRIBE sub_subscription_theme;"
```
