package repository

import (
	"context"
	"errors"
	"strconv"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

func (r *StrategyRepository) QueryStrategyRule(ctx context.Context, strategyID int64, ruleModel string) (strategy.StrategyRuleEntity, bool, error) {
	var strategyRulePO po.StrategyRule
	err := r.defaultDB(ctx).
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

func (r *StrategyRepository) QueryStrategyAwardRuleModels(ctx context.Context, strategyID int64, awardID int) (string, error) {
	var awardPO po.StrategyAward
	err := r.defaultDB(ctx).
		Select("rule_models").
		Where("strategy_id = ? and award_id = ?", strategyID, awardID).
		First(&awardPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return awardPO.RuleModels, nil
}

func (r *StrategyRepository) QueryStrategyRuleValue(ctx context.Context, strategyID int64, ruleModel string) (string, error) {
	var strategyRulePO po.StrategyRule
	err := r.defaultDB(ctx).
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

func (r *StrategyRepository) QueryAwardRuleLockCount(ctx context.Context, treeIDs []string) (map[string]int, error) {
	if len(treeIDs) == 0 {
		return map[string]int{}, nil
	}

	var nodes []po.RuleTreeNode
	err := r.defaultDB(ctx).
		Select("tree_id", "rule_value").
		Where("rule_key = ? and tree_id in ?", "rule_lock", treeIDs).
		Find(&nodes).
		Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]int, len(nodes))
	for _, node := range nodes {
		ruleValue, err := strconv.Atoi(node.RuleValue)
		if err != nil {
			return nil, err
		}
		result[node.TreeID] = ruleValue
	}
	return result, nil
}
