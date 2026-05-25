package repository

import (
	"context"
	"errors"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

func (r *StrategyRepository) QueryStrategyEntityByStrategyID(ctx context.Context, strategyID int64) (strategy.StrategyEntity, error) {
	var strategyPO po.Strategy
	err := r.defaultDB(ctx).
		Select("strategy_id", "strategy_desc", "rule_models").
		Where("strategy_id = ?", strategyID).
		First(&strategyPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return strategy.StrategyEntity{}, nil
	}
	if err != nil {
		return strategy.StrategyEntity{}, err
	}

	return strategy.StrategyEntity{
		StrategyID:   strategyPO.StrategyID,
		StrategyDesc: strategyPO.StrategyDesc,
		RuleModel:    strategyPO.RuleModels,
	}, nil
}

func (r *StrategyRepository) QueryStrategyIDByActivityID(ctx context.Context, activityID int64) (int64, error) {
	var activityPO po.RaffleActivity
	err := r.defaultDB(ctx).
		Select("strategy_id").
		Where("activity_id = ?", activityID).
		First(&activityPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return activityPO.StrategyID, nil
}
