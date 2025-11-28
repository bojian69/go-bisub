# Uhomes 微服务设计规范

**版本：** 1.0
**最后更新：** 2025-02-19
**参考标准：** AWS API Standards, Alibaba Cloud OpenAPI, Google API Design Guide, CNCF Cloud Native Standards

---

## 1. 核心设计原则

### 1.1 可观察、可演进、可治理

- **可观察**：统一的日志、追踪和指标体系，所有服务可监控
- **可演进**：支持向后兼容的变更和版本演进策略
- **可治理**：集中式服务注册、配置管理和策略执行

### 1.2 简单且可预测

- 固定的请求/响应结构，清晰的契约
- 可预测的行为，无隐藏副作用
- 明确的错误处理和状态报告
- 适用场景下的幂等操作

### 1.3 云原生就绪

- 无状态服务设计，支持水平扩展
- 容器优先的部署策略
- 基础设施即代码（IaC）
- 多云兼容（AWS、阿里云、私有云）

---

## 2. API 端点设计规范

### 2.1 RESTful 资源命名

**规则：**
- 使用小写、复数名词表示资源：`/users`、`/orders`、`/properties`
- 多词资源使用短横线连接：`/rental-applications`、`/payment-methods`
- 路径中避免使用动词（使用 HTTP 方法代替）
- 嵌套资源使用层级结构：`/users/{userId}/orders/{orderId}`

**示例：**
```
GET    /v1/users                    # 获取用户列表
GET    /v1/users/{id}               # 获取用户详情
POST   /v1/users                    # 创建用户
PUT    /v1/users/{id}               # 更新用户（完整）
PATCH  /v1/users/{id}               # 更新用户（部分）
DELETE /v1/users/{id}               # 删除用户
GET    /v1/users/{id}/orders        # 获取用户的订单列表
```

### 2.2 基于动作的操作

对于不适合 CRUD 模型的操作，使用动作后缀：

**云原生风格（推荐）：**
```
POST /v1/instances/{id}:start
POST /v1/instances/{id}:stop
POST /v1/instances/{id}:restart
POST /v1/orders/{id}:cancel
POST /v1/users/{id}:resetPassword
```

**查询参数风格（AWS/阿里云兼容）：**
```
POST /v1/api?Action=StartInstance&InstanceId=xxx
POST /v1/api?Action=StopInstance&InstanceId=xxx
```

### 2.3 版本管理

**版本策略：**
- 使用 URL 路径版本：`/v1`、`/v2`、`/v3`
- 仅使用主版本号（URL 中不包含次版本号）
- 仅在破坏性变更时递增版本
- 至少支持 N-1 版本（当前版本 + 前一版本）

**破坏性变更包括：**
- 删除或重命名字段
- 更改字段类型
- 更改响应结构
- 删除端点
- 更改认证机制

**非破坏性变更：**
- 添加新的可选字段
- 添加新端点
- 添加新的可选查询参数
- 扩展枚举值（带适当默认值）

### 2.4 查询参数

**分页：**
```
GET /v1/users?limit=20&offset=100
```

**过滤：**
```
GET /v1/orders?status=pending&createdAfter=2025-01-01T00:00:00Z
GET /v1/properties?city=NewYork&minPrice=1000&maxPrice=5000
```

**排序：**
```
GET /v1/users?sortBy=createdAt&sortOrder=desc
GET /v1/orders?orderBy=-createdAt,+amount
```

**字段选择：**
```
GET /v1/users?fields=id,name,email
GET /v1/users?expand=profile,orders
```

---

## 3. 请求和响应标准

### 3.1 JSON 字段命名规范（重要）

**统一使用 snake_case 作为 JSON 字段名**

详细规范与验证请参考：[Protobuf_JSON命名规范_v1.md](media/17642089832436/Protobuf_JSON%E5%91%BD%E5%90%8D%E8%A7%84%E8%8C%83_v1.md)

