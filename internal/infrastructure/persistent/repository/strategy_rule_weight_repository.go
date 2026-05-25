package repository

import (
	"context"
	"strconv"
	"strings"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/infrastructure/persistent/po"
)

func (r *StrategyRepository) QueryAwardRuleWeight(ctx context.Context, strategyID int64) ([]strategy.RuleWeight, error) {
	ruleValue, err := r.QueryStrategyRuleValue(ctx, strategyID, "rule_weight")
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(ruleValue) == "" {
		return nil, nil
	}

	ruleEntity := strategy.StrategyRuleEntity{
		StrategyID: strategyID,
		RuleModel:  "rule_weight",
		RuleValue:  ruleValue,
	}

	ruleWeightValues := ruleEntity.RuleWeightValues()
	result := make([]strategy.RuleWeight, 0, len(ruleWeightValues))
	for ruleWeightKey, awardIDs := range ruleWeightValues {
		awards, err := r.queryRuleWeightAwards(ctx, strategyID, awardIDs)
		if err != nil {
			return nil, err
		}

		weight, err := strconv.Atoi(strings.Split(ruleWeightKey, ":")[0])
		if err != nil {
			return nil, err
		}

		result = append(result, strategy.RuleWeight{
			RuleValue: ruleValue,
			Weight:    weight,
			AwardIDs:  awardIDs,
			AwardList: awards,
		})
	}
	return result, nil
}

func (r *StrategyRepository) queryRuleWeightAwards(ctx context.Context, strategyID int64, awardIDs []int) ([]strategy.RuleWeightAward, error) {
	if len(awardIDs) == 0 {
		return nil, nil
	}

	var awardPOList []po.StrategyAward
	err := r.defaultDB(ctx).
		Select("award_id", "award_title").
		Where("strategy_id = ? and award_id in ?", strategyID, awardIDs).
		Order("sort asc").
		Find(&awardPOList).
		Error
	if err != nil {
		return nil, err
	}

	awards := make([]strategy.RuleWeightAward, 0, len(awardPOList))
	for _, awardPO := range awardPOList {
		awards = append(awards, strategy.RuleWeightAward{
			AwardID:    awardPO.AwardID,
			AwardTitle: awardPO.AwardTitle,
		})
	}
	return awards, nil
}
