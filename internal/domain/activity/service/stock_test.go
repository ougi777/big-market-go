package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/activity"
)

func TestStockServiceUpdateActivitySkuStock(t *testing.T) {
	repo := &fakeActivityStockRepository{}
	queue := &fakeActivityStockQueue{
		key: activity.ActivitySkuStockKey{SKU: 9011, ActivityID: 100301},
		ok:  true,
	}
	service := NewStockService(repo, queue)

	updated, err := service.UpdateActivitySkuStock(context.Background())
	if err != nil {
		t.Fatalf("update activity sku stock: %v", err)
	}
	if !updated {
		t.Fatal("expected updated")
	}
	if repo.updatedSKU != 9011 {
		t.Fatalf("expected sku 9011, got %d", repo.updatedSKU)
	}
}

func TestStockServiceClearActivitySkuStock(t *testing.T) {
	repo := &fakeActivityStockRepository{}
	queue := &fakeActivityStockQueue{}
	service := NewStockService(repo, queue)

	err := service.ClearActivitySkuStock(context.Background(), 9011)
	if err != nil {
		t.Fatalf("clear activity sku stock: %v", err)
	}
	if repo.clearedSKU != 9011 {
		t.Fatalf("expected cleared sku 9011, got %d", repo.clearedSKU)
	}
	if !queue.cleared {
		t.Fatal("expected queue cleared")
	}
}

type fakeActivityStockRepository struct {
	updatedSKU int64
	clearedSKU int64
}

func (f *fakeActivityStockRepository) UpdateActivitySkuStock(ctx context.Context, sku int64) error {
	f.updatedSKU = sku
	return nil
}

func (f *fakeActivityStockRepository) ClearActivitySkuStock(ctx context.Context, sku int64) error {
	f.clearedSKU = sku
	return nil
}

type fakeActivityStockQueue struct {
	key     activity.ActivitySkuStockKey
	ok      bool
	cleared bool
}

func (f *fakeActivityStockQueue) TakeActivitySkuStock(ctx context.Context) (activity.ActivitySkuStockKey, bool, error) {
	return f.key, f.ok, nil
}

func (f *fakeActivityStockQueue) ClearActivitySkuStockQueue(ctx context.Context) error {
	f.cleared = true
	return nil
}
