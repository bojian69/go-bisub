# 日志系统文档

## 概述

项目实现了完整的文件日志系统，包括API请求日志和SQL执行日志。所有日志按日期自动分割，格式为JSON，便于后续分析和监控。

## 日志文件格式

### 文件命名规则

- **API日志**: `YYMMDD.log` (例如: `251128.log`)
- **SQL日志**: `YYMMDD_sql.log` (例如: `251128_sql.log`)

### 日志目录

默认日志目录: `./logs/`

可通过配置文件修改:
```yaml
logging:
  file_log_dir: "./logs"
```

## API日志格式

每行一个JSON对象，包含以下字段：

```json
{
  "timestamp": "2024-11-28 16:30:45.123",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/subscriptions",
  "client_ip": "127.0.0.1",
  "user_agent": "Mozilla/5.0...",
  "status_code": 200,
  "duration_ms": 125,
  "request_body": {
    "sub_key": "test_subscription",
    "version": 1,
    "title": "测试订阅"
  },
  "response_body": {
    "code": "OK",
    "message": "订阅创建成功",
    "data": {...}
  },
  "error_message": ""
}
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| timestamp | string | 请求时间，格式: YYYY-MM-DD HH:mm:ss.SSS |
| request_id | string | 请求唯一标识符（UUID） |
| method | string | HTTP方法（GET/POST/PUT/DELETE等） |
| path | string | 请求路径 |
| client_ip | string | 客户端IP地址 |
| user_agent | string | 客户端User-Agent |
| status_code | int | HTTP状态码 |
| duration_ms | int64 | 请求处理耗时（毫秒） |
| request_body | object | 请求体（JSON格式） |
| response_body | object | 响应体（JSON格式） |
| error_message | string | 错误信息（如果有） |

## SQL日志格式

每行一个JSON对象，包含以下字段：

```json
{
  "timestamp": "2024-11-28 16:30:45.234",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "sql": "SELECT * FROM sub_subscription_theme WHERE sub_key = ? AND version = ?",
  "duration_ms": 15,
  "rows_affected": 1,
  "error": "",
  "variables": {
    "sub_key": "test_subscription",
    "version": 1
  },
  "database": "go_sub"
}
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| timestamp | string | SQL执行时间 |
| request_id | string | 关联的API请求ID |
| sql | string | 执行的SQL语句 |
| duration_ms | int64 | SQL执行耗时（毫秒） |
| rows_affected | int64 | 影响的行数 |
| error | string | 错误信息（如果有） |
| variables | object | SQL参数（如果有） |
| database | string | 数据库名称 |

### 慢查询标记

当SQL执行时间超过阈值（默认200ms）时，SQL字段会添加 `[SLOW QUERY]` 前缀：

```json
{
  "sql": "[SLOW QUERY] SELECT * FROM large_table WHERE ...",
  "duration_ms": 1250
}
```

## 配置说明

### config.yaml 配置

```yaml
logging:
  level: "debug"                    # 日志级别: debug/info/warn/error
  format: "json"                    # 日志格式: json/text
  file_log_enabled: true            # 是否启用文件日志
  file_log_dir: "./logs"            # 日志文件目录
  log_request_body: true            # 是否记录请求体
  log_response_body: true           # 是否记录响应体
```

### 配置项说明

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| file_log_enabled | bool | false | 是否启用文件日志 |
| file_log_dir | string | "./logs" | 日志文件存储目录 |
| log_request_body | bool | false | 是否记录API请求体 |
| log_response_body | bool | false | 是否记录API响应体 |

**注意**: 
- 记录请求/响应体会增加日志文件大小
- 生产环境建议关闭或选择性开启
- 敏感信息（如密码）应在记录前脱敏

## 日志轮转

### 自动轮转

- 日志文件按日期自动轮转
- 每天00:00自动创建新的日志文件
- 旧日志文件自动保留

### 手动清理

建议定期清理旧日志文件：

```bash
# 删除30天前的日志
find ./logs -name "*.log" -mtime +30 -delete

# 压缩7天前的日志
find ./logs -name "*.log" -mtime +7 -exec gzip {} \;
```

### 使用logrotate（Linux）

创建 `/etc/logrotate.d/go-bisub`:

```
/path/to/go-bisub/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0644 app app
    sharedscripts
    postrotate
        # 可选：重启应用或发送信号
    endscript
}
```

## 日志分析

### 使用jq分析日志

```bash
# 查看所有错误请求
cat 251128.log | jq 'select(.status_code >= 400)'

# 统计API调用次数
cat 251128.log | jq -r '.path' | sort | uniq -c | sort -rn

# 查找慢请求（>1秒）
cat 251128.log | jq 'select(.duration_ms > 1000)'

# 查看特定request_id的所有日志
REQUEST_ID="550e8400-e29b-41d4-a716-446655440000"
cat 251128.log | jq "select(.request_id == \"$REQUEST_ID\")"
cat 251128_sql.log | jq "select(.request_id == \"$REQUEST_ID\")"

# 统计SQL执行时间分布
cat 251128_sql.log | jq '.duration_ms' | awk '{
    if ($1 < 10) fast++
    else if ($1 < 100) normal++
    else if ($1 < 1000) slow++
    else very_slow++
}
END {
    print "Fast (<10ms):", fast
    print "Normal (10-100ms):", normal
    print "Slow (100-1000ms):", slow
    print "Very Slow (>1000ms):", very_slow
}'

# 查找慢查询
cat 251128_sql.log | jq 'select(.sql | contains("[SLOW QUERY]"))'
```

