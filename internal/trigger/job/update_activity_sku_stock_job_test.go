package job

import (
	"context"
	"errors"
	"testing"
)

func TestUpdateActivitySkuStockJobExec(t *testing.T) {
	updater := &fakeActivitySkuStockUpdater{}
	job := NewUpdateActivitySkuStockJob(updater, nil)

	job.Exec()

	if updater.calls != 1 {
		t.Fatalf("expected 1 call, got %d", updater.calls)
	}
}

func TestUpdateActivitySkuStockJobExecError(t *testing.T) {
	updater := &fakeActivitySkuStockUpdater{err: errors.New("update failed")}
	job := NewUpdateActivitySkuStockJob(updater, nil)

	job.Exec()

	if updater.calls != 1 {
		t.Fatalf("expected 1 call, got %d", updater.calls)
	}
}

type fakeActivitySkuStockUpdater struct {
	calls int
	err   error
}

func (f *fakeActivitySkuStockUpdater) UpdateActivitySkuStock(ctx context.Context) (bool, error) {
	f.calls++
	return true, f.err
}
