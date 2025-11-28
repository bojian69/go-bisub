# 日志系统实现总结

## 📋 实现概述

已为项目添加完整的文件日志系统，实现了API请求日志和SQL执行日志的自动记录。

## ✅ 完成的工作

### 1. 核心日志模块

#### 文件日志记录器 (`internal/pkg/logger/file_logger.go`)
- ✅ 支持API和SQL日志分离
- ✅ 按日期自动轮转（YYMMDD格式）
- ✅ JSON格式输出
- ✅ 异步写入，不阻塞主流程
- ✅ 线程安全
- ✅ 自动创建日志目录

#### GORM日志插件 (`internal/pkg/logger/gorm_logger.go`)
- ✅ 集成GORM日志系统
- ✅ 记录所有SQL执行
- ✅ 慢查询标记（>200ms）
- ✅ Request ID追踪
- ✅ 错误记录

#### Gin中间件 (`internal/middleware/logger.go`)
- ✅ 记录所有API请求
- ✅ 捕获请求/响应体
- ✅ 生成Request ID
- ✅ 计算请求耗时
- ✅ 错误信息记录
- ✅ 提供简化版中间件

### 2. 配置系统

#### 配置结构 (`internal/config/config.go`)
```go
type LoggingConfig struct {
    Level          string `mapstructure:"level"`
    Format         string `mapstructure:"format"`
    FileLogEnabled bool   `mapstructure:"file_log_enabled"`
    FileLogDir     string `mapstructure:"file_log_dir"`
    LogRequestBody bool   `mapstructure:"log_request_body"`
    LogResponseBody bool  `mapstructure:"log_response_body"`
}
```

#### 配置文件 (`config.yaml`)
```yaml
logging:
  level: "debug"
  format: "json"
  file_log_enabled: true
  file_log_dir: "./logs"
  log_request_body: true
  log_response_body: true
```

### 3. 集成到应用

#### FX模块集成 (`internal/pkg/fx/modules.go`)
- ✅ 初始化文件日志系统
- ✅ 配置GORM日志
- ✅ 注册Gin中间件
- ✅ 根据配置启用/禁用功能

### 4. 工具和文档

#### 日志查看脚本 (`scripts/view_logs.sh`)
- ✅ 查看API/SQL日志
- ✅ 实时跟踪
- ✅ 错误过滤
- ✅ 慢查询过滤
- ✅ Request ID追踪
- ✅ 统计信息展示

#### 文档
- ✅ 完整的日志系统文档 (`docs/LOGGING_SYSTEM.md`)
- ✅ 日志目录说明 (`logs/README.md`)
- ✅ README更新

### 5. 测试

#### 单元测试 (`internal/pkg/logger/file_logger_test.go`)
- ✅ 文件日志记录测试
- ✅ API日志测试
- ✅ SQL日志测试
- ✅ 简化接口测试

## 📁 文件结构

```
go-bisub/
├── internal/
│   ├── config/
│   │   └── config.go                    # 添加日志配置
│   ├── middleware/
│   │   └── logger.go                    # 新增：API日志中间件
│   └── pkg/
│       └── logger/
│           ├── file_logger.go           # 新增：文件日志记录器
│           ├── file_logger_test.go      # 新增：单元测试
│           └── gorm_logger.go           # 新增：GORM日志插件
├── logs/
│   ├── README.md                        # 新增：日志目录说明
│   ├── YYMMDD.log                       # 自动生成：API日志
│   └── YYMMDD_sql.log                   # 自动生成：SQL日志
├── scripts/
│   └── view_logs.sh                     # 新增：日志查看工具
├── docs/
│   ├── LOGGING_SYSTEM.md                # 新增：日志系统文档
│   └── LOGGING_IMPLEMENTATION.md        # 新增：实现总结
├── config.yaml                          # 更新：添加日志配置
├── config.local.yaml                    # 更新：添加日志配置
└── README.md                            # 更新：添加日志系统说明
```

## 🎯 日志格式

### API日志示例

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
  "request_body": {...},
  "response_body": {...},
  "error_message": ""
}
```

### SQL日志示例

```json
{
  "timestamp": "2024-11-28 16:30:45.234",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "sql": "SELECT * FROM sub_subscription_theme WHERE sub_key = ? AND version = ?",
  "duration_ms": 15,
  "rows_affected": 1,
  "error": "",
  "variables": {...},
  "database": "go_sub"
}
```

## 🚀 使用方法

### 1. 启动应用

日志系统会自动启动，无需额外配置：

```bash
make dev
```

### 2. 查看日志

```bash
# 查看今天的API日志
./scripts/view_logs.sh

