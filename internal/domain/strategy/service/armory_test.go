package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/types"
)

func TestArmoryServiceAssembleLotteryStrategy(t *testing.T) {
	repo := &fakeArmoryRepository{
		awards: []strategy.StrategyAwardEntity{
			{StrategyID: 100001, AwardID: 101, AwardCountSurplus: 10, AwardRate: 0.1},
			{StrategyID: 100001, AwardID: 102, AwardCountSurplus: 20, AwardRate: 0.2},
		},
		strategyEntity: strategy.StrategyEntity{
			StrategyID:    100001,
			RuleModelList: []string{"rule_weight"},
		},
		ruleEntity: strategy.StrategyRuleEntity{
			StrategyID: 100001,
			RuleModel:  "rule_weight",
			RuleValue:  "4000:101",
		},
	}
	store := newFakeRateTableStore()
	armory := NewArmoryService(repo, store)

	if err := armory.AssembleLotteryStrategy(context.Background(), 100001); err != nil {
		t.Fatalf("assemble lottery strategy: %v", err)
	}

	defaultKey := "100001"
	if store.rateRanges[defaultKey] != 3 {
		t.Fatalf("expected default rate range 3, got %d", store.rateRanges[defaultKey])
	}
	if len(store.tables[defaultKey]) != 3 {
		t.Fatalf("expected default rate table size 3, got %d", len(store.tables[defaultKey]))
	}

	weightKey := "100001_4000:101"
	if store.rateRanges[weightKey] != 1 {
		t.Fatalf("expected weight rate range 1, got %d", store.rateRanges[weightKey])
	}
	if len(store.tables[weightKey]) != 1 {
		t.Fatalf("expected weight rate table size 1, got %d", len(store.tables[weightKey]))
	}

	award101StockKey := types.RedisKeyStrategyAwardCount + "100001_101"
	if store.awardCounts[award101StockKey] != 10 {
		t.Fatalf("expected award 101 stock 10, got %d", store.awardCounts[award101StockKey])
	}
	award102StockKey := types.RedisKeyStrategyAwardCount + "100001_102"
	if store.awardCounts[award102StockKey] != 20 {
		t.Fatalf("expected award 102 stock 20, got %d", store.awardCounts[award102StockKey])
	}
}

type fakeArmoryRepository struct {
	awards         []strategy.StrategyAwardEntity
	strategyEntity strategy.StrategyEntity
	ruleEntity     strategy.StrategyRuleEntity
}

func (f *fakeArmoryRepository) QueryStrategyAwardList(ctx context.Context, strategyID int64) ([]strategy.StrategyAwardEntity, error) {
	return f.awards, nil
}

func (f *fakeArmoryRepository) QueryStrategyEntityByStrategyID(ctx context.Context, strategyID int64) (strategy.StrategyEntity, error) {
	return f.strategyEntity, nil
}

func (f *fakeArmoryRepository) QueryStrategyRule(ctx context.Context, strategyID int64, ruleModel string) (strategy.StrategyRuleEntity, bool, error) {
	return f.ruleEntity, true, nil
}

type fakeRateTableStore struct {
	rateRanges  map[string]int
	tables      map[string]map[int]int
	awardCounts map[string]int
}

func newFakeRateTableStore() *fakeRateTableStore {
	return &fakeRateTableStore{
		rateRanges:  make(map[string]int),
		tables:      make(map[string]map[int]int),
		awardCounts: make(map[string]int),
	}
}

func (f *fakeRateTableStore) StoreStrategyAwardSearchRateTable(ctx context.Context, key string, rateRange int, table map[int]int) error {
	f.rateRanges[key] = rateRange
	f.tables[key] = table
	return nil
}

func (f *fakeRateTableStore) CacheStrategyAwardCount(ctx context.Context, key string, awardCount int) error {
	f.awardCounts[key] = awardCount
	return nil
}