**核心规则：**
- ✅ 使用 `snake_case`：`user_id`、`created_at`、`page_size`
- ❌ 不使用 `camelCase`：~~`userId`~~、~~`createdAt`~~、~~`pageSize`~~
- ❌ 不使用 `PascalCase`：~~`UserId`~~、~~`CreatedAt`~~

**选择 snake_case 的理由（简述）：**
1.  **云原生生态标准**：Kubernetes, Envoy, Prometheus 等均使用 snake_case。
2.  **Protobuf 默认行为**：Protobuf 字段定义和 gRPC-Gateway 默认输出均为 snake_case。
3.  **语言无关性**：对后端语言不敏感，跨语言一致性更强。

**Protobuf 定义示例：**
```protobuf
message User {
  string user_id = 1;           // JSON: "user_id"
  string user_name = 2;         // JSON: "user_name"
  string email_address = 3;     // JSON: "email_address"
}
```

**JSON 示例：**
```json
{
  "user_id": "usr_123",
  "user_name": "张三",
  "email_address": "zhangsan@example.com"
}
```

### 3.2 标准请求结构

**HTTP 头部（必需）：**
```
Content-Type: application/json
X-Request-Id: uuid-v4-format
X-Trace-Id: distributed-trace-id
Authorization: Bearer <token> 
X-Client-Version: app-version
X-Platform: ios|android|web
```

**请求体结构：**
```json
{
  "data": {
    "name": "张三",
    "email": "zhangsan@example.com"
  },
  "metadata": {
    "operator": "admin@uhomes.com",
    "client_ip": "192.168.1.1",
    "idempotency_key": "unique-operation-id"
  }
}
```

### 3.3 标准响应结构

**成功响应：**
```json
{
  "code": "OK",
  "message": "操作成功",
  "request_id": "req-uuid-xxx",
  "data": {
    "id": "user-123",
    "name": "张三",
    "email": "zhangsan@example.com",
    "created_at": "2025-02-19T10:30:00Z"
  },
  "metadata": {
    "timestamp": "2025-02-19T10:30:01Z",
    "version": "v1"
  }
}
```

**列表响应（带分页）：**
```json
{
  "code": "OK",
  "message": "操作成功",
  "request_id": "req-uuid-xxx",
  "data": {
    "items": [...],
    "pagination": {
      "total": 1000,
      "page_size": 20,
      "current_page": 1,
      "total_pages": 50,
      "next_page_token": "token-xxx",
      "has_more": true
    }
  }
}
```

**错误响应：**
```json
{
  "code": "INVALID_PARAMETER",
  "message": "邮箱格式不正确",
  "request_id": "req-uuid-xxx",
  "errors": [
    {
      "field": "email",
      "code": "INVALID_FORMAT",
      "message": "邮箱必须是有效的邮箱地址"
    }
  ],
  "metadata": {
    "timestamp": "2025-02-19T10:30:01Z",
    "documentation": "https://docs.uhomes.net/errors/INVALID_PARAMETER"
  }
}
```

### 3.4 标准错误码

**客户端错误（4xx）：**
- `OK` (200)：成功
- `INVALID_PARAMETER` (400)：请求参数无效
- `UNAUTHORIZED` (401)：认证失败
- `FORBIDDEN` (403)：权限不足
- `NOT_FOUND` (404)：资源不存在
- `CONFLICT` (409)：资源冲突（如重复）
- `PRECONDITION_FAILED` (412)：前置条件不满足
- `RATE_LIMITED` (429)：请求过于频繁
- `IDEMPOTENCY_CONFLICT` (422)：幂等键冲突

**服务端错误（5xx）：**
- `INTERNAL_ERROR` (500)：内部服务器错误
- `SERVICE_UNAVAILABLE` (503)：服务暂时不可用
- `GATEWAY_TIMEOUT` (504)：上游服务超时
- `DEPENDENCY_FAILURE` (502)：下游服务故障

