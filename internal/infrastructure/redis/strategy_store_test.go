package redis

import (
	"context"
	"errors"
	"testing"
)

func TestStrategyStoreTakeQueueValueWithoutClient(t *testing.T) {
	var store *StrategyStore

	_, _, err := store.TakeQueueValue(context.Background())
	if !errors.Is(err, ErrClientNotConnected) {
		t.Fatalf("expected not connected error, got %v", err)
	}
}

func TestStrategyStoreStoreStrategyAwardSearchRateTableWithoutClient(t *testing.T) {
	store := NewStrategyStore(nil)

	err := store.StoreStrategyAwardSearchRateTable(context.Background(), "100001", 1, map[int]int{0: 101})
	if !errors.Is(err, ErrClientNotConnected) {
		t.Fatalf("expected not connected error, got %v", err)
	}
}

func TestStrategyStoreCacheStrategyAwardCountWithoutClient(t *testing.T) {
	store := NewStrategyStore(nil)

	err := store.CacheStrategyAwardCount(context.Background(), "strategy_award_count_key_100001_101", 10)
	if !errors.Is(err, ErrClientNotConnected) {
		t.Fatalf("expected not connected error, got %v", err)
	}
}

func TestParseAwardStockQueueValue(t *testing.T) {
	key, err := parseAwardStockQueueValue("100001:101")
	if err != nil {
		t.Fatalf("expected parse success, got %v", err)
	}
	if key.StrategyID != 100001 || key.AwardID != 101 {
		t.Fatalf("unexpected key: %+v", key)
	}
}

func TestParseAwardStockQueueValueInvalid(t *testing.T) {
	cases := []string{"100001", "abc:101", "100001:abc"}

	for _, value := range cases {
		if _, err := parseAwardStockQueueValue(value); err == nil {
			t.Fatalf("expected parse error for %q", value)
		}
	}
}

func TestStrategyStoreAwardStockConsumeSendQueueWithoutClient(t *testing.T) {
	store := NewStrategyStore(nil)

	err := store.AwardStockConsumeSendQueue(context.Background(), 100001, 101)
	if !errors.Is(err, ErrClientNotConnected) {
		t.Fatalf("expected not connected error, got %v", err)
	}
}