### 使用ELK Stack

1. **Filebeat配置** (`filebeat.yml`):

```yaml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /path/to/logs/*.log
  json.keys_under_root: true
  json.add_error_key: true
  fields:
    log_type: api
  fields_under_root: true

- type: log
  enabled: true
  paths:
    - /path/to/logs/*_sql.log
  json.keys_under_root: true
  json.add_error_key: true
  fields:
    log_type: sql
  fields_under_root: true

output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "go-bisub-%{+yyyy.MM.dd}"
```

2. **Kibana可视化**:
   - 创建索引模式: `go-bisub-*`
   - 按 `log_type` 字段区分API和SQL日志
   - 创建Dashboard监控关键指标

### 使用Grafana Loki

1. **Promtail配置** (`promtail.yml`):

```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://localhost:3100/loki/api/v1/push

scrape_configs:
  - job_name: go-bisub-api
    static_configs:
      - targets:
          - localhost
        labels:
          job: go-bisub
          log_type: api
          __path__: /path/to/logs/*.log

  - job_name: go-bisub-sql
    static_configs:
      - targets:
          - localhost
        labels:
          job: go-bisub
          log_type: sql
          __path__: /path/to/logs/*_sql.log
```

## 性能考虑

### 异步写入

日志写入采用异步方式，不会阻塞主请求：

```go
go func() {
    fileLogger := logger.GetFileLogger()
    _ = fileLogger.LogAPI(entry)
}()
```

### 文件缓冲

日志文件使用操作系统的文件缓冲，提高写入性能。

### 日志级别

生产环境建议：
- 关闭 `log_request_body` 和 `log_response_body`
- 使用 `SimpleLoggerMiddleware` 代替 `LoggerMiddleware`
- 设置合理的慢查询阈值

## 故障排查

### 日志文件未生成

1. 检查配置:
```bash
grep -A 5 "logging:" config.yaml
```

2. 检查目录权限:
```bash
ls -ld ./logs
```

3. 检查应用日志:
```bash
# 查看启动日志
./server 2>&1 | grep -i "logger"
```

### 日志文件过大

1. 启用日志轮转
2. 关闭请求/响应体记录
3. 定期清理旧日志

### 性能影响

如果日志记录影响性能：
1. 使用 `SimpleLoggerMiddleware`
2. 关闭SQL日志（不推荐）
3. 使用SSD存储日志
4. 增加文件系统缓存

## 最佳实践

### 1. 生产环境配置

```yaml
logging:
  level: "info"
  format: "json"
  file_log_enabled: true
  file_log_dir: "/var/log/go-bisub"
  log_request_body: false
  log_response_body: false
```

### 2. 开发环境配置

```yaml
logging:
  level: "debug"
  format: "json"
  file_log_enabled: true
  file_log_dir: "./logs"
  log_request_body: true
  log_response_body: true
```

### 3. 日志保留策略

- 保留最近30天的日志
- 压缩7天前的日志
- 归档重要日志到对象存储

### 4. 监控告警

设置以下告警：
- 错误率超过阈值
- 慢查询数量增加
- 日志文件大小异常
- 磁盘空间不足

### 5. 安全考虑

- 不记录敏感信息（密码、token等）
- 限制日志文件访问权限
- 定期审计日志访问
- 加密传输日志到中心化系统

## 示例查询

### 查看今天的API调用统计

```bash
TODAY=$(date +%y%m%d)
cat logs/${TODAY}.log | jq -s '
  group_by(.path) | 
  map({
    path: .[0].path,
    count: length,
    avg_duration: (map(.duration_ms) | add / length),
    max_duration: (map(.duration_ms) | max),
    error_count: (map(select(.status_code >= 400)) | length)
  })
'
```

### 查看今天的慢SQL

```bash
TODAY=$(date +%y%m%d)
cat logs/${TODAY}_sql.log | jq 'select(.duration_ms > 200) | {sql, duration_ms, request_id}'
```

### 追踪特定请求

```bash
REQUEST_ID="your-request-id"
echo "=== API Log ==="
cat logs/*.log | jq "select(.request_id == \"$REQUEST_ID\")"
echo "=== SQL Logs ==="
cat logs/*_sql.log | jq "select(.request_id == \"$REQUEST_ID\")"
```

## 相关文档

- [配置文档](./CONFIGURATION.md)
- [性能优化](./PERFORMANCE.md)
- [监控告警](./MONITORING.md)

## 更新日志

### v1.0.0 (2024-11-28)
- ✅ 实现文件日志系统
- ✅ 支持API和SQL日志分离
- ✅ 自动日志轮转
- ✅ 异步日志写入
- ✅ 慢查询标记
- ✅ Request ID追踪
