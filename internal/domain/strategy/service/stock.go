package service

import (
	"context"

	"bm-go/internal/domain/strategy"
)

type StockService struct {
	repo  strategy.StockRepository
	queue strategy.StockQueue
}

func NewStockService(repo strategy.StockRepository, queue strategy.StockQueue) *StockService {
	return &StockService{repo: repo, queue: queue}
}

func (s *StockService) UpdateAwardStock(ctx context.Context) (bool, error) {
	stockKey, ok, err := s.queue.TakeQueueValue(ctx)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	if err := s.repo.UpdateStrategyAwardStock(ctx, stockKey.StrategyID, stockKey.AwardID); err != nil {
		return false, err
	}
	return true, nil
}
