package service

import (
	"context"

	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"git.uhomes.net/uhs-go/go-bisub/internal/repository"
)

type RefsService struct {
	repo *repository.RefsRepository
}

func NewRefsService(repo *repository.RefsRepository) *RefsService {
	return &RefsService{repo: repo}
}

// GetSubscriptionTypes 获取订阅类型列表
func (s *RefsService) GetSubscriptionTypes(ctx context.Context) ([]*models.RefOption, error) {
	return s.repo.GetRefOptions(ctx, "SUBSCRIPTION_TYPE")
}

// GetSubscriptionStatuses 获取订阅状态列表
func (s *RefsService) GetSubscriptionStatuses(ctx context.Context) ([]*models.RefOption, error) {
	return s.repo.GetRefOptions(ctx, "SUBSCRIPTION_STATUS")
}
