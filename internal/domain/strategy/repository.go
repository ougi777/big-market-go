package strategy

import "context"

type Repository interface {
	QueryStrategyEntityByStrategyID(ctx context.Context, strategyID int64) (StrategyEntity, error)
	QueryStrategyRuleValue(ctx context.Context, strategyID int64, ruleModel string) (string, error)
	QueryActivityAccountTotalUseCount(ctx context.Context, userID string, strategyID int64) (int, error)
}

type Dispatch interface {
	GetRandomAwardID(ctx context.Context, strategyID int64, ruleWeightValue ...string) (int, error)
}
