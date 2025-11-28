package repository

import (
	"context"
	"time"

	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 如果状态为强制兼容，先将同key的低版本设为失效
	if subscription.Status == models.StatusActiveForceCompatible {
		if err := tx.Model(&models.Subscription{}).
			Where("type = ? AND sub_key = ? AND version < ? AND status IN (?, ?)",
				subscription.Type, subscription.SubKey, subscription.Version,
				models.StatusActive, models.StatusActiveForceCompatible).
			Update("status", models.StatusExpired).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Create(subscription).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *SubscriptionRepository) GetByKeyAndVersion(ctx context.Context, subType, key string, version uint8) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.WithContext(ctx).Where("type = ? AND sub_key = ? AND version = ?", subType, key, version).First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) GetActiveByKey(ctx context.Context, subType, key string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.WithContext(ctx).
		Where("type = ? AND sub_key = ? AND status = ?", subType, key, models.StatusActive).
		Order("version DESC").
		First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) List(ctx context.Context, limit, offset int) ([]*models.Subscription, int64, error) {
	var subscriptions []*models.Subscription
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Subscription{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&subscriptions).Error

	return subscriptions, total, err
}

func (r *SubscriptionRepository) Update(ctx context.Context, subscription *models.Subscription) error {
	return r.db.WithContext(ctx).Save(subscription).Error
}

func (r *SubscriptionRepository) UpdateFields(ctx context.Context, subType, key string, version uint8, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Subscription{}).
		Where("type = ? AND sub_key = ? AND version = ?", subType, key, version).
		Updates(updates).Error
}

func (r *SubscriptionRepository) UpdateStatus(ctx context.Context, subType, key string, version uint8, status string) error {
	return r.db.WithContext(ctx).Model(&models.Subscription{}).
		Where("type = ? AND sub_key = ? AND version = ?", subType, key, version).
		Update("status", status).Error
}

func (r *SubscriptionRepository) Delete(ctx context.Context, subType, key string, version uint8) error {
	return r.db.WithContext(ctx).
		Where("type = ? AND sub_key = ? AND version = ?", subType, key, version).
		Delete(&models.Subscription{}).Error
}

type StatsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) Create(ctx context.Context, stats *models.SubscriptionStats) error {
	return r.db.WithContext(ctx).Create(stats).Error
}

func (r *StatsRepository) GetStats(ctx context.Context, startTime, endTime time.Time, limit, offset int) ([]*models.StatsResponse, error) {
	var results []*models.StatsResponse

	query := `
		SELECT 
			s.sub_key,
			s.version,
			COUNT(*) as call_count,
			AVG(s.execution_duration) as avg_execution_time,
			MIN(s.execution_duration) as min_execution_time,
			MAX(s.execution_duration) as max_execution_time,
			(SELECT JSON_EXTRACT(request_response, '$.instance_sql') FROM sub_logs_bidata_response 
			 WHERE sub_key = s.sub_key AND version = s.version 
			 AND created_at BETWEEN ? AND ? 
			 ORDER BY execution_duration ASC LIMIT 1) as fastest_sql,
			(SELECT JSON_EXTRACT(request_response, '$.instance_sql') FROM sub_logs_bidata_response 
			 WHERE sub_key = s.sub_key AND version = s.version 
			 AND created_at BETWEEN ? AND ? 
			 ORDER BY execution_duration DESC LIMIT 1) as slowest_sql,
			sub.created_by
		FROM sub_logs_bidata_response s
		LEFT JOIN sub_subscription_theme sub ON s.sub_key = sub.sub_key AND s.version = sub.version
		WHERE s.created_at BETWEEN ? AND ?
		GROUP BY s.sub_key, s.version, sub.created_by
		ORDER BY avg_execution_time DESC
		LIMIT ? OFFSET ?
	`

	err := r.db.WithContext(ctx).Raw(query,
		startTime, endTime, startTime, endTime, startTime, endTime, limit, offset).
		Scan(&results).Error

	return results, err
}
