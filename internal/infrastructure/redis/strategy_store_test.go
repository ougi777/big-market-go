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

func TestStrategyStoreAwardStockConsumeSendQueueWithoutClient(t *testing.T) {
	store := NewStrategyStore(nil)

	err := store.AwardStockConsumeSendQueue(context.Background(), 100001, 101)
	if !errors.Is(err, ErrClientNotConnected) {
		t.Fatalf("expected not connected error, got %v", err)
	}
}
