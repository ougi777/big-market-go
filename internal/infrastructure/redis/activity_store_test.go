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

func TestParseActivitySkuStockQueueValue(t *testing.T) {
	key, err := parseActivitySkuStockQueueValue(`{"sku":9011,"activityId":100301}`)
	if err != nil {
		t.Fatalf("expected parse success, got %v", err)
	}
	if key.SKU != 9011 || key.ActivityID != 100301 {
		t.Fatalf("unexpected key: %+v", key)
	}
}

func TestParseActivitySkuStockQueueValueInvalid(t *testing.T) {
	if _, err := parseActivitySkuStockQueueValue(`{"sku":`); err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestActivityStoreClearActivitySkuStockQueueWithoutClient(t *testing.T) {
	store := NewActivityStore(nil)

	err := store.ClearActivitySkuStockQueue(context.Background())
	if !errors.Is(err, ErrClientNotConnected) {
		t.Fatalf("expected not connected error, got %v", err)
	}
}
