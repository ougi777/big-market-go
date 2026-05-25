package repository

import (
	"context"
	"errors"
	"strconv"

	"bm-go/internal/domain/strategy"
	treepkg "bm-go/internal/domain/strategy/rule/tree"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/infrastructure/persistent/sharding"

	"gorm.io/gorm"
)

type StrategyRepository struct {
	db         dbRouter
	sharder    sharding.Router
	stockQueue AwardStockQueue
}

type AwardStockQueue interface {
	AwardStockConsumeSendQueue(ctx context.Context, strategyID int64, awardID int) error
}

var _ strategy.Repository = (*StrategyRepository)(nil)
var _ strategy.ArmoryRepository = (*StrategyRepository)(nil)
var _ strategy.QueryRepository = (*StrategyRepository)(nil)
var _ strategy.RaffleRepository = (*StrategyRepository)(nil)
var _ strategy.StockRepository = (*StrategyRepository)(nil)
var _ treepkg.Repository = (*StrategyRepository)(nil)

func NewStrategyRepository(db *gorm.DB, stockQueues ...AwardStockQueue) *StrategyRepository {
	return NewStrategyRepositoryWithDBRouter(singleDBRouter{db: db}, sharding.NewRouter(1), stockQueues...)
}

func NewStrategyRepositoryWithDBRouter(db dbRouter, sharder sharding.Router, stockQueues ...AwardStockQueue) *StrategyRepository {
	var stockQueue AwardStockQueue
	if len(stockQueues) > 0 {
		stockQueue = stockQueues[0]
	}
	return &StrategyRepository{db: db, sharder: sharder, stockQueue: stockQueue}
}

func (r *StrategyRepository) defaultDB(ctx context.Context) *gorm.DB {
	return r.db.Default().WithContext(ctx)
}

func (r *StrategyRepository) shardDB(ctx context.Context, userID string) *gorm.DB {
	return r.db.Shard(r.sharder.DBKey(userID)).WithContext(ctx)
}

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
