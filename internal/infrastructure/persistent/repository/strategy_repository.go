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
var _ strategy.ArmoryRepository = (*StrategyRepository)(nil)

func NewStrategyRepository(db *gorm.DB) *StrategyRepository {
	return &StrategyRepository{db: db}
}

func (r *StrategyRepository) QueryStrategyAwardList(ctx context.Context, strategyID int64) ([]strategy.StrategyAwardEntity, error) {
	var awardPOList []po.StrategyAward
	err := r.db.WithContext(ctx).
		Select("strategy_id", "award_id", "award_title", "award_subtitle", "award_count", "award_count_surplus", "award_rate", "rule_models", "sort").
		Where("strategy_id = ?", strategyID).
		Order("sort asc").
		Find(&awardPOList).
		Error
	if err != nil {
		return nil, err
	}

	awards := make([]strategy.StrategyAwardEntity, 0, len(awardPOList))
	for _, awardPO := range awardPOList {
		awards = append(awards, strategy.StrategyAwardEntity{
			StrategyID:        awardPO.StrategyID,
			AwardID:           awardPO.AwardID,
			AwardTitle:        awardPO.AwardTitle,
			AwardSubtitle:     awardPO.AwardSubtitle,
			AwardCount:        awardPO.AwardCount,
			AwardCountSurplus: awardPO.AwardCountSurplus,
			AwardRate:         awardPO.AwardRate,
			Sort:              awardPO.Sort,
			RuleModels:        awardPO.RuleModels,
		})
	}
	return awards, nil
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
		StrategyID:   strategyPO.StrategyID,
		StrategyDesc: strategyPO.StrategyDesc,
		RuleModel:    strategyPO.RuleModels,
	}, nil
}

func (r *StrategyRepository) QueryStrategyRule(ctx context.Context, strategyID int64, ruleModel string) (strategy.StrategyRuleEntity, bool, error) {
	var strategyRulePO po.StrategyRule
	err := r.db.WithContext(ctx).
		Select("strategy_id", "award_id", "rule_type", "rule_model", "rule_value", "rule_desc").
		Where("strategy_id = ? and rule_model = ?", strategyID, ruleModel).
		First(&strategyRulePO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return strategy.StrategyRuleEntity{}, false, nil
	}
	if err != nil {
		return strategy.StrategyRuleEntity{}, false, err
	}

	return strategy.StrategyRuleEntity{
		StrategyID: strategyRulePO.StrategyID,
		AwardID:    strategyRulePO.AwardID,
		RuleType:   strategyRulePO.RuleType,
		RuleModel:  strategyRulePO.RuleModel,
		RuleValue:  strategyRulePO.RuleValue,
		RuleDesc:   strategyRulePO.RuleDesc,
	}, true, nil
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
