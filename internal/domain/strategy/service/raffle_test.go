package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/domain/strategy/rule/chain"
	"bm-go/internal/domain/strategy/rule/tree"
)

func TestRaffleServicePerformRaffleWithTree(t *testing.T) {
	repo := &fakeRaffleRepository{}
	dispatch := &fakeRaffleDispatch{awardID: 101, stockOK: true}
	chainFactory := chain.NewFactory(repo, dispatch)
	service := NewRaffleService(chainFactory, repo, map[string]tree.Node{
		tree.RuleStock: tree.NewStockNode(repo, dispatch),
	})

	result, err := service.PerformRaffle(context.Background(), "user001", 100001)
	if err != nil {
		t.Fatalf("perform raffle: %v", err)
	}
	if result.AwardID != 101 {
		t.Fatalf("expected award 101, got %d", result.AwardID)
	}
	if result.AwardIndex != 1 {
		t.Fatalf("expected award index 1, got %d", result.AwardIndex)
	}
	if repo.queuedAwardID != 101 {
		t.Fatalf("expected queued award 101, got %d", repo.queuedAwardID)
	}
}

type fakeRaffleRepository struct {
	queuedAwardID int
}

func (f *fakeRaffleRepository) QueryStrategyEntityByStrategyID(ctx context.Context, strategyID int64) (strategy.StrategyEntity, error) {
	return strategy.StrategyEntity{StrategyID: strategyID}, nil
}

func (f *fakeRaffleRepository) QueryStrategyRuleValue(ctx context.Context, strategyID int64, ruleModel string) (string, error) {
	return "", nil
}

func (f *fakeRaffleRepository) QueryActivityAccountTotalUseCount(ctx context.Context, userID string, strategyID int64) (int, error) {
	return 0, nil
}

func (f *fakeRaffleRepository) QueryStrategyAwardRuleModels(ctx context.Context, strategyID int64, awardID int) (string, error) {
	return "tree_stock", nil
}

func (f *fakeRaffleRepository) QueryRuleTreeByTreeID(ctx context.Context, treeID string) (tree.RuleTree, bool, error) {
	return tree.RuleTree{
		TreeID:   treeID,
		RootRule: tree.RuleStock,
		NodeMap: map[string]tree.RuleTreeNode{
			tree.RuleStock: {
				RuleKey:   tree.RuleStock,
				RuleValue: "stock",
			},
		},
	}, true, nil
}

func (f *fakeRaffleRepository) QueryStrategyAwardEntity(ctx context.Context, strategyID int64, awardID int) (strategy.StrategyAwardEntity, bool, error) {
	return strategy.StrategyAwardEntity{
		StrategyID: strategyID,
		AwardID:    awardID,
		AwardTitle: "积分",
		Sort:       1,
		RuleModels: "tree_stock",
	}, true, nil
}

func (f *fakeRaffleRepository) QueryTodayUserRaffleCount(ctx context.Context, userID string, strategyID int64) (int, error) {
	return 0, nil
}

func (f *fakeRaffleRepository) AwardStockConsumeSendQueue(ctx context.Context, strategyID int64, awardID int) error {
	f.queuedAwardID = awardID
	return nil
}

type fakeRaffleDispatch struct {
	awardID int
	stockOK bool
}

func (f *fakeRaffleDispatch) GetRandomAwardID(ctx context.Context, strategyID int64, ruleWeightValue ...string) (int, error) {
	return f.awardID, nil
}

func (f *fakeRaffleDispatch) SubtractionAwardStock(ctx context.Context, strategyID int64, awardID int) (bool, error) {
	return f.stockOK, nil
}
