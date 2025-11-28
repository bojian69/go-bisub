package repository

import (
	"context"

	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"gorm.io/gorm"
)

type RefsRepository struct {
	db *gorm.DB
}

func NewRefsRepository(db *gorm.DB) *RefsRepository {
	return &RefsRepository{db: db}
}

// GetByRefField 根据ref_field获取参考数据
func (r *RefsRepository) GetByRefField(ctx context.Context, refField string) ([]*models.SubRefs, error) {
	var refs []*models.SubRefs
	err := r.db.WithContext(ctx).
		Where("ref_field = ?", refField).
		Order("sort ASC").
		Find(&refs).Error
	return refs, err
}

// GetRefOptions 获取参考选项
func (r *RefsRepository) GetRefOptions(ctx context.Context, refField string) ([]*models.RefOption, error) {
	var refs []*models.SubRefs
	err := r.db.WithContext(ctx).
		Where("ref_field = ?", refField).
		Order("sort ASC").
		Find(&refs).Error
	if err != nil {
		return nil, err
	}

	options := make([]*models.RefOption, len(refs))
	for i, ref := range refs {
		options[i] = &models.RefOption{
			Value: ref.RefValue,
			Label: ref.RefName,
			Sort:  ref.Sort,
		}
	}
	return options, nil
}