**业务错误（自定义）：**
- `INSUFFICIENT_BALANCE`：账户余额不足
- `ORDER_EXPIRED`：订单已过期
- `PROPERTY_UNAVAILABLE`：房源不可用
- `VERIFICATION_FAILED`：验证码错误

### 3.5 HTTP 状态码映射

```
200 OK              -> 成功的 GET、PUT、PATCH
201 Created         -> 成功的 POST（资源已创建）
202 Accepted        -> 异步操作已接受
204 No Content      -> 成功的 DELETE
400 Bad Request     -> INVALID_PARAMETER
401 Unauthorized    -> UNAUTHORIZED
403 Forbidden       -> FORBIDDEN
404 Not Found       -> NOT_FOUND
409 Conflict        -> CONFLICT
429 Too Many        -> RATE_LIMITED
500 Internal Error  -> INTERNAL_ERROR
503 Unavailable     -> SERVICE_UNAVAILABLE
```

---

## 4. 安全与访问控制

### 4.1 认证方式

**JWT Token（面向用户的 API）：**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```
- Token 过期时间：2 小时（访问令牌），30 天（刷新令牌）
- Claims 中包含用户 ID、角色、权限
- 实现 Token 刷新机制

**Access Key / Secret Key（服务间调用）：**
```
Authorization: AK <AccessKeyId>:<Signature>
X-Timestamp: 1645267800
X-Nonce: random-string
```
- 签名 = HMAC-SHA256(SecretKey, CanonicalRequest)
- 请求在时间戳 5 分钟内有效
- 使用 nonce 防止重放攻击

**API Key（第三方集成）：**
```
X-API-Key: uhomes-api-key-xxx
```
- 按 API Key 限流
- 基于作用域的权限
- 密钥轮换策略（90 天）

### 4.2 授权模型

**RBAC（基于角色的访问控制）：**
```json
{
  "roles": ["admin", "property_manager", "tenant"],
  "permissions": [
    "user:read",
    "user:write",
    "order:read",
    "property:manage"
  ]
}
```

**资源级权限：**
- 操作前检查所有权
- 多租户场景实现租户隔离
- 复杂授权规则使用策略引擎

### 4.3 传输安全

- **仅 HTTPS**：所有 API 必须使用 TLS 1.2+
- **证书固定**：移动应用使用证书固定
- **URL 中不包含敏感数据**：使用请求体或头部
- **加密敏感字段**：PII 数据静态和传输加密

### 4.4 安全头部

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
```

---

## 5. 微服务治理

### 5.1 限流

**策略：**
- 按用户：1000 请求/分钟
- 按 IP：5000 请求/分钟
- 按 Access Key：根据等级自定义配额
- 按端点：关键端点有更低限制

**实现（使用 gox/middleware）：**
```go
import "git.uhomes.net/uhs-go/gox/middleware"

middleware.RateLimit(
    middleware.WithQPS(1000),
    middleware.WithBurst(2000),
    middleware.WithKeyFunc(func(c *gin.Context) string {
        return c.GetHeader("X-User-Id")
    }),
)
```

**响应头部：**
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 950
X-RateLimit-Reset: 1645267800
Retry-After: 60
```

**算法：**
- Token Bucket 用于突发流量
- Sliding Window 用于精确计数
- 使用 Redis 实现分布式限流

### 5.2 熔断与重试

**熔断器状态：**
- Closed（关闭）：正常运行
- Open（打开）：快速失败（达到阈值后）
- Half-Open（半开）：测试恢复

**配置（使用 gox/middleware）：**
```go
middleware.CircuitBreaker(
    middleware.WithFailureThreshold(0.5),      // 50% 失败率
    middleware.WithRequestVolumeThreshold(20), // 最少 20 个请求
    middleware.WithSleepWindow(30*time.Second),// 30 秒后尝试恢复
    middleware.WithTimeout(5*time.Second),     // 5 秒超时
)
```

**重试策略：**
- 仅幂等操作（GET、PUT、DELETE）
- 指数退避：1s、2s、4s、8s
- 最多重试 3 次
- POST 操作使用 `X-Idempotency-Key` 头部

### 5.3 超时配置

**服务级超时（使用 gox/transport）：**
```go
import "git.uhomes.net/uhs-go/gox/transport"