# 查看今天的SQL日志
./scripts/view_logs.sh -t sql

# 实时跟踪
./scripts/view_logs.sh -f

# 查看统计
./scripts/view_logs.sh --stats
```

### 3. 分析日志

```bash
# 查看错误请求
cat logs/251128.log | jq 'select(.status_code >= 400)'

# 查看慢查询
cat logs/251128_sql.log | jq 'select(.sql | contains("[SLOW QUERY]"))'

# 追踪特定请求
REQUEST_ID="your-request-id"
cat logs/*.log | jq "select(.request_id == \"$REQUEST_ID\")"
cat logs/*_sql.log | jq "select(.request_id == \"$REQUEST_ID\")"
```

## 🔧 配置选项

### 开发环境配置

```yaml
logging:
  level: "debug"
  format: "json"
  file_log_enabled: true
  file_log_dir: "./logs"
  log_request_body: true      # 记录请求体
  log_response_body: true     # 记录响应体
```

### 生产环境配置

```yaml
logging:
  level: "info"
  format: "json"
  file_log_enabled: true
  file_log_dir: "/var/log/go-bisub"
  log_request_body: false     # 不记录请求体
  log_response_body: false    # 不记录响应体
```

## 📊 性能影响

### 优化措施

1. **异步写入**: 日志写入不阻塞主请求
2. **文件缓冲**: 利用操作系统文件缓冲
3. **可选记录**: 可关闭请求/响应体记录
4. **简化模式**: 提供简化版中间件

### 性能测试

```bash
# 运行性能测试
go test -bench=. -benchmem ./internal/pkg/logger/...
```

## 🔍 监控和告警

### 建议监控指标

1. **错误率**: 4xx/5xx响应占比
2. **响应时间**: P50/P95/P99
3. **慢查询**: SQL执行时间>200ms
4. **日志文件大小**: 防止磁盘满
5. **请求量**: QPS统计

### 告警规则

```yaml
# Prometheus告警规则示例
groups:
  - name: go-bisub-logs
    rules:
      - alert: HighErrorRate
        expr: rate(api_errors_total[5m]) > 0.1
        annotations:
          summary: "High error rate detected"
      
      - alert: SlowQueries
        expr: rate(sql_slow_queries_total[5m]) > 10
        annotations:
          summary: "Too many slow queries"
```

## 🛠️ 故障排查

### 日志文件未生成

1. 检查配置: `grep -A 5 "logging:" config.yaml`
2. 检查目录权限: `ls -ld ./logs`
3. 查看应用日志: `./server 2>&1 | grep -i "logger"`

### 日志文件过大

1. 启用日志轮转
2. 关闭请求/响应体记录
3. 定期清理旧日志

### 性能影响

1. 使用 `SimpleLoggerMiddleware`
2. 关闭SQL日志（不推荐）
3. 使用SSD存储日志

## 📚 相关文档

- [日志系统详细文档](./LOGGING_SYSTEM.md)
- [配置文档](./CONFIGURATION.md)
- [性能优化](./PERFORMANCE.md)

## ✨ 后续优化

### 可选功能

- [ ] 日志压缩
- [ ] 日志上传到对象存储
- [ ] 集成ELK/Loki
- [ ] 日志采样（高流量场景）
- [ ] 敏感信息脱敏
- [ ] 结构化日志查询API
- [ ] 日志告警集成

### 性能优化

- [ ] 批量写入
- [ ] 内存池优化
- [ ] 零拷贝优化
- [ ] 日志分级存储

## 🎉 总结

日志系统已完全集成到项目中，具备以下特点：

1. ✅ **完整性**: 记录所有API请求和SQL执行
2. ✅ **可追踪**: Request ID贯穿整个请求链路
3. ✅ **易用性**: 提供便捷的查看和分析工具
4. ✅ **高性能**: 异步写入，不影响主流程
5. ✅ **可配置**: 灵活的配置选项
6. ✅ **可扩展**: 易于集成到监控系统

现在你可以：
- 查看所有API请求的详细信息
- 追踪SQL执行和性能
- 分析系统行为和性能瓶颈
- 快速定位和解决问题

---

**实现日期**: 2024-11-28  
**版本**: v1.0.0  
**测试状态**: ✅ 通过
