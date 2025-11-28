# 日志和指标使用示例

## 📝 日志使用示例

### 基础日志记录

```go
package main

import (
    "git.uhomes.net/uhs-go/go-bisub/internal/pkg/logger"
)

func main() {
    // INFO 级别日志
    logger.Info("服务启动",
        "port", 8080,
        "environment", "production",
    )
    
    // DEBUG 级别日志
    logger.Debug("配置加载完成",
        "config_file", "./config.yaml",
        "items", 10,
    )
    
    // WARN 级别日志
    logger.Warn("缓存未命中",
        "key", "user:123",
        "fallback", "database",
    )
    
    // ERROR 级别日志
    logger.Error("数据库连接失败",
        "database", "primary",
        "error", "connection timeout",
        "retry_count", 3,
    )
}
```

### 带上下文的日志

```go
func HandleRequest(ctx context.Context) {
    // 从 context 中自动提取 trace_id, span_id, request_id
    logger.InfoContext(ctx, "处理用户请求",
        "user_id", "usr_123",
        "action", "create_subscription",
    )
    
    // 创建带字段的 logger
    log := logger.WithContext(ctx).WithFields(map[string]interface{}{
        "service": "subscription-service",
        "version": "1.0.0",
    })
    
    log.Info("开始处理")
    // ... 业务逻辑
    log.Info("处理完成", "duration_ms", 150)
}
```

### 业务日志示例

```go
// 用户操作日志
func CreateUser(ctx context.Context, user *User) error {
    logger.InfoContext(ctx, "用户创建开始",
        "user_id", user.ID,
        "username", user.Username,
        "operator", getOperator(ctx),
    )
    
    if err := db.Create(user).Error; err != nil {
        logger.ErrorContext(ctx, "用户创建失败",
            "user_id", user.ID,
            "error", err.Error(),
        )
        return err
    }
    
    logger.InfoContext(ctx, "用户创建成功",
        "user_id", user.ID,
        "duration_ms", 50,
    )
    return nil
}

// 订阅执行日志
func ExecuteSubscription(ctx context.Context, key string) error {
    start := time.Now()
    
    logger.InfoContext(ctx, "订阅执行开始",
        "subscription_key", key,
    )
    
    result, err := executeSQL(ctx, key)
    duration := time.Since(start)
    
    if err != nil {
        logger.ErrorContext(ctx, "订阅执行失败",
            "subscription_key", key,
            "duration_ms", duration.Milliseconds(),
            "error", err.Error(),
        )
        return err
    }
    
    logger.InfoContext(ctx, "订阅执行成功",
        "subscription_key", key,
        "duration_ms", duration.Milliseconds(),
        "rows_affected", result.RowsAffected,
    )
    return nil
}

// 数据库操作日志
func QueryDatabase(ctx context.Context, sql string) error {
    start := time.Now()
    
    logger.DebugContext(ctx, "执行SQL查询",
        "sql", sql,
    )
    
    err := db.Raw(sql).Error
    duration := time.Since(start)
    
    if duration > 200*time.Millisecond {
        logger.WarnContext(ctx, "慢查询检测",
            "sql", sql,
            "duration_ms", duration.Milliseconds(),
            "threshold_ms", 200,
        )
    }
    
    if err != nil {
        logger.ErrorContext(ctx, "SQL查询失败",
            "sql", sql,
            "duration_ms", duration.Milliseconds(),
            "error", err.Error(),
        )
        return err
    }
    
    return nil
}
```

## 📊 指标使用示例

### HTTP 请求指标

```go
import (
    "time"
    "git.uhomes.net/uhs-go/go-bisub/internal/pkg/metrics"
)

func HandleHTTPRequest(c *gin.Context) {
    start := time.Now()
    
    // 处理请求
    c.Next()
    
    // 记录指标
    duration := time.Since(start)
    metrics.RecordHTTPRequest(
        "go-bisub",
        c.Request.Method,
        c.FullPath(),
        c.Writer.Status(),
        duration,
        c.Request.ContentLength,
        int64(c.Writer.Size()),
    )
}
```

### 数据库指标

```go
func QueryWithMetrics(ctx context.Context, sql string) error {
    start := time.Now()
    
    err := db.WithContext(ctx).Raw(sql).Error
    duration := time.Since(start)
    
    // 记录指标
    metrics.RecordDBQuery(
        "go-bisub",
        "primary",
        "SELECT",
        duration,
        err,
    )
    
    return err
}

// 更新连接池指标
func UpdateConnectionPoolMetrics() {
    sqlDB, _ := db.DB()
    stats := sqlDB.Stats()
    
    metrics.SetDBConnections("go-bisub", "primary", "idle", stats.Idle)
    metrics.SetDBConnections("go-bisub", "primary", "active", stats.InUse)
    metrics.SetDBConnections("go-bisub", "primary", "total", stats.OpenConnections)
}
```

### 业务指标

```go
// 订阅执行指标
func ExecuteWithMetrics(ctx context.Context, key string) error {
    start := time.Now()
    
    err := execute(ctx, key)
    duration := time.Since(start)
    
    // 记录指标
    metrics.RecordExecution("go-bisub", key, duration, err)
    
    return err
}

// 错误指标
func HandleError(err error, errorType string) {
    metrics.RecordError("go-bisub", errorType, getErrorCode(err))
}
```

### 系统指标

