package repository

import (
	"context"
	"errors"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/infrastructure/persistent/po"

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
	var strategyPO po.Strategy
	err := r.db.WithContext(ctx).
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
		StrategyID: strategyPO.StrategyID,
		RuleModel:  strategyPO.RuleModels,
	}, nil
}

func (r *StrategyRepository) QueryStrategyRuleValue(ctx context.Context, strategyID int64, ruleModel string) (string, error) {
	var strategyRulePO po.StrategyRule
	err := r.db.WithContext(ctx).
		Select("rule_value").
		Where("strategy_id = ? and rule_model = ?", strategyID, ruleModel).
		First(&strategyRulePO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strategyRulePO.RuleValue, nil
}

func (r *StrategyRepository) QueryActivityAccountTotalUseCount(ctx context.Context, userID string, strategyID int64) (int, error) {
	var activityPO po.RaffleActivity
	err := r.db.WithContext(ctx).
		Select("activity_id").
		Where("strategy_id = ?", strategyID).
		First(&activityPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	var accountPO po.RaffleActivityAccount
	err = r.db.WithContext(ctx).
		Select("total_count", "total_count_surplus").
		Where("user_id = ? and activity_id = ?", userID, activityPO.ActivityID).
		First(&accountPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return accountPO.TotalCount - accountPO.TotalCountSurplus, nil
}
