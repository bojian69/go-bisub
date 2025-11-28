package models

import "time"

// SubRefs 字段参考表
type SubRefs struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	RefField  string    `json:"ref_field" gorm:"column:ref_field;size:50;not null;default:''"`
	RefValue  string    `json:"ref_value" gorm:"column:ref_value;size:1;not null;default:''"`
	RefName   string    `json:"ref_name" gorm:"column:ref_name;size:100;not null;default:''"`
	RefNameEn string    `json:"ref_name_en" gorm:"column:ref_name_en;size:100;not null;default:''"`
	Sort      uint8     `json:"sort" gorm:"column:sort;not null;default:0"`
}

func (SubRefs) TableName() string {
	return "sub_refs"
}

// RefOption 参考选项
type RefOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Sort  uint8  `json:"sort"`
}