// HTTP 客户端
client := transport.NewHTTPClient(
    transport.WithConnectTimeout(3*time.Second),
    transport.WithReadTimeout(10*time.Second),
    transport.WithWriteTimeout(10*time.Second),
)

// gRPC 客户端
grpcClient := transport.NewGRPCClient(
    transport.WithTimeout(5*time.Second),
)
```

**推荐超时配置：**
```yaml
timeouts:
  http:
    connection: 3s
    read: 10s
    write: 10s
  grpc:
    call: 5s
    stream: 30s
  database:
    query: 5s
    transaction: 30s
  cache:
    get: 100ms
    set: 200ms
```

### 5.4 分布式追踪

**必需头部：**
```
X-Trace-Id: unique-trace-id（跨服务传播）
X-Span-Id: current-span-id
X-Parent-Span-Id: parent-span-id
X-Sampled: 1（采样决策）
```

**追踪集成（使用 gox/trace）：**
```go
import "git.uhomes.net/uhs-go/gox/trace"

// 初始化追踪
tp := trace.Init(
    trace.WithServiceName("user-service"),
    trace.WithEndpoint("http://jaeger:14268/api/traces"),
    trace.WithSampleRate(0.1), // 10% 采样率
)
defer tp.Shutdown(context.Background())

// 自动追踪 HTTP/gRPC 请求
middleware.Tracing()
```

**追踪标准：**
- OpenTelemetry 标准
- Jaeger 后端
- 采样率：正常流量 10%，错误 100%
- 包含服务名、操作、持续时间、状态

### 5.5 健康检查

**存活探针（Liveness Probe）：**
```
GET /health/live
Response: 200 OK
```

**就绪探针（Readiness Probe）：**
```
GET /health/ready
Response: 
{
  "status": "ready",
  "checks": {
    "database": "ok",
    "cache": "ok",
    "downstream": "ok"
  }
}
```

**启动探针（Startup Probe）：**
```
GET /health/startup
Response: 200 OK（初始化完成后）
```

**实现（使用 gox/transport）：**
```go
srv.RegisterHealthCheck(func(ctx context.Context) error {
    // 检查数据库
    if err := db.Ping(ctx); err != nil {
        return err
    }
    // 检查 Redis
    if err := redis.Ping(ctx); err != nil {
        return err
    }
    return nil
})
```

---

## 6. 数据模型规范

### 6.1 字段命名约定

**JSON 字段名（使用 snake_case）：**
- ✅ 使用 `snake_case`：`first_name`、`created_at`、`user_id`
- ❌ 不使用 `camelCase`：~~`firstName`~~、~~`createdAt`~~、~~`userId`~~
- 避免缩写：`description` 而非 `desc`
- 布尔字段：`is_active`、`has_permission`、`can_edit`
- 数组：复数名词：`users`、`orders`、`items`

**数据库列名：**
- 使用下划线命名：`first_name`、`created_at`、`user_id`
- 与 JSON 字段名保持一致（都是 `snake_case`）
- 与 ORM 映射保持一致

**完整示例：**
```json
{
  "user_id": "usr_123",
  "first_name": "张三",
  "last_name": "李",
  "email_address": "zhangsan@example.com",
  "phone_number": "+86-138-0000-0000",
  "is_active": true,
  "is_verified": false,
  "created_at": "2025-02-19T10:30:00Z",
  "updated_at": "2025-02-19T15:45:30Z",
  "last_login_at": "2025-02-19T16:00:00Z"
}
```

### 6.2 时间和日期格式

**ISO 8601 UTC 格式：**
```json
{
  "created_at": "2025-02-19T10:30:00Z",
  "updated_at": "2025-02-19T15:45:30.123Z",
  "scheduled_at": "2025-02-20T08:00:00Z"
}
```

**时区处理：**
- 所有时间戳存储为 UTC
- 客户端转换为本地时区
- 面向用户的日期包含时区信息

**持续时间格式：**
```json
{
  "duration": "PT2H30M",  // ISO 8601 持续时间
  "expires_in": 3600      // 秒（用于 TTL）
}
```

### 6.3 ID 生成

**全局唯一 ID：**
- UUID v4 用于分布式系统
- Snowflake ID 用于带时间戳的有序 ID
- ULID 用于可排序、URL 安全的 ID

**格式：**
```json
{
  "id": "usr_2Nq8xYz9K3mP7wR",           // 带前缀的类型安全 ID
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "snowflake_id": "1234567890123456789"
}
```

**禁止暴露：**
- 自增数据库 ID
- 内部序列号
- 可预测的标识符

### 6.4 枚举和常量

**字符串枚举（推荐）：**
```json
{
  "status": "PENDING",  // PENDING, APPROVED, REJECTED
  "type": "RESIDENTIAL" // RESIDENTIAL, COMMERCIAL, LAND
}
```

**可扩展性：**
- 使用 UNKNOWN 表示未识别的值
- 在 OpenAPI 中记录所有可能的值
- 添加新值不破坏客户端

### 6.5 货币和金额

**十进制精度：**
```json
{
  "amount": "1234.56",
  "currency": "USD",
  "amount_in_cents": 123456
}
```

**规则：**
- 使用字符串或整数（分）表示金额
- 永远不要使用浮点数表示货币
- 始终包含货币代码（ISO 4217）

### 6.6 本地化

**多语言支持：**
```json
{
  "name": {
    "en": "Property Name",
    "zh": "房产名称"
  },
  "description": {
    "en": "Description",
    "zh": "描述"
  }
}
```

**Accept-Language 头部：**
```
Accept-Language: zh-CN,zh;q=0.9,en;q=0.8
```

---

## 7. 错误处理标准

### 7.1 错误响应结构

**详细错误响应：**
```json
{
  "code": "VALIDATION_ERROR",
  "message": "请求验证失败",
  "request_id": "req-uuid-xxx",
  "errors": [
    {
      "field": "email",
      "code": "INVALID_FORMAT",
      "message": "邮箱必须是有效的邮箱地址",
      "value": "invalid-email"
    },
    {
      "field": "age",
      "code": "OUT_OF_RANGE",
      "message": "年龄必须在 18 到 100 之间",
      "value": 150
    }
  ],
  "metadata": {
    "timestamp": "2025-02-19T10:30:01Z",
    "documentation": "https://docs.uhomes.net/errors/VALIDATION_ERROR",
    "supportContact": "support@uhomes.com"
  }
}
```

### 7.2 错误码层次结构

**格式：** `类别_具体错误`

```
AUTH_INVALID_TOKEN
AUTH_TOKEN_EXPIRED
AUTH_INSUFFICIENT_PERMISSION

