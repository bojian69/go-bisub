package models

import (
	"encoding/json"
	"time"
)

// OperationType 操作类型
const (
	OpTypeCreate  = "CREATE"  // 创建
	OpTypeUpdate  = "UPDATE"  // 更新
	OpTypeDelete  = "DELETE"  // 删除
	OpTypeExecute = "EXECUTE" // 执行
	OpTypeQuery   = "QUERY"   // 查询
	OpTypeLogin   = "LOGIN"   // 登录
	OpTypeLogout  = "LOGOUT"  // 登出
)

// OperationStatus 操作状态
const (
	OpStatusSuccess = "SUCCESS" // 成功
	OpStatusFailed  = "FAILED"  // 失败
)

// OperationLog 操作日志模型
type OperationLog struct {
	ID           uint64          `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	UserID       uint64          `json:"user_id" gorm:"column:user_id;not null;default:0"`    // 操作用户ID
	Username     string          `json:"username" gorm:"column:username;size:120;not null"`   // 操作用户名
	Operation    string          `json:"operation" gorm:"column:operation;size:50;not null"`  // 操作类型
	Resource     string          `json:"resource" gorm:"column:resource;size:200;not null"`   // 操作资源
	ResourceID   string          `json:"resource_id" gorm:"column:resource_id;size:120"`      // 资源ID
	Status       string          `json:"status" gorm:"column:status;size:20;not null"`        // 操作状态
	ClientIP     string          `json:"client_ip" gorm:"column:client_ip;size:45;not null"`  // 客户端IP
	UserAgent    string          `json:"user_agent" gorm:"column:user_agent;size:500"`        // 用户代理
	RequestURL   string          `json:"request_url" gorm:"column:request_url;size:1000"`     // 请求URL
	Method       string          `json:"method" gorm:"column:method;size:10;not null"`        // HTTP方法
	Duration     uint32          `json:"duration" gorm:"column:duration;not null;default:0"`  // 执行耗时(毫秒)
	ErrorMsg     string          `json:"error_msg" gorm:"column:error_msg;type:text"`         // 错误信息
	RequestData  json.RawMessage `json:"request_data" gorm:"column:request_data;type:json"`   // 请求数据
	ResponseData json.RawMessage `json:"response_data" gorm:"column:response_data;type:json"` // 响应数据
}

func (OperationLog) TableName() string {
	return "sub_logs_operation"
}

// OperationLogRequest 操作日志查询请求
type OperationLogRequest struct {
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	UserID    uint64 `form:"user_id"`
	Username  string `form:"username"`
	Operation string `form:"operation"`
	Resource  string `form:"resource"`
	Status    string `form:"status"`
	ClientIP  string `form:"client_ip"`
	Limit     int    `form:"limit"`
	Offset    int    `form:"offset"`
}