```go
import (
    "runtime"
    "git.uhomes.net/uhs-go/go-bisub/internal/pkg/metrics"
)

// 定期更新系统指标
func UpdateSystemMetrics() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        // CPU 使用率
        cpuPercent := getCPUUsage()
        metrics.SetCPUUsage("go-bisub", cpuPercent)
        
        // 内存使用
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        metrics.SetMemoryUsage("go-bisub", "used", int64(m.Alloc))
        metrics.SetMemoryUsage("go-bisub", "total", int64(m.Sys))
        
        // 磁盘使用率
        diskPercent := getDiskUsage("/")
        metrics.SetDiskUsage("go-bisub", "/", diskPercent)
    }
}
```

## 🔍 查询示例

### Prometheus 查询

```promql
# 查询 QPS
rate(http_requests_total{service="go-bisub"}[5m])

# 查询错误率
sum(rate(http_requests_total{service="go-bisub",status=~"5.."}[5m]))
/
sum(rate(http_requests_total{service="go-bisub"}[5m]))

# 查询 P95 延迟
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket{service="go-bisub"}[5m])) by (le)
)

# 查询慢查询数量
rate(db_slow_queries_total{service="go-bisub"}[5m])

# 查询数据库连接池使用率
(
  db_connections{service="go-bisub",state="active"}
  /
  db_connections{service="go-bisub",state="total"}
) * 100
```

### 日志查询

```bash
# 查询错误日志
cat logs/241128.log | jq 'select(.level == "ERROR")'

# 查询特定用户的操作
cat logs/241128.log | jq 'select(.user_id == "usr_123")'

# 查询慢操作（>1秒）
cat logs/241128.log | jq 'select(.duration_ms > 1000)'

# 追踪特定请求
REQUEST_ID="req-xxx"
cat logs/*.log | jq "select(.request_id == \"$REQUEST_ID\")"

# 统计错误类型
cat logs/241128.log | jq -r 'select(.level == "ERROR") | .error' | sort | uniq -c
```

## 🎯 完整示例

### 订阅服务示例

```go
package service

import (
    "context"
    "time"
    
    "git.uhomes.net/uhs-go/go-bisub/internal/pkg/logger"
    "git.uhomes.net/uhs-go/go-bisub/internal/pkg/metrics"
)

type SubscriptionService struct {
    repo *SubscriptionRepository
}

func (s *SubscriptionService) Execute(ctx context.Context, key string, params map[string]interface{}) error {
    start := time.Now()
    
    // 记录开始日志
    logger.InfoContext(ctx, "订阅执行开始",
        "subscription_key", key,
        "params", params,
    )
    
    // 获取订阅配置
    sub, err := s.repo.GetByKey(ctx, key)
    if err != nil {
        logger.ErrorContext(ctx, "获取订阅配置失败",
            "subscription_key", key,
            "error", err.Error(),
        )
        metrics.RecordError("go-bisub", "subscription_not_found", "404")
        return err
    }
    
    // 执行 SQL
    sqlStart := time.Now()
    result, err := s.executeSQL(ctx, sub.SQLContent, params)
    sqlDuration := time.Since(sqlStart)
    
    // 记录 SQL 指标
    metrics.RecordDBQuery("go-bisub", "primary", "SELECT", sqlDuration, err)
    
    if err != nil {
        logger.ErrorContext(ctx, "SQL执行失败",
            "subscription_key", key,
            "sql", sub.SQLContent,
            "duration_ms", sqlDuration.Milliseconds(),
            "error", err.Error(),
        )
        metrics.RecordExecution("go-bisub", key, time.Since(start), err)
        return err
    }
    
    // 记录成功日志
    duration := time.Since(start)
    logger.InfoContext(ctx, "订阅执行成功",
        "subscription_key", key,
        "duration_ms", duration.Milliseconds(),
        "rows_affected", result.RowsAffected,
    )
    
    // 记录业务指标
    metrics.RecordExecution("go-bisub", key, duration, nil)
    
    return nil
}
```

### HTTP Handler 示例

```go
package handler

import (
    "time"
    
    "git.uhomes.net/uhs-go/go-bisub/internal/pkg/logger"
    "git.uhomes.net/uhs-go/go-bisub/internal/pkg/metrics"
    "github.com/gin-gonic/gin"
)

func (h *Handler) CreateSubscription(c *gin.Context) {
    start := time.Now()
    ctx := c.Request.Context()
    
    var req CreateSubscriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.ErrorContext(ctx, "请求参数解析失败",
            "error", err.Error(),
        )
        metrics.RecordError("go-bisub", "invalid_request", "400")
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    logger.InfoContext(ctx, "创建订阅请求",
        "subscription_key", req.SubKey,
        "version", req.Version,
    )
    
    sub, err := h.service.Create(ctx, &req)
    if err != nil {
        logger.ErrorContext(ctx, "创建订阅失败",
            "subscription_key", req.SubKey,
            "error", err.Error(),
        )
        metrics.RecordError("go-bisub", "create_failed", "500")
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    duration := time.Since(start)
    logger.InfoContext(ctx, "创建订阅成功",
        "subscription_key", req.SubKey,
        "version", req.Version,
        "duration_ms", duration.Milliseconds(),
    )
    
    c.JSON(200, gin.H{"data": sub})
}
```

## 📚 相关文档

- [监控标准规范](./MONITORING_STANDARDS.md)
- [日志系统架构](./LOGGING_ARCHITECTURE.md)
- [Prometheus 告警规则](./prometheus-alerts.yml)

---

**版本**: v1.0.0  
**更新日期**: 2024-11-28
