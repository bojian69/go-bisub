package repository

import (
	"context"
	"time"

	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"gorm.io/gorm"
)

type OperationLogRepository struct {
	db *gorm.DB
}

func NewOperationLogRepository(db *gorm.DB) *OperationLogRepository {
	return &OperationLogRepository{db: db}
}

func (r *OperationLogRepository) Create(ctx context.Context, log *models.OperationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *OperationLogRepository) List(ctx context.Context, req *models.OperationLogRequest) ([]*models.OperationLog, int64, error) {
	var logs []*models.OperationLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.OperationLog{})

	// 时间范围过滤
	if req.StartTime != "" {
		if startTime, err := time.Parse("2006-01-02", req.StartTime); err == nil {
			query = query.Where("created_at >= ?", startTime)
		}
	}
	if req.EndTime != "" {
		if endTime, err := time.Parse("2006-01-02", req.EndTime); err == nil {
			query = query.Where("created_at <= ?", endTime.Add(24*time.Hour))
		}
	}

	// 其他过滤条件
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Operation != "" {
		query = query.Where("operation = ?", req.Operation)
	}
	if req.Resource != "" {
		query = query.Where("resource LIKE ?", "%"+req.Resource+"%")
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.ClientIP != "" {
		query = query.Where("client_ip = ?", req.ClientIP)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}
