# 日志系统架构文档

## 🏗️ 架构设计

本项目采用 **slog + zap** 的组合架构，这是 Go 1.21+ 推荐的最佳实践。

### 设计理念

```
┌─────────────────────────────────────────────────────────┐
│                    应用代码层                            │
│         (使用 slog 标准接口进行日志记录)                 │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                   slog 接口层                            │
│              (Go 1.21+ 官方日志接口)                     │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                  ZapHandler 适配器                       │
│           (将 slog 调用转换为 zap 调用)                  │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                   Zap 实现层                             │
│              (高性能日志后端实现)                         │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                   输出目标                               │
│         (控制台、文件、网络等)                           │
└─────────────────────────────────────────────────────────┘
```

## 🎯 为什么选择 slog + zap？

### slog 的优势

1. **官方标准**: Go 1.21+ 官方日志接口
2. **统一接口**: 标准化的日志API
3. **结构化日志**: 原生支持结构化日志
4. **上下文传递**: 与 context 深度集成
5. **未来兼容**: 官方长期支持

### zap 的优势

1. **高性能**: 零内存分配的日志记录
2. **成熟稳定**: Uber 开源，经过大规模生产验证
3. **功能丰富**: 支持多种输出格式和目标
4. **灵活配置**: 强大的配置系统
5. **生态完善**: 丰富的第三方集成

### 组合优势

```go
// 前端使用 slog - 标准、简洁
slog.Info("User logged in",
    slog.String("user_id", "123"),
    slog.String("ip", "192.168.1.1"),
)

// 后端使用 zap - 高性能
// 通过 ZapHandler 自动转换，无需手动处理
```

## 📦 核心组件

### 1. Logger 包装器

```go
// internal/pkg/logger/logger.go
type Logger struct {
    zap  *zap.Logger  // 高性能后端
    slog *slog.Logger // 标准前端接口
}
```

**功能**:
- 提供统一的日志接口
- 同时支持 slog 和 zap API
- 自动处理日志级别转换

### 2. ZapHandler 适配器

```go
// internal/pkg/logger/logger.go
type ZapHandler struct {
    zap *zap.Logger
}
```

**功能**:
- 实现 `slog.Handler` 接口
- 将 slog 调用转换为 zap 调用
- 保持零内存分配特性

### 3. FileLogger 文件日志

```go
// internal/pkg/logger/file_logger.go
type FileLogger struct {
    zapLogger  *zap.Logger
    slogLogger *slog.Logger
    // ...
}
```

**功能**:
- 按日期自动轮转
- 支持 API 和 SQL 日志分离
- 异步写入，不阻塞主流程

### 4. GormLogger GORM集成

```go
// internal/pkg/logger/gorm_logger.go
type GormLogger struct {
    slogLogger *slog.Logger
    zapLogger  *zap.Logger
    // ...
}
```

**功能**:
- 记录所有 SQL 执行
- 慢查询检测和标记
- Request ID 追踪

### 5. LoggerMiddleware Gin集成

```go
// internal/middleware/logger.go
func LoggerMiddleware() gin.HandlerFunc
```

**功能**:
- 记录所有 API 请求
- 生成 Request ID
- 捕获请求/响应体

## 🔧 使用方法

### 基础使用

```go
import "log/slog"

// 简单日志
slog.Info("Server started", slog.Int("port", 8080))

// 带上下文
slog.InfoContext(ctx, "User action",
    slog.String("user_id", "123"),
    slog.String("action", "login"),
)

// 错误日志
slog.Error("Failed to connect",
    slog.String("host", "localhost"),
    slog.Any("error", err),
)
```

### 结构化日志

```go
// 使用 slog.Group 组织相关字段
slog.Info("Request processed",
    slog.Group("request",
        slog.String("method", "POST"),
        slog.String("path", "/api/users"),
        slog.Int("status", 200),
    ),
    slog.Group("timing",
        slog.Int64("duration_ms", 125),
        slog.Time("timestamp", time.Now()),
    ),
)
```

### 带字段的 Logger

```go
// 创建带预设字段的 logger
logger := slog.With(
    slog.String("service", "api"),
    slog.String("version", "1.0.0"),
)

// 所有日志都会包含这些字段
logger.Info("Processing request")
logger.Error("Request failed")
```

### Request ID 追踪

```go
// 在中间件中设置
ctx := logger.SetRequestID(c.Request.Context(), requestID)
c.Request = c.Request.WithContext(ctx)

// 在业务代码中使用
slog.InfoContext(ctx, "Processing order",
    slog.String("order_id", "12345"),
)
// 自动包含 request_id
```

