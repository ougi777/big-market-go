package service

import (
	"context"

	"bm-go/internal/domain/activity"
)

type StockService struct {
	repo  activity.SkuStockRepository
	queue activity.SkuStockQueue
}

func NewStockService(repo activity.SkuStockRepository, queue activity.SkuStockQueue) *StockService {
	return &StockService{
		repo:  repo,
		queue: queue,
	}
}

func (s *StockService) UpdateActivitySkuStock(ctx context.Context) (bool, error) {
	stockKey, ok, err := s.queue.TakeActivitySkuStock(ctx)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	if err := s.repo.UpdateActivitySkuStock(ctx, stockKey.SKU); err != nil {
		return false, err
	}
	return true, nil
}

func (s *StockService) ClearActivitySkuStock(ctx context.Context, sku int64) error {
	if err := s.repo.ClearActivitySkuStock(ctx, sku); err != nil {
		return err
	}
	return s.queue.ClearActivitySkuStockQueue(ctx)
}
