package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/strategy"
)

func TestQueryServiceQueryRaffleAwardList(t *testing.T) {
	repo := &fakeStrategyQueryRepository{
		strategyID:      100006,
		dayPartakeCount: 1,
		awards: []strategy.StrategyAwardEntity{
			{
				AwardID:       101,
				AwardTitle:    "积分",
				AwardSubtitle: "抽奖1次后解锁",
				Sort:          1,
				RuleModels:    "tree_lock_1",
			},
			{
				AwardID:    102,
				AwardTitle: "优惠券",
				Sort:       2,
			},
		},
		lockCounts: map[string]int{
			"tree_lock_1": 2,
		},
	}
	service := NewQueryService(repo)

	awards, err := service.QueryRaffleAwardList(context.Background(), 100301, "xiaofuge")
	if err != nil {
		t.Fatalf("query raffle award list: %v", err)
	}

	if len(awards) != 2 {
		t.Fatalf("expected 2 awards, got %d", len(awards))
	}
	if awards[0].IsAwardUnlock || awards[0].WaitUnlockCount != 1 {
		t.Fatalf("expected first award locked with wait count 1, got %+v", awards[0])
	}
	if !awards[1].IsAwardUnlock || awards[1].HasAwardRuleLock {
		t.Fatalf("expected second award unlocked without lock, got %+v", awards[1])
	}
}

func TestQueryServiceQueryRaffleStrategyRuleWeight(t *testing.T) {
	repo := &fakeStrategyQueryRepository{
		strategyID:    100006,
		totalUseCount: 4500,
		ruleWeightList: []strategy.RuleWeight{
			{
				Weight: 4000,
				AwardList: []strategy.RuleWeightAward{
					{AwardID: 101, AwardTitle: "积分"},
				},
			},
		},
	}
	service := NewQueryService(repo)

	ruleWeights, err := service.QueryRaffleStrategyRuleWeight(context.Background(), 100301, "xiaofuge")
	if err != nil {
		t.Fatalf("query raffle strategy rule weight: %v", err)
	}

	if len(ruleWeights) != 1 {
		t.Fatalf("expected 1 rule weight, got %d", len(ruleWeights))
	}
	if ruleWeights[0].UserActivityAccountTotalUseCount != 4500 || ruleWeights[0].RuleWeightCount != 4000 {
		t.Fatalf("expected count data, got %+v", ruleWeights[0])
	}
	if ruleWeights[0].StrategyAwards[0].AwardID != 101 {
		t.Fatalf("expected award 101, got %+v", ruleWeights[0].StrategyAwards)
	}
}

type fakeStrategyQueryRepository struct {
	strategyID      int64
	awards          []strategy.StrategyAwardEntity
	lockCounts      map[string]int
	dayPartakeCount int
	totalUseCount   int
	ruleWeightList  []strategy.RuleWeight
}

func (f *fakeStrategyQueryRepository) QueryStrategyIDByActivityID(ctx context.Context, activityID int64) (int64, error) {
	return f.strategyID, nil
}

func (f *fakeStrategyQueryRepository) QueryStrategyAwardList(ctx context.Context, strategyID int64) ([]strategy.StrategyAwardEntity, error) {
	return f.awards, nil
}

func (f *fakeStrategyQueryRepository) QueryAwardRuleLockCount(ctx context.Context, treeIDs []string) (map[string]int, error) {
	return f.lockCounts, nil
}

func (f *fakeStrategyQueryRepository) QueryRaffleActivityAccountDayPartakeCount(ctx context.Context, activityID int64, userID string) (int, error) {
	return f.dayPartakeCount, nil
}

func (f *fakeStrategyQueryRepository) QueryRaffleActivityAccountPartakeCount(ctx context.Context, activityID int64, userID string) (int, error) {
	return f.totalUseCount, nil
}

func (f *fakeStrategyQueryRepository) QueryAwardRuleWeight(ctx context.Context, strategyID int64) ([]strategy.RuleWeight, error) {
	return f.ruleWeightList, nil
}