PAYMENT_INSUFFICIENT_BALANCE
PAYMENT_CARD_DECLINED
PAYMENT_GATEWAY_ERROR

PROPERTY_NOT_AVAILABLE
PROPERTY_ALREADY_BOOKED
PROPERTY_INVALID_DATES
```

**使用 gox/utils 错误封装：**
```go
import "git.uhomes.net/uhs-go/gox/utils/errorx"

// 定义业务错误
var (
    ErrInsufficientBalance = errorx.New("PAYMENT_INSUFFICIENT_BALANCE", "账户余额不足")
    ErrPropertyNotAvailable = errorx.New("PROPERTY_NOT_AVAILABLE", "房源不可用")
)

// 使用
if balance < amount {
    return nil, ErrInsufficientBalance.WithDetails(map[string]interface{}{
        "balance": balance,
        "required": amount,
    })
}
```

### 7.3 客户端错误处理

**重试逻辑：**
- 5xx 错误：使用指数退避重试
- 429 限流：根据 `Retry-After` 头部重试
- 4xx 错误：不重试（除了 408、409 带幂等性）

**用户友好的消息：**
- 提供可操作的错误消息
- 关键错误包含支持联系方式
- 服务端记录详细错误，向用户显示简化消息

---

## 8. OpenAPI（Swagger）规范

### 8.1 OpenAPI 3.0 要求

**每个微服务必须提供：**
- 完整的 OpenAPI 3.0 规范
- 带示例的请求/响应模式
- 错误响应文档
- 认证要求
- 限流信息

**文件位置：**
```
/api/openapi.yaml
/api/v1/openapi.yaml
/docs/api-spec.yaml
```

### 8.2 OpenAPI 示例

```yaml
openapi: 3.0.3
info:
  title: Uhomes 用户服务 API
  version: 1.0.0
  description: 用户管理微服务
  contact:
    email: api@uhomes.com

