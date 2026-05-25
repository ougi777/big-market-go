package redis

import (
	"context"
	"errors"
	"testing"
)

func TestActivityStoreSubtractActivitySkuStockWithoutClient(t *testing.T) {
	var store *ActivityStore

	_, err := store.SubtractActivitySkuStock(context.Background(), "activity_sku_stock_count_key_9011")
	if !errors.Is(err, ErrClientNotConnected) {
		t.Fatalf("expected not connected error, got %v", err)
	}
}

func TestActivityStoreClearActivitySkuStockQueueWithoutClient(t *testing.T) {
	store := NewActivityStore(nil)

	err := store.ClearActivitySkuStockQueue(context.Background())
	if !errors.Is(err, ErrClientNotConnected) {
		t.Fatalf("expected not connected error, got %v", err)
	}
}
