package repository

import (
	"context"
	"errors"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

func (r *StrategyRepository) QueryStrategyAwardList(ctx context.Context, strategyID int64) ([]strategy.StrategyAwardEntity, error) {
	var awardPOList []po.StrategyAward
	err := r.defaultDB(ctx).
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

func (r *StrategyRepository) QueryStrategyAwardEntity(ctx context.Context, strategyID int64, awardID int) (strategy.StrategyAwardEntity, bool, error) {
	var awardPO po.StrategyAward
	err := r.defaultDB(ctx).
		Select("strategy_id", "award_id", "award_title", "award_subtitle", "award_count", "award_count_surplus", "award_rate", "rule_models", "sort").
		Where("strategy_id = ? and award_id = ?", strategyID, awardID).
		First(&awardPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return strategy.StrategyAwardEntity{}, false, nil
	}
	if err != nil {
		return strategy.StrategyAwardEntity{}, false, err
	}

	return strategy.StrategyAwardEntity{
		StrategyID:        awardPO.StrategyID,
		AwardID:           awardPO.AwardID,
		AwardTitle:        awardPO.AwardTitle,
		AwardSubtitle:     awardPO.AwardSubtitle,
		AwardCount:        awardPO.AwardCount,
		AwardCountSurplus: awardPO.AwardCountSurplus,
		AwardRate:         awardPO.AwardRate,
		Sort:              awardPO.Sort,
		RuleModels:        awardPO.RuleModels,
	}, true, nil
}
