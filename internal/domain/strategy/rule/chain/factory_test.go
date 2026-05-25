package chain

import (
	"context"
	"testing"

	"bm-go/internal/domain/strategy"
)

func TestFactoryBlacklistTakesOver(t *testing.T) {
	repo := &fakeRepository{
		strategy: strategy.StrategyEntity{StrategyID: 100001, RuleModelList: []string{RuleBlacklist, RuleWeight}},
		rules: map[string]string{
			RuleBlacklist: "101:user001,user002",
		},
	}
	dispatch := &fakeDispatch{defaultAwardID: 102}
	factory := NewFactory(repo, dispatch)

	logicChain, err := factory.OpenLogicChain(context.Background(), 100001)
	if err != nil {
		t.Fatalf("open logic chain: %v", err)
	}
	result, err := logicChain.Logic(context.Background(), "user001", 100001)
	if err != nil {
		t.Fatalf("logic: %v", err)
	}

	if result.AwardID != 101 {
		t.Fatalf("expected award 101, got %d", result.AwardID)
	}
	if result.LogicModel != RuleBlacklist {
		t.Fatalf("expected logic model %s, got %s", RuleBlacklist, result.LogicModel)
	}
}

func TestFactoryWeightTakesOver(t *testing.T) {
	repo := &fakeRepository{
		strategy: strategy.StrategyEntity{StrategyID: 100001, RuleModelList: []string{RuleWeight}},
		rules: map[string]string{
			RuleWeight: "4000:102,103 5000:104,105",
		},
		userScore: 5500,
	}
	dispatch := &fakeDispatch{weightedAwardID: 104}
	factory := NewFactory(repo, dispatch)

	logicChain, err := factory.OpenLogicChain(context.Background(), 100001)
	if err != nil {
		t.Fatalf("open logic chain: %v", err)
	}
	result, err := logicChain.Logic(context.Background(), "user003", 100001)
	if err != nil {
		t.Fatalf("logic: %v", err)
	}

	if result.AwardID != 104 {
		t.Fatalf("expected award 104, got %d", result.AwardID)
	}
	if result.LogicModel != RuleWeight {
		t.Fatalf("expected logic model %s, got %s", RuleWeight, result.LogicModel)
	}
	if dispatch.lastWeightValue != "5000:104,105" {
		t.Fatalf("expected weight value 5000:104,105, got %s", dispatch.lastWeightValue)
	}
}

func TestFactoryDefaultFallback(t *testing.T) {
	repo := &fakeRepository{
		strategy: strategy.StrategyEntity{StrategyID: 100001},
	}
	dispatch := &fakeDispatch{defaultAwardID: 108}
	factory := NewFactory(repo, dispatch)

	logicChain, err := factory.OpenLogicChain(context.Background(), 100001)
	if err != nil {
		t.Fatalf("open logic chain: %v", err)
	}
	result, err := logicChain.Logic(context.Background(), "user004", 100001)
	if err != nil {
		t.Fatalf("logic: %v", err)
	}

	if result.AwardID != 108 {
		t.Fatalf("expected award 108, got %d", result.AwardID)
	}
	if result.LogicModel != RuleDefault {
		t.Fatalf("expected logic model %s, got %s", RuleDefault, result.LogicModel)
	}
}

type fakeRepository struct {
	strategy  strategy.StrategyEntity
	rules     map[string]string
	userScore int
}

func (f *fakeRepository) QueryStrategyEntityByStrategyID(ctx context.Context, strategyID int64) (strategy.StrategyEntity, error) {
	return f.strategy, nil
}

func (f *fakeRepository) QueryStrategyRuleValue(ctx context.Context, strategyID int64, ruleModel string) (string, error) {
	return f.rules[ruleModel], nil
}

func (f *fakeRepository) QueryActivityAccountTotalUseCount(ctx context.Context, userID string, strategyID int64) (int, error) {
	return f.userScore, nil
}

type fakeDispatch struct {
	defaultAwardID  int
	weightedAwardID int
	lastWeightValue string
}

func (f *fakeDispatch) GetRandomAwardID(ctx context.Context, strategyID int64, ruleWeightValue ...string) (int, error) {
	if len(ruleWeightValue) > 0 {
		f.lastWeightValue = ruleWeightValue[0]
		return f.weightedAwardID, nil
	}
	return f.defaultAwardID, nil
}
