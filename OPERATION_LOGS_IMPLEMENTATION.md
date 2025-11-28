# 操作日志功能实现说明

## 概述

为 go-bisub 项目添加了完整的操作日志功能，实现对所有用户操作的审计跟踪。

## 新增文件

### 1. 数据模型
- `internal/models/operation_log.go` - 操作日志数据模型

### 2. 数据访问层
- `internal/repository/operation_log.go` - 操作日志数据库操作

### 3. 业务逻辑层
- `internal/service/operation_log.go` - 操作日志业务逻辑

### 4. 控制器层
- `internal/handler/operation_log.go` - 操作日志HTTP接口

### 5. 中间件
- `internal/middleware/operation_log.go` - 操作日志记录中间件

### 6. 数据库脚本
- `init_operation_logs.sql` - 操作日志表创建脚本

### 7. Web界面
- `web/templates/operation_logs.html` - 操作日志管理界面

## 修改文件

### 1. 订阅处理器
- `internal/handler/subscription.go` - 添加操作日志记录

### 2. 数据库初始化
- `init.sql` - 添加操作日志表

### 3. 项目文档
- `README.md` - 更新功能特性和API文档

## 功能特性

### 1. 操作类型支持
- CREATE - 创建操作
- UPDATE - 更新操作  
- DELETE - 删除操作
- EXECUTE - 执行操作
- QUERY - 查询操作
- LOGIN - 登录操作
- LOGOUT - 登出操作

### 2. 记录内容
- 用户信息（ID、用户名）
- 操作详情（类型、资源、资源ID）
- 请求信息（IP、User-Agent、URL、HTTP方法）
- 执行信息（耗时、状态、错误信息）
- 数据内容（请求数据、响应数据）

### 3. 查询功能
- 时间范围过滤
- 用户过滤
- 操作类型过滤
- 资源类型过滤
- 状态过滤
- IP地址过滤
- 分页查询

### 4. Web界面
- 操作日志列表展示
- 多条件搜索过滤
- 分页浏览
- 操作详情查看

## API接口

### 获取操作日志
```
GET /v1/operation-logs
```

支持查询参数：
- `start_time` - 开始时间 (YYYY-MM-DD)
- `end_time` - 结束时间 (YYYY-MM-DD)
- `user_id` - 用户ID
- `username` - 用户名（模糊匹配）
- `operation` - 操作类型
- `resource` - 资源类型（模糊匹配）
- `status` - 操作状态 (SUCCESS/FAILED)
- `client_ip` - 客户端IP
- `limit` - 每页数量（默认20，最大100）
- `offset` - 偏移量（默认0）

## 数据库表结构

### sub_logs_operation 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT UNSIGNED | 主键ID |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| user_id | BIGINT UNSIGNED | 操作用户ID |
| username | VARCHAR(120) | 操作用户名 |
| operation | VARCHAR(50) | 操作类型 |
| resource | VARCHAR(200) | 操作资源 |
| resource_id | VARCHAR(120) | 资源ID |
| status | VARCHAR(20) | 操作状态 |
| client_ip | VARCHAR(45) | 客户端IP |
| user_agent | VARCHAR(500) | 用户代理 |
| request_url | VARCHAR(1000) | 请求URL |
| method | VARCHAR(10) | HTTP方法 |
| duration | INT UNSIGNED | 执行耗时(毫秒) |
| error_msg | TEXT | 错误信息 |
| request_data | JSON | 请求数据 |
| response_data | JSON | 响应数据 |

## 使用方式

### 1. 自动记录
通过中间件自动记录所有API操作，无需手动调用。

### 2. 手动记录
```go
// 创建操作日志
log := logService.CreateOperationLog(
    userID,
    username,
    operation,
    resource,
    resourceID,
    status,
    clientIP,
    userAgent,
    requestURL,
    method,
    duration,
    errorMsg,
    requestData,
    responseData,
)

// 异步记录
logService.LogOperation(ctx, log)
```

### 3. Web界面访问
访问 `http://localhost:8080/admin/operation-logs` 查看操作日志。

## 安全考虑

1. **数据脱敏** - 敏感数据在记录前进行脱敏处理
2. **异步记录** - 日志记录不影响主业务流程
3. **存储优化** - 大数据量时考虑分表或归档策略
4. **访问控制** - 操作日志查看需要相应权限

## 性能优化

1. **异步处理** - 所有日志记录都是异步执行
2. **索引优化** - 为常用查询字段建立索引
3. **数据清理** - 定期清理过期日志数据
4. **批量处理** - 高并发时可考虑批量写入

## 部署说明

1. 执行数据库迁移脚本 `init_operation_logs.sql`
2. 更新应用配置，启用操作日志中间件
3. 重启应用服务
4. 访问Web界面验证功能

## 监控告警

建议对以下指标进行监控：
- 操作失败率
- 异常操作频率
- 敏感操作审计
- 系统性能影响