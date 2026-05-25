package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/strategy"
)

func TestStockServiceUpdateAwardStock(t *testing.T) {
	repo := &fakeStockRepository{}
	queue := &fakeStockQueue{
		value: strategy.AwardStockKey{StrategyID: 100001, AwardID: 101},
		ok:    true,
	}
	service := NewStockService(repo, queue)

	updated, err := service.UpdateAwardStock(context.Background())
	if err != nil {
		t.Fatalf("update award stock: %v", err)
	}
	if !updated {
		t.Fatal("expected updated")
	}
	if repo.strategyID != 100001 || repo.awardID != 101 {
		t.Fatalf("unexpected update target %d:%d", repo.strategyID, repo.awardID)
	}
}

func TestStockServiceEmptyQueue(t *testing.T) {
	service := NewStockService(&fakeStockRepository{}, &fakeStockQueue{})

	updated, err := service.UpdateAwardStock(context.Background())
	if err != nil {
		t.Fatalf("update award stock: %v", err)
	}
	if updated {
		t.Fatal("expected empty queue")
	}
}

type fakeStockRepository struct {
	strategyID int64
	awardID    int
}

func (f *fakeStockRepository) UpdateStrategyAwardStock(ctx context.Context, strategyID int64, awardID int) error {
	f.strategyID = strategyID
	f.awardID = awardID
	return nil
}

type fakeStockQueue struct {
	value strategy.AwardStockKey
	ok    bool
}

func (f *fakeStockQueue) TakeQueueValue(ctx context.Context) (strategy.AwardStockKey, bool, error) {
	return f.value, f.ok, nil
}
