# 日志和监控快速开始

## 🚀 5分钟快速上手

### 1. 启动服务

```bash
make dev
```

服务启动后，日志系统自动启用。

### 2. 查看日志

```bash
# 查看今天的 API 日志
./scripts/view_logs.sh

# 查看统计信息
./scripts/view_logs.sh --stats

# 实时跟踪日志
./scripts/view_logs.sh -f
```

### 3. 查看指标

```bash
# 访问 Prometheus 指标端点
curl http://localhost:8080/metrics

# 查看 HTTP 请求指标
curl -s http://localhost:8080/metrics | grep "^http_"

# 查看数据库指标
curl -s http://localhost:8080/metrics | grep "^db_"
```

### 4. 访问 Web UI

打开浏览器访问: http://localhost:8080/admin

- 用户名: `admin`
- 密码: `admin123`

## 📝 代码示例

### 记录日志

```go
import "git.uhomes.net/uhs-go/go-bisub/internal/pkg/logger"

// 基础日志
logger.Info("订阅创建成功",
    "subscription_key", "daily_report",
    "version", 1,
    "creator", "admin",
)

// 带上下文的日志（自动包含 request_id, trace_id）
logger.InfoContext(ctx, "订阅执行完成",
    "subscription_key", "daily_report",
    "duration_ms", 2500,
    "rows_affected", 1000,
)

// 错误日志
logger.Error("数据库连接失败",
    "database", "primary",
    "error", err.Error(),
    "retry_count", 3,
)
```

### 记录指标

指标会通过中间件自动记录，无需手动调用。

如需手动记录业务指标：

```go
import "git.uhomes.net/uhs-go/go-bisub/internal/pkg/metrics"

// 记录订阅执行
metrics.RecordExecution(
    "go-bisub",
    "daily_report",
    2500*time.Millisecond,
    nil,
)

// 记录错误
metrics.RecordError("go-bisub", "database_error", "timeout")
```

## 🔍 日志查询

### 使用脚本查询

```bash
# 查看错误日志
./scripts/view_logs.sh -e

# 查看慢查询
./scripts/view_logs.sh -t sql -s

# 按 request_id 查询
./scripts/view_logs.sh -r <request-id>

# 查看最后 50 行
./scripts/view_logs.sh --tail 50
```

### 使用 jq 查询

```bash
# 查看所有错误请求
cat logs/251128.log | jq 'select(.status_code >= 400)'

# 统计 API 调用次数
cat logs/251128.log | jq -r '.path' | sort | uniq -c | sort -rn

# 查找慢请求（>1秒）
cat logs/251128.log | jq 'select(.duration_ms > 1000)'

# 查找慢 SQL
cat logs/251128_sql.log | jq 'select(.sql | contains("[SLOW QUERY]"))'

# 追踪特定请求
REQUEST_ID="your-request-id"
cat logs/*.log | jq "select(.request_id == \"$REQUEST_ID\")"
```

## 📊 Prometheus 查询

### 基础查询

```promql
# QPS（每秒请求数）
rate(http_requests_total{service="go-bisub"}[5m])

# 错误率
sum(rate(http_requests_total{service="go-bisub",status=~"5.."}[5m]))
/
sum(rate(http_requests_total{service="go-bisub"}[5m]))

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

### 数据库查询

```promql
# 数据库查询 QPS
rate(db_queries_total{service="go-bisub"}[5m])

# 慢查询数量
rate(db_slow_queries_total{service="go-bisub"}[5m])

# 数据库连接池使用率
(
  db_connections{service="go-bisub",state="active"}
  /
  db_connections{service="go-bisub",state="total"}
) * 100
```

### 业务查询

```promql
# 订阅执行 QPS
rate(execution_total{service="go-bisub"}[5m])

# 订阅执行失败率
sum(rate(execution_total{service="go-bisub",status="error"}[5m]))
/
sum(rate(execution_total{service="go-bisub"}[5m]))

# 订阅执行 P95 延迟
histogram_quantile(0.95,
  sum(rate(execution_duration_seconds_bucket{service="go-bisub"}[5m])) by (le, subscription_key)
)
```

## 🎯 常用场景

### 场景1: 排查慢请求

```bash
# 1. 查看慢请求日志
cat logs/251128.log | jq 'select(.duration_ms > 1000)'

# 2. 查看慢 SQL
cat logs/251128_sql.log | jq 'select(.sql | contains("[SLOW QUERY]"))'

# 3. 查看 Prometheus 指标
curl -s http://localhost:8080/metrics | grep "http_request_duration"
```

### 场景2: 追踪特定请求

```bash
# 1. 从 API 响应获取 request_id
curl -i http://localhost:8080/health
# X-Request-Id: b96a441b-3069-4313-b6b2-f06b3aa5918f

# 2. 查询该请求的所有日志
REQUEST_ID="b96a441b-3069-4313-b6b2-f06b3aa5918f"
cat logs/*.log | jq "select(.request_id == \"$REQUEST_ID\")"
cat logs/*_sql.log | jq "select(.request_id == \"$REQUEST_ID\")"
```

### 场景3: 监控错误率

```bash
# 1. 查看错误日志
./scripts/view_logs.sh -e

# 2. 统计错误类型
cat logs/251128.log | jq -r 'select(.status_code >= 400) | .path' | sort | uniq -c

# 3. 查看 Prometheus 错误率
curl -s http://localhost:8080/metrics | grep "http_requests_total.*5.."
```

## 🔧 配置调整

### 开发环境（详细日志）

```yaml
logging:
  level: "debug"
  format: "console"
  file_log_enabled: true
  file_log_dir: "./logs"
  log_request_body: true
  log_response_body: true
```

### 生产环境（精简日志）

```yaml
logging:
  level: "info"
  format: "json"
  file_log_enabled: true
  file_log_dir: "/var/log/go-bisub"
  log_request_body: false
  log_response_body: false
```

## 📚 下一步

1. **集成 Grafana**: 创建可视化 Dashboard
2. **配置告警**: 导入 `docs/prometheus-alerts.yml`
3. **日志聚合**: 集成 ELK 或 Loki
4. **分布式追踪**: 集成 Jaeger 或 Zipkin

## 🆘 故障排查

### 日志文件未生成

```bash
# 检查配置
grep -A 5 "logging:" config.yaml

# 检查目录权限
ls -ld ./logs

# 查看应用日志
tail -f logs/251128.log
```

### 指标端点无响应

```bash
# 检查服务是否运行
curl http://localhost:8080/health

# 检查路由注册
curl http://localhost:8080/metrics | head -10
```

### 性能问题

```bash
# 关闭请求/响应体记录
# 在 config.yaml 中设置:
# log_request_body: false
# log_response_body: false

# 使用简化日志中间件
# 在代码中使用 SimpleLoggerMiddleware
```

---

**快速开始版本**: v1.0.0  
**更新日期**: 2024-11-28  
**适用版本**: Go 1.21+
