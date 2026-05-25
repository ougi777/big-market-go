package job

import (
	"context"
	"errors"
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

func TestUpdateAwardStockJobExecError(t *testing.T) {
	updater := &fakeAwardStockUpdater{err: errors.New("update failed")}
	job := NewUpdateAwardStockJob(updater, nil)

	job.Exec()

	if updater.calls != 1 {
		t.Fatalf("expected 1 call, got %d", updater.calls)
	}
}

type fakeAwardStockUpdater struct {
	calls int
	err   error
}

func (f *fakeAwardStockUpdater) UpdateAwardStock(ctx context.Context) (bool, error) {
	f.calls++
	return true, f.err
}