servers:
  - url: https://api.uhomes.com/v1
    description: 生产环境
  - url: https://api-staging.uhomes.com/v1
    description: 预发布环境

security:
  - bearerAuth: []
  - apiKey: []

paths:
  /users:
    get:
      summary: 获取用户列表
      operationId: listUsers
      tags: [Users]
      parameters:
        - name: pageSize
          in: query
          schema:
            type: integer
            default: 20
            maximum: 100
        - name: pageToken
          in: query
          schema:
            type: string
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserListResponse'
              example:
                code: "OK"
                data:
                  items:
                    - id: "usr_123"
                      name: "张三"
                      email: "zhangsan@example.com"
        '401':
          $ref: '#/components/responses/Unauthorized'
        '429':
          $ref: '#/components/responses/RateLimited'

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
    apiKey:
      type: apiKey
      in: header
      name: X-API-Key

  schemas:
    User:
      type: object
      required: [id, name, email]
      properties:
        id:
          type: string
          example: "usr_123"
        name:
          type: string
          example: "张三"
        email:
          type: string
          format: email
          example: "zhangsan@example.com"
        createdAt:
          type: string
          format: date-time
          example: "2025-02-19T10:30:00Z"

  responses:
    Unauthorized:
      description: 认证失败
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
```

### 8.3 契约测试

**确保 API 实现与 OpenAPI 规范匹配：**
- 使用 Dredd、Schemathesis 或 Prism 等工具
- 在 CI/CD 流程中运行契约测试
- 检测到破坏性变更时阻止部署
- 版本兼容性检查

---

## 9. 日志和监控标准

### 9.1 结构化日志

**日志格式（JSON，使用 gox/log）：**
```go
import "git.uhomes.net/uhs-go/gox/log"

logger := log.New(
    log.WithLevel("info"),
    log.WithFormat("json"),
)

logger.Info("用户创建成功",
    "user_id", "usr_123",
    "operator", "admin@uhomes.com",
    "duration", 150,
)
```

**输出：**
```json
{
  "timestamp": "2025-02-19T10:30:00Z",
  "level": "INFO",
  "service": "user-service",
  "traceId": "trace-xxx",
  "span_id": "span-xxx",
  "message": "用户创建成功",
  "user_id": "usr_123",
  "operator": "admin@uhomes.com",
  "duration": 150
}
```

**日志级别：**
- `DEBUG`：详细调试信息
- `INFO`：一般信息消息
- `WARN`：警告消息（性能下降、废弃使用）
- `ERROR`：错误消息（已处理的错误）
- `FATAL`：严重错误（服务崩溃）

### 9.2 指标采集

**关键指标（使用 gox/metrics）：**
```go
import "git.uhomes.net/uhs-go/gox/metrics"

// 请求计数
metrics.RequestTotal.WithLabelValues("user-service", "GET", "/users", "200").Inc()

// 请求延迟
metrics.RequestDuration.WithLabelValues("user-service", "GET", "/users").Observe(0.15)

