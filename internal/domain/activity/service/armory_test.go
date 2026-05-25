package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"
)

func TestArmoryServiceAssembleActivitySkuByActivityID(t *testing.T) {
	repo := &fakeArmorySkuProductRepository{
		products: []activity.SkuProductEntity{
			{
				SKU:               9011,
				ActivityID:        100301,
				ActivityCountID:   11101,
				StockCountSurplus: 99890,
			},
		},
	}
	store := &fakeActivitySkuStockStore{
		stocks: make(map[string]int),
	}
	armory := NewArmoryService(repo, store)

	if err := armory.AssembleActivitySkuByActivityID(context.Background(), 100301); err != nil {
		t.Fatalf("assemble activity sku: %v", err)
	}

	cacheKey := types.RedisKeyActivitySkuStockCount + "9011"
	if store.stocks[cacheKey] != 99890 {
		t.Fatalf("expected sku stock 99890, got %d", store.stocks[cacheKey])
	}
	if repo.activityID != 100301 {
		t.Fatalf("expected activity id 100301, got %d", repo.activityID)
	}
}

type fakeArmorySkuProductRepository struct {
	activityID int64
	products   []activity.SkuProductEntity
}

func (f *fakeArmorySkuProductRepository) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	f.activityID = activityID
	return f.products, nil
}

type fakeActivitySkuStockStore struct {
	stocks map[string]int
}

func (f *fakeActivitySkuStockStore) CacheActivitySkuStockCount(ctx context.Context, key string, stockCount int) error {
	f.stocks[key] = stockCount
	return nil
}
