package service

import (
	"context"
	"strconv"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"
)

type ArmoryService struct {
	repo  activity.SkuProductRepository
	store activity.SkuStockStore
}

func NewArmoryService(repo activity.SkuProductRepository, store activity.SkuStockStore) *ArmoryService {
	return &ArmoryService{
		repo:  repo,
		store: store,
	}
}

func (s *ArmoryService) AssembleActivitySkuByActivityID(ctx context.Context, activityID int64) error {
	products, err := s.repo.QuerySkuProductListByActivityID(ctx, activityID)
	if err != nil {
		return err
	}

	for _, product := range products {
		cacheKey := types.RedisKeyActivitySkuStockCount + strconv.FormatInt(product.SKU, 10)
		if err := s.store.CacheActivitySkuStockCount(ctx, cacheKey, product.StockCountSurplus); err != nil {
			return err
		}
	}
	return nil
}
