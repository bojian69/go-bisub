# 日志和监控标准规范

## 📋 概述

本文档定义了 go-bisub 项目的日志和监控标准，遵循 Uhomes 微服务设计规范。

## 🎯 设计原则

1. **结构化日志**: 使用 JSON 格式，便于解析和查询
2. **标准化字段**: 统一的字段命名和格式
3. **可追踪性**: 通过 trace_id 和 span_id 追踪请求链路
4. **高性能**: 使用 slog + zap 组合，零内存分配
5. **可观测性**: 完整的指标采集和告警规则

## 📝 结构化日志

### 9.1 日志格式

#### 使用方式

```go
import "git.uhomes.net/uhs-go/go-bisub/internal/pkg/logger"

// 基础日志
logger.Info("用户创建成功",
    "user_id", "usr_123",
    "operator", "admin@uhomes.com",
    "duration", 150,
)

// 带上下文的日志
logger.InfoContext(ctx, "订阅执行完成",
    "subscription_key", "daily_report",
    "rows_affected", 1000,
    "duration_ms", 2500,
)

// 带字段的日志记录器
log := logger.WithFields(map[string]interface{}{
    "service": "go-bisub",
    "version": "1.0.0",
})
log.Info("服务启动")
```

#### 输出格式

```json
{
  "timestamp": "2024-11-28T10:30:00Z",
  "level": "INFO",
  "service": "go-bisub",
  "trace_id": "trace-xxx",
  "span_id": "span-xxx",
  "request_id": "req-xxx",
  "message": "用户创建成功",
  "user_id": "usr_123",
  "operator": "admin@uhomes.com",
  "duration": 150
}
```

### 9.2 日志级别

| 级别 | 用途 | 示例 |
|------|------|------|
| DEBUG | 详细调试信息 | 变量值、函数调用、缓存命中 |
| INFO | 一般信息消息 | 请求处理、业务操作、状态变更 |
| WARN | 警告消息 | 性能下降、废弃使用、配置问题 |
| ERROR | 错误消息 | 已处理的错误、业务异常 |
| FATAL | 严重错误 | 服务崩溃、无法恢复的错误 |

#### 使用示例

```go
// DEBUG: 详细调试信息
logger.Debug("缓存查询",
    "key", cacheKey,
    "hit", true,
    "ttl", 3600,
)

// INFO: 一般信息
logger.Info("订阅创建成功",
    "subscription_key", "daily_report",
    "version", 1,
    "creator", "admin",
)

// WARN: 警告信息
logger.Warn("慢查询检测",
    "sql", sqlQuery,
    "duration_ms", 1500,
    "threshold_ms", 200,
)

// ERROR: 错误信息
logger.Error("数据库连接失败",
    "database", "primary",
    "error", err.Error(),
    "retry_count", 3,
)

// FATAL: 严重错误
logger.Fatal("配置文件加载失败",
    "config_path", "./config.yaml",
    "error", err.Error(),
)
```

### 9.3 标准字段

#### 必需字段

| 字段 | 类型 | 说明 |
|------|------|------|
| timestamp | string | ISO 8601 格式时间戳 |
| level | string | 日志级别 |
| service | string | 服务名称 |
| message | string | 日志消息 |

#### 推荐字段

| 字段 | 类型 | 说明 |
|------|------|------|
| trace_id | string | 分布式追踪 ID |
| span_id | string | Span ID |
| request_id | string | 请求唯一标识 |
| user_id | string | 用户 ID |
| operator | string | 操作人 |
| duration | int | 操作耗时（毫秒） |
| error | string | 错误信息 |

### 9.4 上下文传递

```go
import (
    "context"
    "git.uhomes.net/uhs-go/go-bisub/internal/pkg/logger"
)

// 设置 trace_id 到 context
ctx = context.WithValue(ctx, "trace_id", "trace-xxx")
ctx = context.WithValue(ctx, "span_id", "span-xxx")
ctx = context.WithValue(ctx, "request_id", "req-xxx")

// 使用带上下文的日志
logger.InfoContext(ctx, "处理请求",
    "endpoint", "/api/subscriptions",
    "method", "POST",
)
// 自动包含 trace_id, span_id, request_id
```

