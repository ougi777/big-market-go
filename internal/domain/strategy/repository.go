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

type RateTableStore interface {
	StoreStrategyAwardSearchRateTable(ctx context.Context, key string, rateRange int, table map[int]int) error
	CacheStrategyAwardCount(ctx context.Context, key string, awardCount int) error
}

type Dispatch interface {
	GetRandomAwardID(ctx context.Context, strategyID int64, ruleWeightValue ...string) (int, error)
}
