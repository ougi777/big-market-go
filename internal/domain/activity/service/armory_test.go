package service

import (
	"context"
	"errors"
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
			{
				SKU:               9012,
				ActivityID:        100301,
				ActivityCountID:   11102,
				StockCountSurplus: 100,
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
	secondCacheKey := types.RedisKeyActivitySkuStockCount + "9012"
	if store.stocks[secondCacheKey] != 100 {
		t.Fatalf("expected sku stock 100, got %d", store.stocks[secondCacheKey])
	}
	if repo.activityID != 100301 {
		t.Fatalf("expected activity id 100301, got %d", repo.activityID)
	}
}

func TestArmoryServiceAssembleActivitySkuRepositoryError(t *testing.T) {
	repo := &fakeArmorySkuProductRepository{err: errors.New("query failed")}
	store := &fakeActivitySkuStockStore{stocks: make(map[string]int)}
	armory := NewArmoryService(repo, store)

	err := armory.AssembleActivitySkuByActivityID(context.Background(), 100301)
	if err == nil {
		t.Fatal("expected repository error")
	}
	if len(store.stocks) != 0 {
		t.Fatalf("expected no cache, got %+v", store.stocks)
	}
}

func TestArmoryServiceAssembleActivitySkuCacheError(t *testing.T) {
	repo := &fakeArmorySkuProductRepository{
		products: []activity.SkuProductEntity{
			{SKU: 9011, ActivityID: 100301, StockCountSurplus: 99890},
		},
	}
	store := &fakeActivitySkuStockStore{
		stocks: make(map[string]int),
		err:    errors.New("cache failed"),
	}
	armory := NewArmoryService(repo, store)

	err := armory.AssembleActivitySkuByActivityID(context.Background(), 100301)
	if err == nil {
		t.Fatal("expected cache error")
	}
}

type fakeArmorySkuProductRepository struct {
	activityID int64
	products   []activity.SkuProductEntity
	err        error
}

func (f *fakeArmorySkuProductRepository) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	f.activityID = activityID
	return f.products, f.err
}

func (f *fakeArmorySkuProductRepository) QuerySkuProductBySKU(ctx context.Context, sku int64) (activity.SkuProductEntity, bool, error) {
	return activity.SkuProductEntity{}, false, nil
}

type fakeActivitySkuStockStore struct {
	stocks map[string]int
	err    error
}

func (f *fakeActivitySkuStockStore) CacheActivitySkuStockCount(ctx context.Context, key string, stockCount int) error {
	f.stocks[key] = stockCount
	return f.err
}

func (f *fakeActivitySkuStockStore) SubtractActivitySkuStock(ctx context.Context, key string) (int64, error) {
	return 0, nil
}
