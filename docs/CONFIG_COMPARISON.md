# 配置文件对比速查表

## 🎯 一句话总结

- **本地开发**：只需要 `config.yaml`，不需要 `.env`
- **Docker 部署**：需要 `.env` + `config.docker.yaml`

## 📊 配置文件对比

### 文件用途

```
┌─────────────────────────────────────────────────────────────┐
│ .env.example                                                │
│ ├─ 用途：Docker Compose 环境变量模板                        │
│ ├─ 读取：docker-compose.yml                                │
│ ├─ 提交：✅ 是（模板文件）                                   │
│ └─ 场景：Docker 部署                                        │
├─────────────────────────────────────────────────────────────┤
│ .env                                                        │
│ ├─ 用途：Docker Compose 实际环境变量                        │
│ ├─ 读取：docker-compose.yml                                │
│ ├─ 提交：❌ 否（包含密码）                                   │
│ └─ 场景：Docker 部署                                        │
├─────────────────────────────────────────────────────────────┤
│ config.local.yaml                                           │
│ ├─ 用途：本地开发配置模板                                   │
│ ├─ 读取：无（仅作为模板）                                   │
│ ├─ 提交：✅ 是（模板文件）                                   │
│ └─ 场景：本地开发参考                                       │
├─────────────────────────────────────────────────────────────┤
│ config.docker.yaml                                          │
│ ├─ 用途：Docker 环境应用配置                                │
│ ├─ 读取：Go 应用（挂载为 config.yaml）                      │
│ ├─ 提交：✅ 是（不含敏感信息）                               │
│ └─ 场景：Docker 部署                                        │
├─────────────────────────────────────────────────────────────┤
│ config.yaml                                                 │
│ ├─ 用途：实际应用配置                                       │
│ ├─ 读取：Go 应用                                            │
│ ├─ 提交：❌ 否（包含密码）                                   │
│ └─ 场景：本地开发                                           │
└─────────────────────────────────────────────────────────────┘
```

## 🔄 不同场景的配置

### 场景 1：本地开发（直接运行 Go）

```bash
# 需要的文件
config.yaml  # ✅ 必需

# 不需要的文件
.env         # ❌ 不需要（Docker Compose 才用）

# 配置示例
database:
  host: 127.0.0.1  # 本地 MySQL
redis:
  host: 127.0.0.1  # 本地 Redis
```

### 场景 2：Docker 开发

```bash
# 需要的文件
.env                # ✅ 必需（Docker Compose 变量）
config.docker.yaml  # ✅ 必需（应用配置）

# 不需要的文件
config.yaml         # ❌ 不需要（被 config.docker.yaml 替代）

# 配置示例
# config.docker.yaml
database:
  host: mysql  # Docker 服务名
redis:
  host: redis  # Docker 服务名
```

### 场景 3：混合模式（Docker 数据库 + 本地应用）

```bash
# 需要的文件
.env         # ✅ 必需（启动 Docker 数据库）
config.yaml  # ✅ 必需（应用配置）

# 配置示例
# config.yaml
database:
  host: 127.0.0.1  # Docker 映射到本地
  port: 3306       # 映射端口
  password: password  # 来自 .env 的 MYSQL_ROOT_PASSWORD
```

## 📝 配置内容对比

### 数据库连接

| 场景 | 文件 | host | 说明 |
|------|------|------|------|
| 本地开发 | config.yaml | `127.0.0.1` | 本机 MySQL |
| Docker 内 | config.docker.yaml | `mysql` | Docker 服务名 |
| Docker 外 | config.yaml | `127.0.0.1` | 端口映射 |
| 远程数据库 | config.yaml | `db.example.com` | 实际地址 |

### Redis 连接

| 场景 | 文件 | host | 说明 |
|------|------|------|------|
| 本地开发 | config.yaml | `127.0.0.1` | 本机 Redis |
| Docker 内 | config.docker.yaml | `redis` | Docker 服务名 |
| Docker 外 | config.yaml | `127.0.0.1` | 端口映射 |
| 远程 Redis | config.yaml | `redis.example.com` | 实际地址 |

## 🎯 快速决策树

```
你要做什么？
│
├─ 本地开发（直接运行 Go）
│  └─ 使用：config.yaml
│     └─ host: 127.0.0.1
│
├─ Docker 部署
│  └─ 使用：.env + config.docker.yaml
│     └─ host: mysql / redis
│
└─ Docker 数据库 + 本地应用
   └─ 使用：.env + config.yaml
      └─ host: 127.0.0.1
```

## ⚡ 快速配置命令

### 本地开发

```bash
# 1. 复制配置
cp config.local.yaml config.yaml

# 2. 修改数据库连接
vim config.yaml
# host: 127.0.0.1
# password: 你的密码

# 3. 启动
make dev
```

### Docker 部署

```bash
# 1. 复制环境变量
cp .env.example .env

# 2. 修改密码
vim .env
# MYSQL_ROOT_PASSWORD=your_password
# JWT_SECRET=your_secret

# 3. 启动
make docker-up
```

### 混合模式

```bash
# 1. 启动 Docker 数据库
docker-compose up -d mysql redis

# 2. 配置应用
cp config.local.yaml config.yaml
vim config.yaml
# host: 127.0.0.1
# password: password (来自 .env)

# 3. 启动应用
make dev
```

## 🔍 配置验证

### 检查你的配置是否正确

```bash
# 本地开发
✅ 有 config.yaml
✅ config.yaml 中 host: 127.0.0.1
✅ MySQL 和 Redis 在本地运行

# Docker 部署
✅ 有 .env
✅ 有 config.docker.yaml
✅ config.docker.yaml 中 host: mysql/redis
✅ docker-compose.yml 挂载 config.docker.yaml

# 混合模式
✅ 有 .env
✅ 有 config.yaml
✅ config.yaml 中 host: 127.0.0.1
✅ Docker 数据库在运行
✅ 端口已映射（3306, 6379）
```

## ❌ 常见错误

### 错误 1：本地开发使用 Docker 服务名

```yaml
# ❌ 错误
database:
  host: mysql  # 本地无法解析 Docker 服务名

# ✅ 正确
database:
  host: 127.0.0.1
```

### 错误 2：Docker 环境使用 localhost

```yaml
# ❌ 错误（在 config.docker.yaml 中）
database:
  host: 127.0.0.1  # Docker 容器内的 localhost 是容器自己

# ✅ 正确
database:
  host: mysql  # 使用 Docker 服务名
```

### 错误 3：以为 .env 会被 Go 应用读取

```bash
# ❌ 错误理解
.env 中的 DB_HOST 会被 Go 应用使用

# ✅ 正确理解
.env 只被 docker-compose.yml 使用
Go 应用读取 config.yaml
```

## 📚 详细文档

- [本地环境配置](docs/LOCAL_ENV_SETUP.md)
- [配置文件指南](docs/CONFIGURATION_GUIDE.md)
- [Docker 快速开始](DOCKER_QUICKSTART.md)

## 🆘 还是不清楚？

### 问自己三个问题：

1. **我要在哪里运行应用？**
   - 本地 → 用 config.yaml
   - Docker → 用 .env + config.docker.yaml

2. **数据库在哪里？**
   - 本地 → host: 127.0.0.1
   - Docker → host: mysql
   - 远程 → host: 实际地址

3. **我需要 .env 吗？**
   - 本地开发 → 不需要
   - Docker 部署 → 需要