## 📊 指标采集

### 9.5 核心指标

#### HTTP 请求指标

```go
import "git.uhomes.net/uhs-go/go-bisub/internal/pkg/metrics"

// 自动记录（通过中间件）
// 或手动记录
metrics.RecordHTTPRequest(
    "go-bisub",           // service
    "POST",               // method
    "/api/subscriptions", // endpoint
    200,                  // status
    150*time.Millisecond, // duration
    1024,                 // request_size
    2048,                 // response_size
)
```

#### 数据库指标

```go
// 记录数据库查询
metrics.RecordDBQuery(
    "go-bisub",           // service
    "primary",            // database
    "SELECT",             // operation
    50*time.Millisecond,  // duration
    nil,                  // error
)

// 设置连接池状态
metrics.SetDBConnections("go-bisub", "primary", "idle", 10)
metrics.SetDBConnections("go-bisub", "primary", "active", 5)
metrics.SetDBConnections("go-bisub", "primary", "total", 15)
```

#### 业务指标

```go
// 记录订阅执行
metrics.RecordExecution(
    "go-bisub",
    "daily_report",
    2500*time.Millisecond,
    nil,
)

// 记录错误
metrics.RecordError("go-bisub", "database_error", "connection_timeout")
```

### 9.6 Prometheus 格式

#### HTTP 请求指标

```prometheus
# 请求持续时间直方图
http_request_duration_seconds{
  service="go-bisub",
  method="POST",
  endpoint="/api/subscriptions",
  status="2xx"
}

# 请求计数器
http_requests_total{
  service="go-bisub",
  method="POST",
  endpoint="/api/subscriptions",
  status="2xx"
}

# 活跃连接数
http_active_connections{service="go-bisub"}
```

#### 数据库指标

```prometheus
# 查询持续时间
db_query_duration_seconds{
  service="go-bisub",
  database="primary",
  operation="SELECT"
}

# 查询计数
db_queries_total{
  service="go-bisub",
  database="primary",
  operation="SELECT",
  status="success"
}

# 慢查询计数
db_slow_queries_total{
  service="go-bisub",
  database="primary"
}

# 连接池状态
db_connections{
  service="go-bisub",
  database="primary",
  state="idle"
}
```

#### 系统指标

```prometheus
# CPU 使用率
system_cpu_usage_percent{service="go-bisub"}

# 内存使用
system_memory_usage_bytes{
  service="go-bisub",
  type="used"
}

# 磁盘使用率
system_disk_usage_percent{
  service="go-bisub",
  mount="/"
}
```

### 9.7 关键指标（RED Method）

#### Rate（请求速率）

```promql
# 每秒请求数
rate(http_requests_total{service="go-bisub"}[5m])

# 按端点分组
sum(rate(http_requests_total{service="go-bisub"}[5m])) by (endpoint)
```

#### Errors（错误率）

```promql
# 错误率
sum(rate(http_requests_total{service="go-bisub",status=~"5.."}[5m]))
/
sum(rate(http_requests_total{service="go-bisub"}[5m]))

# 错误数
sum(rate(http_requests_total{service="go-bisub",status=~"5.."}[5m]))
```

#### Duration（延迟）

```promql
# P50 延迟
histogram_quantile(0.50,
  sum(rate(http_request_duration_seconds_bucket{service="go-bisub"}[5m])) by (le)
)

# P95 延迟
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket{service="go-bisub"}[5m])) by (le)
)

# P99 延迟
histogram_quantile(0.99,
  sum(rate(http_request_duration_seconds_bucket{service="go-bisub"}[5m])) by (le)
)
```

#### Saturation（饱和度）

```promql
# CPU 使用率
system_cpu_usage_percent{service="go-bisub"}

# 内存使用率
(
  system_memory_usage_bytes{service="go-bisub",type="used"}
  /
  system_memory_usage_bytes{service="go-bisub",type="total"}
) * 100

# 磁盘使用率
system_disk_usage_percent{service="go-bisub"}

# 数据库连接池使用率
(
  db_connections{service="go-bisub",state="active"}
  /
  db_connections{service="go-bisub",state="total"}
) * 100
```

