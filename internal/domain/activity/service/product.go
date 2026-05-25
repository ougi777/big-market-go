package service

import (
	"context"

	"bm-go/internal/domain/activity"
)

type SkuProductService struct {
	repo activity.SkuProductRepository
}

func NewSkuProductService(repo activity.SkuProductRepository) *SkuProductService {
	return &SkuProductService{repo: repo}
}

func (s *SkuProductService) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	return s.repo.QuerySkuProductListByActivityID(ctx, activityID)
}
