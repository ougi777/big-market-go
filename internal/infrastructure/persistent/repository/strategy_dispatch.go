package repository

import (
	"context"

	"bm-go/internal/domain/strategy"
)

type StrategyDispatch struct{}

var _ strategy.Dispatch = (*StrategyDispatch)(nil)

func NewStrategyDispatch() *StrategyDispatch {
	return &StrategyDispatch{}
}

func (d *StrategyDispatch) GetRandomAwardID(ctx context.Context, strategyID int64, ruleWeightValue ...string) (int, error) {
	return 0, errRepositoryNotImplemented
}
