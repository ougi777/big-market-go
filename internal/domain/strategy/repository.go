package strategy

import (
	"context"

	"bm-go/internal/domain/strategy/rule/tree"
)

type Repository interface {
	QueryStrategyEntityByStrategyID(ctx context.Context, strategyID int64) (StrategyEntity, error)
	QueryStrategyRuleValue(ctx context.Context, strategyID int64, ruleModel string) (string, error)
	QueryActivityAccountTotalUseCount(ctx context.Context, userID string, strategyID int64) (int, error)
}

type ArmoryRepository interface {
	QueryStrategyIDByActivityID(ctx context.Context, activityID int64) (int64, error)
	QueryStrategyAwardList(ctx context.Context, strategyID int64) ([]StrategyAwardEntity, error)
	QueryStrategyEntityByStrategyID(ctx context.Context, strategyID int64) (StrategyEntity, error)
	QueryStrategyRule(ctx context.Context, strategyID int64, ruleModel string) (StrategyRuleEntity, bool, error)
}

type QueryRepository interface {
	QueryStrategyIDByActivityID(ctx context.Context, activityID int64) (int64, error)
	QueryStrategyAwardList(ctx context.Context, strategyID int64) ([]StrategyAwardEntity, error)
	QueryAwardRuleLockCount(ctx context.Context, treeIDs []string) (map[string]int, error)
	QueryRaffleActivityAccountDayPartakeCount(ctx context.Context, activityID int64, userID string) (int, error)
	QueryRaffleActivityAccountPartakeCount(ctx context.Context, activityID int64, userID string) (int, error)
	QueryAwardRuleWeight(ctx context.Context, strategyID int64) ([]RuleWeight, error)
}

type RateTableStore interface {
	StoreStrategyAwardSearchRateTable(ctx context.Context, key string, rateRange int, table map[int]int) error
	CacheStrategyAwardCount(ctx context.Context, key string, awardCount int) error
}

type Dispatch interface {
	GetRandomAwardID(ctx context.Context, strategyID int64, ruleWeightValue ...string) (int, error)
}

type RaffleRepository interface {
	QueryStrategyAwardRuleModels(ctx context.Context, strategyID int64, awardID int) (string, error)
	QueryRuleTreeByTreeID(ctx context.Context, treeID string) (tree.RuleTree, bool, error)
	QueryStrategyAwardEntity(ctx context.Context, strategyID int64, awardID int) (StrategyAwardEntity, bool, error)
}

type StockRepository interface {
	UpdateStrategyAwardStock(ctx context.Context, strategyID int64, awardID int) error
}

type StockQueue interface {
	TakeQueueValue(ctx context.Context) (AwardStockKey, bool, error)
}
