# 更新日志

## 表结构优化 (2025-02-19)

### 主要变更

#### 1. 数据库名称变更
- 旧: `go_bisub` → 新: `go_sub`

#### 2. 订阅表 (sub_subscription_theme)
**字段变更：**
- `subscription_key` → `sub_key` (VARCHAR(120))
- `description` → `abstract` (TINYTEXT)
- `creator_id` → `created_by` (BIGINT UNSIGNED)
- `sql_content`, `sql_variables` → 合并到 `extra_config` (JSON)
- 新增 `type` (CHAR(1)) - 订阅类型

**状态码变更：**
- 旧: 数字状态 (0,1,2,3)
- 新: 字符状态 (A,B,C,D)
  - A: 待生效
  - B: 生效中
  - C: 生效中-强制兼容低版本
  - D: 已失效

**extra_config 结构：**
```json
{
  "sql_content": "订阅数据SQL",
  "sql_replace": {"variable_replace": "变量说明"},
  "example": "示例说明"
}
```

#### 3. 统计表 (sub_logs_bidata_response)
**字段变更：**
- `subscription_key` → `sub_key` (VARCHAR(120))
- `subscription_version` → `version` (TINYINT UNSIGNED)
- `api_url` → `request_url` (VARCHAR(1000))
- `execution_time` → `execution_duration` (MEDIUMINT UNSIGNED)
- `request_params`, `executed_sql`, `client_ip` → 合并到 `request_response` (JSON)
- 新增 `instance_source` (VARCHAR(120)) - 数据实例标识

**request_response 结构：**
```json
{
  "params": "请求参数",
  "instance_sql": "执行实例SQL",
  "instance_source": "实例来源",
  "request_ip": "请求来源IP",
  "version": "版本号"
}
```

#### 4. 新增参考表 (sub_refs)
用于存储字段参考值和说明：
- SUBSCRIPTION_TYPE: 订阅类型
- SUBSCRIPTION_STATUS: 订阅状态

### API 变更

#### 创建订阅请求
```json
// 旧格式
{
  "subscription_key": "house_report",
  "version": 1,
  "title": "房源报表",
  "description": "获取房源基本信息",
  "sql_content": "SELECT * FROM houses WHERE id = house_id_replace",
  "sql_variables": {"house_id_replace": "房源ID"},
  "status": 1
}

// 新格式
{
  "type": "A",
  "sub_key": "house_report",
  "version": 1,
  "title": "房源报表",
  "abstract": "获取房源基本信息",
  "status": "B",
  "extra_config": {
    "sql_content": "SELECT * FROM houses WHERE id = house_id_replace",
    "sql_replace": {"house_id_replace": "房源ID"},
    "example": "{\"house_id_replace\": 1}"
  }
}
```

#### 执行订阅请求
新增 `type` 查询参数（默认为 "A"）：
```
POST /v1/subscriptions/{key}:execute?type=A
```

### 代码变更

#### 模型层 (internal/models)
- 更新 `Subscription` 结构体字段
- 更新 `SubscriptionStats` 结构体字段
- 新增 `ExtraConfig` 和 `RequestResponse` 结构体
- 状态常量从数字改为字符串

#### 仓库层 (internal/repository)
- 更新所有查询条件，增加 `type` 字段
- 更新表名和字段名
- 优化统计查询SQL

#### 服务层 (internal/service)
- 更新方法签名，增加 `subType` 参数
- 解析和处理 `extra_config` JSON
- 更新统计记录逻辑

#### 处理器层 (internal/handler)
- 更新API参数解析
- 增加 `type` 查询参数支持
- 更新数据类型 (int64 → uint64, int8 → uint8)

### 配置变更
- `config.yaml`: 数据库名称更新为 `go_sub`
- `docker-compose.yml`: 数据库名称更新
- `init.sql`: 完整的新表结构

### 兼容性说明
此次更新为破坏性变更，需要：
1. 重新创建数据库表
2. 更新所有API调用代码
3. 迁移现有数据（如有）

### 迁移步骤
1. 备份现有数据
2. 执行新的 `init.sql` 创建表结构
3. 如需迁移数据，需要编写数据转换脚本
4. 更新应用配置文件
5. 重新部署应用
