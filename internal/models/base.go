package models

import (
	"time"

	"git.uhomes.net/uhs-go/go-bisub/internal/utils"
	"gorm.io/gorm"
)

// BaseModel 基础模型 - 使用Snowflake ID
type BaseModel struct {
	ID        int64     `json:"id" gorm:"primaryKey;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate GORM钩子，创建前生成分布式ID
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == 0 {
		m.ID = utils.GenerateID()
	}
	return nil
}

// StringIDModel 字符串ID基础模型 - 使用UUID v7
type StringIDModel struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate 生成UUID v7
func (m *StringIDModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = utils.GenerateUUIDv7()
	}
	return nil
}

// GetIDTimestamp 获取ID中的时间戳信息
func (m *BaseModel) GetIDTimestamp() time.Time {
	return utils.ParseSnowflakeID(m.ID)
}

// ValidateID 验证ID格式
func (m *BaseModel) ValidateID() error {
	return utils.ValidateSnowflakeID(m.ID)
}