## 🚨 告警规则

### 9.8 严重告警（Critical）

#### 错误率 > 5% 持续 5 分钟

```yaml
alert: HighErrorRate
expr: |
  (
    sum(rate(http_requests_total{service="go-bisub",status=~"5.."}[5m]))
    /
    sum(rate(http_requests_total{service="go-bisub"}[5m]))
  ) > 0.05
for: 5m
severity: critical
```

#### P99 延迟 > 1s 持续 5 分钟

```yaml
alert: HighP99Latency
expr: |
  histogram_quantile(0.99,
    sum(rate(http_request_duration_seconds_bucket{service="go-bisub"}[5m])) by (le)
  ) > 1
for: 5m
severity: critical
```

#### 服务可用性 < 99.9%

```yaml
alert: LowServiceAvailability
expr: |
  (
    sum(rate(http_requests_total{service="go-bisub",status!~"5.."}[5m]))
    /
    sum(rate(http_requests_total{service="go-bisub"}[5m]))
  ) < 0.999
for: 5m
severity: critical
```

#### 数据库连接池耗尽

```yaml
alert: DatabaseConnectionPoolExhausted
expr: |
  db_connections{service="go-bisub",state="idle"} < 2
for: 2m
severity: critical
```

### 9.9 警告告警（Warning）

#### 错误率 > 1% 持续 10 分钟

```yaml
alert: ElevatedErrorRate
expr: |
  (
    sum(rate(http_requests_total{service="go-bisub",status=~"5.."}[10m]))
    /
    sum(rate(http_requests_total{service="go-bisub"}[10m]))
  ) > 0.01
for: 10m
severity: warning
```

#### P95 延迟 > 500ms 持续 10 分钟

```yaml
alert: ElevatedP95Latency
expr: |
  histogram_quantile(0.95,
    sum(rate(http_request_duration_seconds_bucket{service="go-bisub"}[10m])) by (le)
  ) > 0.5
for: 10m
severity: warning
```

#### 内存使用 > 80%

```yaml
alert: HighMemoryUsage
expr: |
  (
    system_memory_usage_bytes{service="go-bisub",type="used"}
    /
    system_memory_usage_bytes{service="go-bisub",type="total"}
  ) > 0.8
for: 5m
severity: warning
```

#### 磁盘使用 > 85%

```yaml
alert: HighDiskUsage
expr: |
  system_disk_usage_percent{service="go-bisub"} > 85
for: 5m
severity: warning
```

## 🔧 集成配置

### Prometheus 配置

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'go-bisub'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

### Grafana Dashboard

推荐使用以下面板：

1. **请求概览**: QPS、错误率、延迟
2. **数据库监控**: 查询延迟、慢查询、连接池
3. **系统资源**: CPU、内存、磁盘
4. **业务指标**: 订阅执行、错误分布

## 📚 相关文档

- [日志系统架构](./LOGGING_ARCHITECTURE.md)
- [日志系统使用](./LOGGING_SYSTEM.md)
- [Prometheus 告警规则](./prometheus-alerts.yml)

## ✅ 检查清单

### 日志规范

- [ ] 使用结构化日志格式
- [ ] 包含必需字段（timestamp, level, service, message）
- [ ] 使用正确的日志级别
- [ ] 传递 trace_id 和 request_id
- [ ] 记录关键业务操作
- [ ] 错误日志包含详细信息

### 指标规范

- [ ] 记录 HTTP 请求指标
- [ ] 记录数据库查询指标
- [ ] 记录业务指标
- [ ] 设置合理的 bucket 范围
- [ ] 使用标准的标签名称

### 告警规范

- [ ] 配置严重告警规则
- [ ] 配置警告告警规则
- [ ] 设置合理的阈值
- [ ] 包含详细的告警描述
- [ ] 测试告警规则

---

**版本**: v1.0.0  
**更新日期**: 2024-11-28  
**遵循规范**: Uhomes 微服务设计规范 v2.0
