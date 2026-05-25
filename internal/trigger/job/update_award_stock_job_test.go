package job

import (
	"context"
	"testing"
)

func TestUpdateAwardStockJobExec(t *testing.T) {
	updater := &fakeAwardStockUpdater{}
	job := NewUpdateAwardStockJob(updater, nil)

	job.Exec()

	if updater.calls != 1 {
		t.Fatalf("expected 1 call, got %d", updater.calls)
	}
}

type fakeAwardStockUpdater struct {
	calls int
}

func (f *fakeAwardStockUpdater) UpdateAwardStock(ctx context.Context) (bool, error) {
	f.calls++
	return true, nil
}
