# Protobuf JSON 命名规范

**文档版本：** 1.0
**最后更新：** 2025-02-19
**状态：** 正式标准 (Standard)

## 0. 规范结论 (Normative Conclusion)

**Uhomes 微服务架构强制标准：**

1.  **Protobuf 字段名**：必须使用 `snake_case` (e.g., `user_id`, `created_at`)。
2.  **JSON 输出**：自动映射为 `snake_case`，与 Protobuf 字段名保持一致。
3.  **禁止行为**：
    *   ❌ 禁止在 Protobuf 中使用 `camelCase` (e.g., `userId`)。
    *   ❌ 禁止依赖 Go 结构体的 `PascalCase` 字段名来推断 JSON 字段名。
    *   ❌ 除非特殊兼容需求，否则不建议使用 `json_name` 选项。

---

## 1. 官方规则总结

**核心规则：Protobuf 字段名 → JSON 字段名是 1:1 映射，不会自动转换命名风格**

来源：[Protocol Buffers - JSON Mapping](https://protobuf.dev/programming-guides/proto3/#json)

---

## 测试案例

### 案例 1：snake_case 字段（推荐）

**Protobuf 定义：**
```protobuf
syntax = "proto3";
package pb.user.v1;

message User {
  string user_id = 1;
  string user_name = 2;
  string email_address = 3;
  int64 created_at = 4;
  bool is_active = 5;
}
```

**生成的 Go 代码：**
```go
type User struct {
    UserId       string `protobuf:"bytes,1,opt,name=user_id,json=user_id,proto3" json:"user_id,omitempty"`
    UserName     string `protobuf:"bytes,2,opt,name=user_name,json=user_name,proto3" json:"user_name,omitempty"`
    EmailAddress string `protobuf:"bytes,3,opt,name=email_address,json=email_address,proto3" json:"email_address,omitempty"`
    CreatedAt    int64  `protobuf:"varint,4,opt,name=created_at,json=created_at,proto3" json:"created_at,omitempty"`
    IsActive     bool   `protobuf:"varint,5,opt,name=is_active,json=is_active,proto3" json:"is_active,omitempty"`
}
```

**JSON 输出：**
```json
{
  "user_id": "usr_123",
  "user_name": "张三",
  "email_address": "zhangsan@example.com",
  "created_at": 1708336200,
  "is_active": true
}
```

**结论：✅ snake_case → snake_case（保持不变）**

---

### 案例 2：camelCase 字段（不推荐）

**Protobuf 定义：**
```protobuf
syntax = "proto3";
package pb.user.v1;

message User {
  string userId = 1;
  string userName = 2;
  string emailAddress = 3;
  int64 createdAt = 4;
  bool isActive = 5;
}
```

**生成的 Go 代码：**
```go
type User struct {
    UserId       string `protobuf:"bytes,1,opt,name=userId,json=userId,proto3" json:"userId,omitempty"`
    UserName     string `protobuf:"bytes,2,opt,name=userName,json=userName,proto3" json:"userName,omitempty"`
    EmailAddress string `protobuf:"bytes,3,opt,name=emailAddress,json=emailAddress,proto3" json:"emailAddress,omitempty"`
    CreatedAt    int64  `protobuf:"varint,4,opt,name=createdAt,json=createdAt,proto3" json:"createdAt,omitempty"`
    IsActive     bool   `protobuf:"varint,5,opt,name=isActive,json=isActive,proto3" json:"isActive,omitempty"`
}
```

**JSON 输出：**
```json
{
  "userId": "usr_123",
  "userName": "张三",
  "emailAddress": "zhangsan@example.com",
  "createdAt": 1708336200,
  "isActive": true
}
```

**结论：❌ camelCase → camelCase（保持不变，但不符合云原生标准）**

---

### 案例 3：使用 json_name 自定义（兼容场景）

**Protobuf 定义：**
```protobuf
syntax = "proto3";
package pb.user.v1;

message User {
  // Go 代码中是 UserId，但 JSON 输出是 user_id
  string userId = 1 [json_name = "user_id"];
  string userName = 2 [json_name = "user_name"];
  string emailAddress = 3 [json_name = "email_address"];
  int64 createdAt = 4 [json_name = "created_at"];
  bool isActive = 5 [json_name = "is_active"];
}
```

**生成的 Go 代码：**
```go
type User struct {
    UserId       string `protobuf:"bytes,1,opt,name=userId,json=user_id,proto3" json:"user_id,omitempty"`
    UserName     string `protobuf:"bytes,2,opt,name=userName,json=user_name,proto3" json:"user_name,omitempty"`
    EmailAddress string `protobuf:"bytes,3,opt,name=emailAddress,json=email_address,proto3" json:"email_address,omitempty"`
    CreatedAt    int64  `protobuf:"varint,4,opt,name=createdAt,json=created_at,proto3" json:"created_at,omitempty"`
    IsActive     bool   `protobuf:"varint,5,opt,name=isActive,json=is_active,proto3" json:"is_active,omitempty"`
}
```

**JSON 输出：**
```json
{
  "user_id": "usr_123",
  "user_name": "张三",
  "email_address": "zhangsan@example.com",
  "created_at": 1708336200,
  "is_active": true
}
```

**结论：✅ 可以通过 json_name 自定义，但增加维护成本**

---

## Protobuf 官方风格指南

来源：[Protobuf Style Guide](https://protobuf.dev/programming-guides/style/)

### 字段命名规则

**推荐：**
```protobuf
✅ 使用 lower_snake_case（小写下划线）

message User {
  string user_id = 1;
  string first_name = 2;
  string last_name = 3;
  string email_address = 4;
  int64 created_at = 5;
  int64 updated_at = 6;
  bool is_active = 7;
  bool is_verified = 8;
}
```

**不推荐：**
```protobuf
❌ camelCase
message User {
  string userId = 1;
  string firstName = 2;
}

❌ PascalCase
message User {
  string UserId = 1;
  string FirstName = 2;
}

❌ UPPER_SNAKE_CASE（仅用于枚举值）
message User {
  string USER_ID = 1;  // 错误
}
```

### 枚举命名规则

```protobuf
enum Status {
  STATUS_UNSPECIFIED = 0;  // 必须有 0 值
  STATUS_PENDING = 1;
  STATUS_APPROVED = 2;
  STATUS_REJECTED = 3;
}
```

---

## 不同语言的生成代码对比

### Go

**Protobuf：**
```protobuf
message User {
  string user_id = 1;
  string user_name = 2;
}
```

**生成的 Go 代码：**
```go
type User struct {
    UserId   string `json:"user_id,omitempty"`  // Go 使用 PascalCase
    UserName string `json:"user_name,omitempty"` // JSON 使用 snake_case
}
```

### Python

**生成的 Python 代码：**
```python
class User:
    user_id: str      # Python 使用 snake_case
    user_name: str    # 与 Protobuf 字段名一致
```

### Java

**生成的 Java 代码：**
```java
public class User {
    private String userId;     // Java 使用 camelCase
    private String userName;   // 但 JSON 序列化是 user_id
    
    @JsonProperty("user_id")
    public String getUserId() { return userId; }
}
```

### TypeScript

**生成的 TypeScript 代码：**
```typescript
interface User {
  user_id: string;    // TypeScript 保持 snake_case
  user_name: string;  // 与 Protobuf 字段名一致
}
```

---

## 实际测试代码

### 测试 1：验证 snake_case

```go
package main

import (
    "encoding/json"
    "fmt"
    pb "your-project/proto/user/v1"
)

func main() {
    user := &pb.User{
        UserId:       "usr_123",
        UserName:     "张三",
        EmailAddress: "zhangsan@example.com",
        CreatedAt:    1708336200,
        IsActive:     true,
    }
    
    jsonData, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println(string(jsonData))
}
```

**输出：**
```json
{
  "user_id": "usr_123",
  "user_name": "张三",
  "email_address": "zhangsan@example.com",
  "created_at": 1708336200,
  "is_active": true
}
```

### 测试 2：验证 gRPC-Gateway

```go
// gRPC-Gateway 自动将 Protobuf 转换为 JSON
// 使用 Protobuf 字段名的 json tag

// 如果 Protobuf 定义是 snake_case
message User {
  string user_id = 1;
}

// gRPC-Gateway 输出的 JSON 就是
{
  "user_id": "123"
}
```

---

## 最佳实践建议

### ✅ 推荐做法

1. **Protobuf 字段名使用 snake_case**
   ```protobuf
   message User {
     string user_id = 1;
     string created_at = 2;
   }
   ```

2. **不需要使用 json_name**
   - 字段名已经是 snake_case
   - JSON 输出自动是 snake_case
   - 减少维护成本

3. **符合官方风格指南**
   - Protobuf Style Guide 推荐
   - 云原生生态标准
   - 跨语言一致性

### ❌ 避免做法

1. **不要在 Protobuf 中使用 camelCase**
   ```protobuf
   ❌ 错误示例
   message User {
     string userId = 1;      // 不推荐
     string userName = 2;    // 不推荐
   }
   ```

2. **不要过度使用 json_name**
   ```protobuf
   ❌ 不必要的复杂性
   message User {
     string user_id = 1 [json_name = "user_id"];  // 多余
   }
   
   ✅ 简洁写法
   message User {
     string user_id = 1;  // 已经是 snake_case，无需 json_name
   }
   ```

---

## 常见误解澄清

### 误解 1：Protobuf 会自动转换命名风格

**错误认知：**
> "我在 Protobuf 中写 userId，JSON 会自动变成 user_id"

**正确理解：**
> Protobuf 不会自动转换命名风格。字段名是什么，JSON 就是什么。

### 误解 2：Go 的 PascalCase 会影响 JSON

**错误认知：**
> "Go 代码中是 UserId，所以 JSON 也应该是 UserId"

**正确理解：**
> Go 的结构体字段名（UserId）是 Go 语言的要求（首字母大写才能导出）。
> JSON 字段名由 Protobuf 定义和 json tag 决定，不受 Go 字段名影响。

### 误解 3：需要为每个字段添加 json_name

**错误认知：**
> "为了输出 snake_case JSON，需要给每个字段加 json_name"

**正确理解：**
> 如果 Protobuf 字段名已经是 snake_case，不需要 json_name。
> 只有当 Protobuf 字段名和期望的 JSON 字段名不同时，才需要 json_name。

---

## 参考资源

1. **Protobuf 官方文档**
   - [JSON Mapping](https://protobuf.dev/programming-guides/proto3/#json)
   - [Style Guide](https://protobuf.dev/programming-guides/style/)

2. **gRPC-Gateway 文档**
   - [JSON Naming](https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/)

3. **云原生项目示例**
   - [Kubernetes API Conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
   - [Envoy API](https://www.envoyproxy.io/docs/envoy/latest/api-v3/api)

---

## 结论

**Uhomes 微服务标准：**

1. ✅ Protobuf 字段名使用 `snake_case`
2. ✅ JSON 输出自动是 `snake_case`
3. ✅ 不需要使用 `json_name`（除非有特殊需求）
4. ✅ 符合云原生生态标准
5. ✅ 符合 Protobuf 官方风格指南

**示例：**
```protobuf
syntax = "proto3";
package pb.user.v1;

option go_package = "git.uhomes.net/uhs-go/proto-contracts/user/v1;userv1";

message User {
  string user_id = 1;
  string user_name = 2;
  string email_address = 3;
  int64 created_at = 4;
  int64 updated_at = 5;
  bool is_active = 6;
}
```

**JSON 输出：**
```json
{
  "user_id": "usr_123",
  "user_name": "张三",
  "email_address": "zhangsan@example.com",
  "created_at": 1708336200,
  "updated_at": 1708336200,
  "is_active": true
}
```
