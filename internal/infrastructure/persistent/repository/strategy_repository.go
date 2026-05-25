package repository

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"bm-go/internal/domain/strategy"
	treepkg "bm-go/internal/domain/strategy/rule/tree"
	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

var errRepositoryNotImplemented = errors.New("repository method is not implemented")

type StrategyRepository struct {
	db         *gorm.DB
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
	var stockQueue AwardStockQueue
	if len(stockQueues) > 0 {
		stockQueue = stockQueues[0]
	}
	return &StrategyRepository{db: db, stockQueue: stockQueue}
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

func (r *StrategyRepository) QueryStrategyAwardRuleModels(ctx context.Context, strategyID int64, awardID int) (string, error) {
	var awardPO po.StrategyAward
	err := r.db.WithContext(ctx).
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
	err := r.db.WithContext(ctx).
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

func (r *StrategyRepository) QueryRuleTreeByTreeID(ctx context.Context, treeID string) (treepkg.RuleTree, bool, error) {
	var ruleTreePO po.RuleTree
	err := r.db.WithContext(ctx).
		Select("tree_id", "tree_name", "tree_desc", "tree_node_rule_key").
		Where("tree_id = ?", treeID).
		First(&ruleTreePO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return treepkg.RuleTree{}, false, nil
	}
	if err != nil {
		return treepkg.RuleTree{}, false, err
	}

	var nodePOList []po.RuleTreeNode
	if err := r.db.WithContext(ctx).
		Select("tree_id", "rule_key", "rule_desc", "rule_value").
		Where("tree_id = ?", treeID).
		Find(&nodePOList).
		Error; err != nil {
		return treepkg.RuleTree{}, false, err
	}

	var linePOList []po.RuleTreeNodeLine
	if err := r.db.WithContext(ctx).
		Select("tree_id", "rule_node_from", "rule_node_to", "rule_limit_type", "rule_limit_value").
		Where("tree_id = ?", treeID).
		Find(&linePOList).
		Error; err != nil {
		return treepkg.RuleTree{}, false, err
	}

	lineMap := make(map[string][]treepkg.RuleTreeNodeLine)
	for _, linePO := range linePOList {
		lineMap[linePO.RuleNodeFrom] = append(lineMap[linePO.RuleNodeFrom], treepkg.RuleTreeNodeLine{
			TreeID:         linePO.TreeID,
			RuleNodeFrom:   linePO.RuleNodeFrom,
			RuleNodeTo:     linePO.RuleNodeTo,
			RuleLimitType:  linePO.RuleLimitType,
			RuleLimitValue: linePO.RuleLimitValue,
		})
	}

	nodeMap := make(map[string]treepkg.RuleTreeNode, len(nodePOList))
	for _, nodePO := range nodePOList {
		nodeMap[nodePO.RuleKey] = treepkg.RuleTreeNode{
			TreeID:    nodePO.TreeID,
			RuleKey:   nodePO.RuleKey,
			RuleDesc:  nodePO.RuleDesc,
			RuleValue: nodePO.RuleValue,
			Lines:     lineMap[nodePO.RuleKey],
		}
	}

	return treepkg.RuleTree{
		TreeID:   ruleTreePO.TreeID,
		TreeName: ruleTreePO.TreeName,
		TreeDesc: ruleTreePO.TreeDesc,
		RootRule: ruleTreePO.TreeRootRuleKey,
		NodeMap:  nodeMap,
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

func (r *StrategyRepository) QueryStrategyIDByActivityID(ctx context.Context, activityID int64) (int64, error) {
	var activityPO po.RaffleActivity
	err := r.db.WithContext(ctx).
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
	err := r.db.WithContext(ctx).
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

func (r *StrategyRepository) QueryRaffleActivityAccountDayPartakeCount(ctx context.Context, activityID int64, userID string) (int, error) {
	var accountDayPO po.RaffleActivityAccountDay
	err := r.db.WithContext(ctx).
		Select("day_count", "day_count_surplus").
		Where("user_id = ? and activity_id = ? and day = ?", userID, activityID, time.Now().Format("2006-01-02")).
		First(&accountDayPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return accountDayPO.DayCount - accountDayPO.DayCountSurplus, nil
}

func (r *StrategyRepository) QueryTodayUserRaffleCount(ctx context.Context, userID string, strategyID int64) (int, error) {
	activityID, err := r.queryActivityIDByStrategyID(ctx, strategyID)
	if err != nil {
		return 0, err
	}
	if activityID == 0 {
		return 0, nil
	}
	return r.QueryRaffleActivityAccountDayPartakeCount(ctx, activityID, userID)
}

func (r *StrategyRepository) AwardStockConsumeSendQueue(ctx context.Context, strategyID int64, awardID int) error {
	if r.stockQueue != nil {
		return r.stockQueue.AwardStockConsumeSendQueue(ctx, strategyID, awardID)
	}
	return nil
}

func (r *StrategyRepository) UpdateStrategyAwardStock(ctx context.Context, strategyID int64, awardID int) error {
	return r.db.WithContext(ctx).
		Model(&po.StrategyAward{}).
		Where("strategy_id = ? and award_id = ? and award_count_surplus > 0", strategyID, awardID).
		UpdateColumn("award_count_surplus", gorm.Expr("award_count_surplus - ?", 1)).
		Error
}

func (r *StrategyRepository) queryActivityIDByStrategyID(ctx context.Context, strategyID int64) (int64, error) {
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
	return activityPO.ActivityID, nil
}

func (r *StrategyRepository) QueryRaffleActivityAccountPartakeCount(ctx context.Context, activityID int64, userID string) (int, error) {
	var accountPO po.RaffleActivityAccount
	err := r.db.WithContext(ctx).
		Select("total_count", "total_count_surplus").
		Where("user_id = ? and activity_id = ?", userID, activityID).
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
	err := r.db.WithContext(ctx).
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