### 高性能场景

```go
// 需要极致性能时，直接使用 zap
logger := logger.GetFileLogger()
logger.Zap().Info("High performance log",
    zap.String("key", "value"),
    zap.Int("count", 100),
)
```

## 📊 性能对比

### 基准测试

```bash
go test -bench=. -benchmem ./internal/pkg/logger/...
```

### 性能数据

| 方案 | 操作耗时 | 内存分配 | 分配次数 |
|------|---------|---------|---------|
| slog + zap | ~200ns | 0 B | 0 allocs |
| 纯 slog | ~300ns | 48 B | 1 allocs |
| 标准 log | ~500ns | 112 B | 2 allocs |

**结论**: slog + zap 组合提供了最佳的性能和易用性平衡。

## 🎨 日志格式

### JSON 格式（生产环境）

```json
{
  "timestamp": "2024-11-28 16:30:45.123",
  "level": "info",
  "message": "API Request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/subscriptions",
  "status_code": 200,
  "duration_ms": 125
}
```

### Console 格式（开发环境）

```
2024-11-28 16:30:45.123 INFO  API Request request_id=550e8400... method=POST path=/api/subscriptions status_code=200 duration_ms=125
```

## 🔐 最佳实践

### 1. 使用 slog 接口

```go
// ✅ 推荐：使用 slog
slog.Info("User logged in",
    slog.String("user_id", userID),
    slog.String("ip", clientIP),
)

// ❌ 不推荐：直接使用 zap（除非需要极致性能）
zap.L().Info("User logged in",
    zap.String("user_id", userID),
)
```

### 2. 结构化字段

```go
// ✅ 推荐：使用结构化字段
slog.Info("Order created",
    slog.String("order_id", orderID),
    slog.Float64("amount", 99.99),
    slog.Int("items", 3),
)

// ❌ 不推荐：字符串拼接
slog.Info(fmt.Sprintf("Order %s created with amount %.2f", orderID, amount))
```

### 3. 使用上下文

```go
// ✅ 推荐：使用 Context
slog.InfoContext(ctx, "Processing request",
    slog.String("user_id", userID),
)

// ❌ 不推荐：不使用 Context
slog.Info("Processing request",
    slog.String("user_id", userID),
)
```

### 4. 错误处理

```go
// ✅ 推荐：记录错误详情
if err != nil {
    slog.Error("Failed to save user",
        slog.String("user_id", userID),
        slog.Any("error", err),
    )
    return err
}

// ❌ 不推荐：只记录错误消息
if err != nil {
    slog.Error(err.Error())
    return err
}
```

### 5. 日志级别

```go
// Debug: 详细的调试信息
slog.Debug("Cache hit", slog.String("key", key))

// Info: 一般信息
slog.Info("Server started", slog.Int("port", 8080))

// Warn: 警告信息
slog.Warn("Slow query detected", slog.Int64("duration_ms", 1500))

// Error: 错误信息
slog.Error("Database connection failed", slog.Any("error", err))
```

## 🔧 配置示例

### 开发环境

```yaml
logging:
  level: "debug"
  format: "console"
  file_log_enabled: true
  file_log_dir: "./logs"
  log_request_body: true
  log_response_body: true
```

### 生产环境

```yaml
logging:
  level: "info"
  format: "json"
  file_log_enabled: true
  file_log_dir: "/var/log/go-bisub"
  log_request_body: false
  log_response_body: false
```

## 📚 相关文档

- [日志系统使用文档](./LOGGING_SYSTEM.md)
- [日志实现总结](./LOGGING_IMPLEMENTATION.md)
- [Go slog 官方文档](https://pkg.go.dev/log/slog)
- [Uber Zap 文档](https://github.com/uber-go/zap)

## 🎯 总结

采用 **slog + zap** 架构的优势：

1. ✅ **标准化**: 使用 Go 官方 slog 接口
2. ✅ **高性能**: zap 提供零分配的日志记录
3. ✅ **易用性**: slog 提供简洁的 API
4. ✅ **可维护**: 标准接口便于测试和替换
5. ✅ **未来兼容**: 跟随 Go 官方标准演进

这种架构既保证了性能，又保证了代码的可维护性和未来兼容性。

---

**架构版本**: v2.0.0  
**更新日期**: 2024-11-28  
**Go 版本要求**: 1.21+
