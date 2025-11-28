package models

import (
	"encoding/json"
	"time"
)

// SubscriptionStatus 订阅状态
const (
	StatusPending              = "A" // 待生效
	StatusActive               = "B" // 生效中
	StatusActiveForceCompatible = "C" // 生效中（强制兼容低版本）
	StatusExpired              = "D" // 已失效
)

// SubscriptionType 订阅类型
const (
	TypeAnalysisData = "A" // 分析数据
)

// ExtraConfig 扩展配置
type ExtraConfig struct {
	SQLContent  string            `json:"sql_content"`  // 订阅数据SQL
	SQLReplace  map[string]string `json:"sql_replace"`  // SQL替换变量说明
	Example     string            `json:"example"`      // 示例说明
}

// Subscription 订阅模型
type Subscription struct {
	ID          uint64          `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Type        string          `json:"type" gorm:"column:type;size:1;not null;default:''"`
	SubKey      string          `json:"sub_key" gorm:"column:sub_key;size:120;not null;default:''"`
	Version     uint8           `json:"version" gorm:"column:version;not null;default:1"`
	Title       string          `json:"title" gorm:"column:title;size:240;not null"`
	Abstract    string          `json:"abstract" gorm:"column:abstract;type:tinytext;not null"`
	Status      string          `json:"status" gorm:"column:status;size:1;not null;default:''"`
	CreatedBy   uint64          `json:"created_by" gorm:"column:created_by;not null;default:0"`
	ExtraConfig json.RawMessage `json:"extra_config" gorm:"column:extra_config;type:json;not null"`
}

func (Subscription) TableName() string {
	return "sub_subscription_theme"
}

// RequestResponse 请求响应详情
type RequestResponse struct {
	Params         interface{} `json:"params"`          // 请求参数
	InstanceSQL    string      `json:"instance_sql"`    // 执行实例SQL
	InstanceSource string      `json:"instance_source"` // 实例来源
	RequestIP      string      `json:"request_ip"`      // 请求来源IP
	Version        uint8       `json:"version"`         // 版本号
}

// SubscriptionStats 订阅统计模型
type SubscriptionStats struct {
	ID                uint64          `json:"id" gorm:"primaryKey"`
	CreatedAt         time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	SubKey            string          `json:"sub_key" gorm:"column:sub_key;size:120;not null;default:''"`
	Version           uint8           `json:"version" gorm:"column:version;not null;default:1"`
	ExecutionDuration uint32          `json:"execution_duration" gorm:"column:execution_duration;not null;default:0"` // 毫秒
	RequestURL        string          `json:"request_url" gorm:"column:request_url;size:1000;not null;default:''"`
	RequestResponse   json.RawMessage `json:"request_response" gorm:"column:request_response;type:json;not null"`
	InstanceSource    string          `json:"instance_source" gorm:"column:instance_source;size:120;not null;default:''"`
}

func (SubscriptionStats) TableName() string {
	return "sub_logs_bidata_response"
}

// CreateSubscriptionRequest 创建订阅请求
type CreateSubscriptionRequest struct {
	Type        string          `json:"type" binding:"required,len=1"`
	SubKey      string          `json:"sub_key" binding:"required"`
	Version     uint8           `json:"version" binding:"required,min=1"`
	Title       string          `json:"title" binding:"required"`
	Abstract    string          `json:"abstract" binding:"required"`
	Status      string          `json:"status" binding:"required,len=1"`
	ExtraConfig json.RawMessage `json:"extra_config" binding:"required"`
}

// ExecuteSubscriptionRequest 执行订阅请求
type ExecuteSubscriptionRequest struct {
	Variables  map[string]interface{} `json:"variables"`
	Timeout    int                    `json:"timeout"`    // 毫秒，默认120000
	DataSource string                 `json:"data_source"` // 数据源名称，默认default
}

// StatsQueryRequest 统计查询请求
type StatsQueryRequest struct {
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Limit     int    `form:"limit"`
	Offset    int    `form:"offset"`
}

// StatsResponse 统计响应
type StatsResponse struct {
	SubKey              string  `json:"sub_key"`
	Version             uint8   `json:"version"`
	CallCount           int64   `json:"call_count"`
	AvgExecutionTime    float64 `json:"avg_execution_time"`
	MinExecutionTime    uint32  `json:"min_execution_time"`
	MaxExecutionTime    uint32  `json:"max_execution_time"`
	FastestSQL          string  `json:"fastest_sql"`
	SlowestSQL          string  `json:"slowest_sql"`
	CreatedBy           uint64  `json:"created_by"`
}