package service

import (
	"context"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"
)

type SkuProductService struct {
	repo activity.SkuProductRepository
}

func NewSkuProductService(repo activity.SkuProductRepository) *SkuProductService {
	return &SkuProductService{repo: repo}
}

func (s *SkuProductService) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	if activityID <= 0 {
		return nil, types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}
	return s.repo.QuerySkuProductListByActivityID(ctx, activityID)
}
