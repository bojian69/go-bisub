# uhomes Go 技术栈选型与依赖清单

**文档版本：** 1.0
**最后更新：** 2025-02-19
**状态：** 草案 (Draft)

---

## 1. 概述

本文档旨在规范 uhomes 内部 Go 项目的技术选型，明确“核心技术栈”与“推荐工具包”。

- **核心技术栈 (Core Stack)**：所有微服务项目**必须**使用的基础库，由架构组统一维护，确保架构一致性。
- **推荐工具包 (Recommended Toolkit)**：经过验证的、高质量的第三方库，建议优先选用，避免重复造轮子。

---

## 2. 核心技术栈 (Core Stack)

这些库构成了我们微服务架构的基石，通常已封装在 `gox` 基础库中。

| 领域 | 库/框架 | 版本策略 | 说明 |
| :--- | :--- | :--- | :--- |
| **Web 框架** | [Gin](https://github.com/gin-gonic/gin) | v1.11+ | 高性能 HTTP Web 框架，生态成熟，中间件丰富。 |
| **RPC 框架** | [gRPC-Go](https://github.com/grpc/grpc-go) | v1.77+ | Google 官方 Go 实现，用于服务间高性能通信。 |
| **日志** | [slog](https://pkg.go.dev/log/slog) + [zap](https://github.com/uber-go/zap) | Go 1.21+ | **推荐组合**：前端使用官方 `slog` 接口，后端使用 `zap` 作为高性能实现（通过 `exp/slog` 适配器）。 |
| **配置管理** | [Viper](https://github.com/spf13/viper) | v1.21+ | 支持多种格式（YAML/JSON）和远程配置中心（Nacos）。 |
| **ORM** | [GORM](https://github.com/go-gorm/gorm) | v1.31+ | 功能强大的 ORM，支持 MySQL/PostgreSQL。**注意**：复杂查询建议手写 SQL 或使用 sqlx。 |
| **缓存** | [go-redis](https://github.com/redis/go-redis) | v9+ | 官方推荐的 Redis 客户端，支持 Cluster 和 Sentinel。 |
| **链路追踪** | [OpenTelemetry](https://pkg.go.dev/go.opentelemetry.io/otel) | v1.38+ | 云原生可观测性标准，用于分布式链路追踪。 |
| **依赖注入** | [Fx](https://github.com/uber-go/fx) | v1.24+ | Uber 出品的依赖注入框架，提供模块化和生命周期管理。 |
| **Protobuf** | [Protobuf-Go](https://google.golang.org/protobuf) | v1.36+ | 官方 Protobuf 运行时库。 |

---

## 3. 推荐工具包 (Recommended Toolkit)

以下库在特定场景下强烈推荐使用，它们能显著提升开发效率和代码质量。

### 3.1 基础工具 (Utilities)

| 库名 | 用途 | 推荐理由 |
| :--- | :--- | :--- |
| **[lo](https://github.com/samber/lo)** | Lodash for Go | 提供了大量切片、Map 操作的泛型函数（如 Filter, Map, Uniq），避免手写循环。 |
| **[cast](https://github.com/spf13/cast)** | 类型转换 | 安全且方便地在不同类型间转换（如 string 转 int，interface{} 转 string）。 |
| **[mergo](https://github.com/imdario/mergo)** | 结构体合并 | 用于合并配置或结构体，支持深度合并。 |
| **[uuid](https://github.com/google/uuid)** | UUID 生成 | Google 官方 UUID 库，标准且高效。 |
| **[validator](https://github.com/go-playground/validator)** | 数据校验 | 结构体标签校验库，Gin 默认集成。 |

### 3.2 并发与流程控制

| 库名 | 用途 | 推荐理由 |
| :--- | :--- | :--- |
| **[conc](https://github.com/sourcegraph/conc)** | 结构化并发 | 比 `sync.WaitGroup` 更易用、更安全的并发池，自动处理 Panic 恢复。 |
| **[ants](https://github.com/panjf2000/ants)** | 协程池 | 高性能协程池，适用于需要限制 Goroutine 数量的高并发场景。 |
| **[retry-go](https://github.com/avast/retry-go)** | 重试机制 | 简单易用的函数重试库，支持自定义策略（Backoff, Jitter）。 |

### 3.3 测试工具

| 库名 | 用途 | 推荐理由 |
| :--- | :--- | :--- |
| **[testify](https://github.com/stretchr/testify)** | 断言库 | 提供了 `assert` 和 `require`，让测试代码更可读。 |
| **[mock](https://github.com/uber-go/mock)** | Mock 工具 | Uber 维护的 gomock 分支，用于生成接口的 Mock 实现。 |

### 3.4 开发效率工具 (Dev Tools)

| 工具名 | 用途 | 安装方式 |
| :--- | :--- | :--- |
| **[Air](https://github.com/cosmtrek/air)** | 热重载 | `go install github.com/cosmtrek/air@latest` <br> 开发时自动监听文件变更并重启服务。 |
| **[Buf](https://buf.build)** | Proto 管理 | `brew install bufbuild/buf/buf` <br> 现代化的 Protobuf/gRPC 构建与 Lint 工具。 |
| **[GolangCI-Lint](https://golangci-lint.run)** | 代码检查 | `brew install golangci-lint` <br> 聚合了数十种 Linter，CI 流程必备。 |

---

## 4. 引入新依赖的原则

在引入上述清单之外的第三方库时，请遵循以下原则：

1.  **许可证合规**：必须是 MIT, Apache 2.0, BSD 等商业友好协议。**严禁使用 GPL/AGPL 协议库**。
2.  **维护活跃度**：检查 GitHub Stars (>100)、最近提交时间（半年内）、Issue 响应情况。
3.  **依赖树复杂度**：避免引入依赖过重（Dependency Hell）的库。
4.  **功能重叠**：如果标准库或 `gox` 已有类似功能，优先使用现有方案。

---

## 5. 附录：版本锁定文件 (go.mod)

所有项目必须提交 `go.mod` 和 `go.sum`，严禁手动修改 `go.sum`。
建议在 CI 中添加 `go mod verify` 步骤以确保依赖完整性。
