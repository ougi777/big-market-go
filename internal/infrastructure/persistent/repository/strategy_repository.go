package repository

import (
	"context"
	"errors"

	"bm-go/internal/domain/strategy"

	"gorm.io/gorm"
)

var errRepositoryNotImplemented = errors.New("repository method is not implemented")

type StrategyRepository struct {
	db *gorm.DB
}

var _ strategy.Repository = (*StrategyRepository)(nil)

func NewStrategyRepository(db *gorm.DB) *StrategyRepository {
	return &StrategyRepository{db: db}
}

func (r *StrategyRepository) QueryStrategyEntityByStrategyID(ctx context.Context, strategyID int64) (strategy.StrategyEntity, error) {
	return strategy.StrategyEntity{}, errRepositoryNotImplemented
}

func (r *StrategyRepository) QueryStrategyRuleValue(ctx context.Context, strategyID int64, ruleModel string) (string, error) {
	return "", errRepositoryNotImplemented
}

func (r *StrategyRepository) QueryActivityAccountTotalUseCount(ctx context.Context, userID string, strategyID int64) (int, error) {
	return 0, errRepositoryNotImplemented
}
