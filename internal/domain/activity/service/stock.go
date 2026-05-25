package service

import (
	"context"
	"strconv"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/award"
	"bm-go/internal/types"
)

type StockService struct {
	repo      activity.SkuStockRepository
	queue     activity.SkuStockQueue
	store     activity.SkuStockStore
	publisher award.MessagePublisher
}

func NewStockService(repo activity.SkuStockRepository, queue activity.SkuStockQueue, store activity.SkuStockStore, publisher award.MessagePublisher) *StockService {
	return &StockService{
		repo:      repo,
		queue:     queue,
		store:     store,
		publisher: publisher,
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

func (s *StockService) SubtractActivitySkuStock(ctx context.Context, sku int64, activityID int64) (bool, error) {
	cacheKey := types.RedisKeyActivitySkuStockCount + strconv.FormatInt(sku, 10)
	surplus, err := s.store.SubtractActivitySkuStock(ctx, cacheKey)
	if err != nil {
		return false, err
	}
	if surplus < 0 {
		return false, nil
	}
	if surplus == 0 && s.publisher != nil {
		message, err := BuildActivitySkuStockZeroMessage(sku)
		if err != nil {
			return false, err
		}
		if err := s.publisher.Publish(ctx, activity.TopicActivitySkuStockZero, message); err != nil {
			return false, err
		}
	}
	if err := s.queue.SendActivitySkuStockConsumeQueue(ctx, activity.ActivitySkuStockKey{SKU: sku, ActivityID: activityID}); err != nil {
		return false, err
	}
	return true, nil
}
