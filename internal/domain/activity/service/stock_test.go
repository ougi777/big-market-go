package service

import (
	"context"
	"errors"
	"testing"

	"bm-go/internal/domain/activity"
)

func TestStockServiceUpdateActivitySkuStock(t *testing.T) {
	repo := &fakeActivityStockRepository{}
	queue := &fakeActivityStockQueue{
		key: activity.ActivitySkuStockKey{SKU: 9011, ActivityID: 100301},
		ok:  true,
	}
	service := NewStockService(repo, queue, nil, nil)

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

func TestStockServiceUpdateActivitySkuStockEmptyQueue(t *testing.T) {
	repo := &fakeActivityStockRepository{}
	queue := &fakeActivityStockQueue{}
	service := NewStockService(repo, queue, nil, nil)

	updated, err := service.UpdateActivitySkuStock(context.Background())
	if err != nil {
		t.Fatalf("update activity sku stock: %v", err)
	}
	if updated {
		t.Fatal("expected no update")
	}
	if repo.updatedSKU != 0 {
		t.Fatalf("expected no repo update, got %d", repo.updatedSKU)
	}
}

func TestStockServiceUpdateActivitySkuStockRepositoryError(t *testing.T) {
	repo := &fakeActivityStockRepository{updateErr: errors.New("update failed")}
	queue := &fakeActivityStockQueue{
		key: activity.ActivitySkuStockKey{SKU: 9011, ActivityID: 100301},
		ok:  true,
	}
	service := NewStockService(repo, queue, nil, nil)

	_, err := service.UpdateActivitySkuStock(context.Background())
	if err == nil {
		t.Fatal("expected update error")
	}
}

func TestStockServiceClearActivitySkuStock(t *testing.T) {
	repo := &fakeActivityStockRepository{}
	queue := &fakeActivityStockQueue{}
	service := NewStockService(repo, queue, nil, nil)

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

func TestStockServiceClearActivitySkuStockRepositoryError(t *testing.T) {
	repo := &fakeActivityStockRepository{clearErr: errors.New("clear failed")}
	queue := &fakeActivityStockQueue{}
	service := NewStockService(repo, queue, nil, nil)

	err := service.ClearActivitySkuStock(context.Background(), 9011)
	if err == nil {
		t.Fatal("expected clear error")
	}
	if queue.cleared {
		t.Fatal("expected queue not cleared")
	}
}

func TestStockServiceSubtractActivitySkuStock(t *testing.T) {
	repo := &fakeActivityStockRepository{}
	queue := &fakeActivityStockQueue{}
	store := &fakeActivityStockStore{surplus: 8}
	service := NewStockService(repo, queue, store, nil)

	ok, err := service.SubtractActivitySkuStock(context.Background(), 9011, 100301)
	if err != nil {
		t.Fatalf("subtract activity sku stock: %v", err)
	}
	if !ok {
		t.Fatal("expected subtract ok")
	}
	if store.key == "" {
		t.Fatal("expected redis key recorded")
	}
	if queue.sent.SKU != 9011 || queue.sent.ActivityID != 100301 {
		t.Fatalf("expected stock queue item, got %+v", queue.sent)
	}
}

func TestStockServiceSubtractActivitySkuStockStoreError(t *testing.T) {
	store := &fakeActivityStockStore{subtractErr: errors.New("redis failed")}
	service := NewStockService(&fakeActivityStockRepository{}, &fakeActivityStockQueue{}, store, nil)

	ok, err := service.SubtractActivitySkuStock(context.Background(), 9011, 100301)
	if err == nil {
		t.Fatal("expected store error")
	}
	if ok {
		t.Fatal("expected subtract failed")
	}
}

func TestStockServiceSubtractActivitySkuStockSoldOut(t *testing.T) {
	queue := &fakeActivityStockQueue{}
	store := &fakeActivityStockStore{surplus: -1}
	service := NewStockService(&fakeActivityStockRepository{}, queue, store, nil)

	ok, err := service.SubtractActivitySkuStock(context.Background(), 9011, 100301)
	if err != nil {
		t.Fatalf("subtract activity sku stock: %v", err)
	}
	if ok {
		t.Fatal("expected sold out")
	}
	if queue.sent.SKU != 0 {
		t.Fatalf("expected no stock queue item, got %+v", queue.sent)
	}
}

func TestStockServiceSubtractActivitySkuStockPublishesStockZero(t *testing.T) {
	queue := &fakeActivityStockQueue{}
	store := &fakeActivityStockStore{surplus: 0}
	publisher := &fakeActivityStockPublisher{}
	service := NewStockService(&fakeActivityStockRepository{}, queue, store, publisher)

	ok, err := service.SubtractActivitySkuStock(context.Background(), 9011, 100301)
	if err != nil {
		t.Fatalf("subtract activity sku stock: %v", err)
	}
	if !ok {
		t.Fatal("expected subtract ok")
	}
	if publisher.topic != activity.TopicActivitySkuStockZero || publisher.message == "" {
		t.Fatalf("expected stock zero message, got topic=%s message=%s", publisher.topic, publisher.message)
	}
	if queue.sent.SKU != 9011 {
		t.Fatalf("expected stock queue item, got %+v", queue.sent)
	}
}

func TestStockServiceSubtractActivitySkuStockPublishError(t *testing.T) {
	queue := &fakeActivityStockQueue{}
	store := &fakeActivityStockStore{surplus: 0}
	publisher := &fakeActivityStockPublisher{err: errors.New("publish failed")}
	service := NewStockService(&fakeActivityStockRepository{}, queue, store, publisher)

	ok, err := service.SubtractActivitySkuStock(context.Background(), 9011, 100301)
	if err == nil {
		t.Fatal("expected publish error")
	}
	if ok {
		t.Fatal("expected subtract failed")
	}
	if queue.sent.SKU != 0 {
		t.Fatalf("expected no stock queue item, got %+v", queue.sent)
	}
}

type fakeActivityStockRepository struct {
	updatedSKU int64
	clearedSKU int64
	updateErr  error
	clearErr   error
}

func (f *fakeActivityStockRepository) UpdateActivitySkuStock(ctx context.Context, sku int64) error {
	f.updatedSKU = sku
	return f.updateErr
}

func (f *fakeActivityStockRepository) ClearActivitySkuStock(ctx context.Context, sku int64) error {
	f.clearedSKU = sku
	return f.clearErr
}

type fakeActivityStockQueue struct {
	key     activity.ActivitySkuStockKey
	ok      bool
	cleared bool
	sent    activity.ActivitySkuStockKey
}

func (f *fakeActivityStockQueue) SendActivitySkuStockConsumeQueue(ctx context.Context, stockKey activity.ActivitySkuStockKey) error {
	f.sent = stockKey
	return nil
}

func (f *fakeActivityStockQueue) TakeActivitySkuStock(ctx context.Context) (activity.ActivitySkuStockKey, bool, error) {
	return f.key, f.ok, nil
}

func (f *fakeActivityStockQueue) ClearActivitySkuStockQueue(ctx context.Context) error {
	f.cleared = true
	return nil
}

type fakeActivityStockStore struct {
	key         string
	surplus     int64
	subtractErr error
}

func (f *fakeActivityStockStore) CacheActivitySkuStockCount(ctx context.Context, key string, stockCount int) error {
	return nil
}

func (f *fakeActivityStockStore) SubtractActivitySkuStock(ctx context.Context, key string) (int64, error) {
	f.key = key
	return f.surplus, f.subtractErr
}

type fakeActivityStockPublisher struct {
	topic   string
	message string
	err     error
}

func (f *fakeActivityStockPublisher) Publish(ctx context.Context, topic string, message string) error {
	f.topic = topic
	f.message = message
	return f.err
}