// 活跃连接
metrics.ActiveConnections.WithLabelValues("user-service").Set(100)
```

**Prometheus 格式：**
```
# 请求持续时间直方图
http_request_duration_seconds{service="user-service",method="GET",endpoint="/users",status="200"}

# 请求计数器
http_requests_total{service="user-service",method="GET",endpoint="/users",status="200"}

# 活跃连接数
http_active_connections{service="user-service"}
```

**核心指标：**
- 请求速率（requests/second）
- 错误率（errors/total requests）
- 延迟（p50、p95、p99）
- 饱和度（CPU、内存、磁盘、网络）

### 9.3 告警规则

**严重告警：**
- 错误率 > 5% 持续 5 分钟
- P99 延迟 > 1s 持续 5 分钟
- 服务可用性 < 99.9%
- 数据库连接池耗尽

**警告告警：**
- 错误率 > 1% 持续 10 分钟
- P95 延迟 > 500ms 持续 10 分钟
- 内存使用 > 80%
- 磁盘使用 > 85%

---

## 10. 测试标准

### 10.1 测试覆盖率要求

- 单元测试：> 80% 代码覆盖率
- 集成测试：覆盖关键路径
- 契约测试：所有公共 API
- 端到端测试：关键用户旅程

### 10.2 测试命名约定

```go
func TestUserService_CreateUser_Success(t *testing.T)
func TestUserService_CreateUser_DuplicateEmail(t *testing.T)
func TestUserService_CreateUser_InvalidInput(t *testing.T)
```

**模式：** `Test<组件>_<方法>_<场景>`

### 10.3 集成测试最佳实践

- 使用测试容器管理依赖（DB、Redis 等）
- 每个测试后清理测试数据
- 使用 fixtures 保证测试数据一致性
- Mock 外部服务调用

---

## 11. 部署和运维

### 11.1 容器标准

**Dockerfile 最佳实践：**
- 多阶段构建以减小镜像大小
- 非 root 用户运行以提高安全性
- 健康检查指令
- 最小化基础镜像（Alpine、Distroless）

### 11.2 Kubernetes 资源

**必需标签：**
```yaml
metadata:
  labels:
    app: user-service
    version: v1.2.3
    environment: production
    team: backend
```

**资源限制：**
```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

### 11.3 配置管理

**使用 gox/config：**
```go
import "git.uhomes.net/uhs-go/gox/config"

// 从 Nacos 加载配置
cfg := config.MustLoad(
    config.WithNacos(nacosConfig),
    config.WithDataID("user-service"),
    config.WithGroup("DEFAULT_GROUP"),
)

// 监听配置变更
cfg.Watch(func(newCfg *Config) {
    logger.Info("配置已更新", "config", newCfg)
})
```

- 使用环境变量进行运行时配置
- 使用 ConfigMap 存储非敏感配置
- 使用 Secrets 存储敏感数据
- 支持配置热重载

### 11.4 蓝绿部署

- 零停机部署
- 金丝雀发布逐步推出
- 失败时自动回滚
- 功能开关用于 A/B 测试

---

## 12. 文档要求

### 12.1 服务 README

必须包含：
- 服务目的和职责
- 架构概览
- API 文档链接
- 本地开发设置
- 部署说明
- 故障排除指南

### 12.2 API 文档

- OpenAPI 规范（机器可读）
- 交互式 API 浏览器（Swagger UI）
- 多语言代码示例
- 认证指南
- 限流信息
- API 版本变更日志

### 12.3 运维手册

运维指南包括：
- 服务依赖关系
- 健康检查端点
- 常见问题和解决方案
- 扩容指南
- 备份和恢复流程
- 事故响应流程

---

## 13. 合规与治理

### 13.1 数据隐私

- 符合 GDPR（欧盟用户）
- 数据保留策略
- 实现删除权
- 数据导出功能
- 敏感操作的审计日志

### 13.2 API 废弃策略

**废弃流程：**
1. 提前 6 个月宣布废弃
2. 在 OpenAPI 规范中添加 `Deprecated: true`
3. 返回 `Sunset` 头部和结束日期
4. 提供迁移指南
5. 监控使用情况并联系活跃用户
6. 废弃期后删除

