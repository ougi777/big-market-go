package strategy

import "context"

type Repository interface {
	QueryStrategyEntityByStrategyID(ctx context.Context, strategyID int64) (StrategyEntity, error)
	QueryStrategyRuleValue(ctx context.Context, strategyID int64, ruleModel string) (string, error)
	QueryActivityAccountTotalUseCount(ctx context.Context, userID string, strategyID int64) (int, error)
}

type ArmoryRepository interface {
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
