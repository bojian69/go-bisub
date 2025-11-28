# 🎯 从这里开始

欢迎使用 GO BI Subscription 订阅服务！

## 最快启动方式

```bash
# 一键启动（需要本地 MySQL 和 Redis）
make install-tools && make deps && cp config.local.yaml config.yaml && make db-init && make dev
```

然后访问 http://localhost:8080/admin (admin/admin123)

---

## 📚 文档导航

### 新手入门
1. **[README.md](../README.md)** - 项目概览和完整文档
2. **[快速启动指南](QUICKSTART.md)** - 快速启动指南（推荐首先阅读）
3. **[命令速查表](COMMANDS.md)** - 命令速查表

### 开发指南
4. **[本地开发指南](LOCAL_DEVELOPMENT.md)** - 本地开发详细指南
5. **[数据库迁移指南](DATABASE_MIGRATION.md)** - 数据库迁移说明

### 技术文档
6. **[技术文档](./)** - 详细技术文档

---

## 🚀 三种启动方式

### 方式 1: 本地开发（推荐）

适合日常开发，性能最好。

```bash
# 1. 准备环境
make install-tools
make deps

# 2. 配置
cp config.local.yaml config.yaml
# 编辑 config.yaml，修改数据库密码

# 3. 初始化数据库
make db-init

# 4. 启动
make dev
```

**前提条件**: 本地已安装 MySQL 和 Redis

### 方式 2: Docker Compose

适合快速体验，无需本地安装数据库。

```bash
# 一键启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f go-bisub
```

**前提条件**: 已安装 Docker 和 Docker Compose

### 方式 3: 混合模式

本地代码 + Docker 数据库。

```bash
# 1. 启动数据库
docker-compose up -d mysql redis

# 2. 初始化
mysql -h 127.0.0.1 -u root -ppassword < init.sql

# 3. 启动应用
make dev
```

---

## 🔧 常用命令

```bash
# 启动服务
make dev                # 热重载开发
make start              # 快速启动

# 数据库
make db-check           # 检查连接
make db-init            # 初始化

# 代码质量
make check              # 完整检查
make test               # 运行测试

# 查看帮助
make help               # 所有命令
```

更多命令查看 [命令速查表](COMMANDS.md)

---

## 🌐 访问地址

启动成功后访问：

- **API 文档**: http://localhost:8080
- **管理界面**: http://localhost:8080/admin
  - 用户名: `admin`
  - 密码: `admin123`
- **健康检查**: http://localhost:8080/health

---

## ❓ 遇到问题？

### 常见问题快速解决

| 问题 | 解决方案 |
|------|---------|
| `air: command not found` | `make install-tools` |
| MySQL 连接失败 | `make db-check` 检查状态 |
| 数据库不存在 | `make db-init` 初始化 |
| 端口被占用 | 修改 `config.yaml` 中的端口 |

详细故障排查: [快速启动指南](QUICKSTART.md#故障排查)

---

## 📖 学习路径

### 第一天：快速上手
1. 阅读 [README.md](README.md) 了解项目
2. 按照 [QUICKSTART.md](QUICKSTART.md) 启动项目
3. 访问管理界面，创建第一个订阅

### 第二天：深入开发
1. 阅读 [LOCAL_DEVELOPMENT.md](LOCAL_DEVELOPMENT.md)
2. 了解项目结构和代码规范
3. 尝试修改代码并测试

### 第三天：进阶使用
1. 学习 API 文档
2. 了解数据库表结构
3. 编写测试用例

---

## 🎓 核心概念

### 订阅 (Subscription)
- 定义 SQL 查询模板
- 支持变量替换
- 版本控制

### 执行 (Execute)
- 传入变量执行订阅
- 返回查询结果
- 记录执行统计

### 统计 (Stats)
- 执行次数
- 平均耗时
- 性能分析

---

## 🛠️ 开发工具推荐

### IDE
- **VSCode** + Go 插件
- **GoLand** (JetBrains)

### 数据库客户端
- **DBeaver** (免费)
- **MySQL Workbench**
- **TablePlus** (macOS)

### API 测试
- **Postman**
- **Insomnia**
- **curl**

---

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交代码
4. 运行 `make check` 和 `make test`
5. 提交 Pull Request

---

## 📞 获取帮助

- 查看文档: [README.md](../README.md)
- 命令速查: [命令速查表](COMMANDS.md)
- 故障排查: [快速启动指南](QUICKSTART.md)
- 提交 Issue

---

## ⭐ 下一步

- [ ] 启动项目
- [ ] 访问管理界面
- [ ] 创建第一个订阅
- [ ] 执行订阅并查看结果
- [ ] 查看统计报表
- [ ] 阅读 API 文档
- [ ] 编写测试用例

**祝你使用愉快！** 🎉