**Sunset 头部：**
```
Sunset: Sat, 31 Dec 2025 23:59:59 GMT
Link: <https://docs.uhomes.net/migration/v2>; rel="sunset"
```

---

## 14. 性能优化

### 14.1 缓存策略

**缓存层级：**
- CDN：静态资源、公共 API
- 应用缓存：热数据（Redis）
- 数据库查询缓存：频繁查询
- HTTP 缓存头部：浏览器缓存

**缓存头部：**
```
Cache-Control: public, max-age=3600
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
Last-Modified: Wed, 19 Feb 2025 10:30:00 GMT
```

### 14.2 数据库优化

- 使用连接池
- 读写分离（读副本用于读密集型工作负载）
- 添加适当的索引
- 使用预编译语句
- 实现查询超时
- 监控慢查询

### 14.3 API 响应优化

- 实现字段过滤（`?fields=id,name`）
- 列表端点使用分页
- 压缩响应（gzip、brotli）
- 实现 ETags 用于条件请求
- 复杂数据需求使用 GraphQL

---

## 15. 迁移与向后兼容

### 15.1 数据库迁移

- 使用迁移工具（golang-migrate）
- 版本化所有迁移
- 先在预发布环境测试迁移
- 支持失败迁移的回滚
- 永远不要修改现有迁移

### 15.2 API 向后兼容

**安全变更：**
- 添加新的可选字段
- 添加新端点
- 添加新枚举值（带默认值）
- 放宽验证规则

**破坏性变更（需要新版本）：**
- 删除字段
- 重命名字段
- 更改字段类型
- 将可选字段改为必需
- 更改响应结构

---

## 16. 事故响应

### 16.1 严重级别

- **P0（严重）**：完全服务中断、数据丢失
- **P1（高）**：主要功能损坏、显著用户影响
- **P2（中）**：次要功能损坏、有变通方法
- **P3（低）**：外观问题、影响最小

### 16.2 响应时间 SLA

- P0：15 分钟
- P1：1 小时
- P2：4 小时
- P3：24 小时

### 16.3 事后审查

- 根本原因分析
- 事件时间线
- 防止再次发生的行动项
- 文档更新

---

## 附录 A：参考实现

### AWS API 标准
- https://docs.aws.amazon.com/general/latest/gr/api-conventions.html

### 阿里云 OpenAPI
- https://www.alibabacloud.com/help/zh/openapi

### Google API 设计指南
- https://cloud.google.com/apis/design

### Microsoft REST API 指南
- https://github.com/microsoft/api-guidelines

---

## 附录 B：工具和库

**API 开发：**
- OpenAPI Generator
- Swagger UI / Redoc
- Postman / Insomnia

**测试：**
- Dredd（契约测试）
- Schemathesis（基于属性的测试）
- k6 / Locust（负载测试）

**监控：**
- Prometheus + Grafana
- Jaeger / Zipkin（追踪）
- Loki（日志）

**安全：**
- OWASP ZAP（安全测试）
- Trivy（容器扫描）
- SonarQube（代码质量）

---

## 附录 C：gox 基础库使用

所有微服务应统一使用 gox 基础库提供的能力：

**日志：** `git.uhomes.net/uhs-go/gox/log`  
**配置：** `git.uhomes.net/uhs-go/gox/config`  
**服务发现：** `git.uhomes.net/uhs-go/gox/discovery`  
**链路追踪：** `git.uhomes.net/uhs-go/gox/trace`  
**指标：** `git.uhomes.net/uhs-go/gox/metrics`  
**中间件：** `git.uhomes.net/uhs-go/gox/middleware`  
**传输层：** `git.uhomes.net/uhs-go/gox/transport`  
**工具：** `git.uhomes.net/uhs-go/gox/utils`

详细使用方法请参考 [gox 基础库文档](./gox_blueprint.md)
