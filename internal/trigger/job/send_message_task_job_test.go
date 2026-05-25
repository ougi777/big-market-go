package job

import (
	"context"
	"testing"
)

func TestSendMessageTaskJobExec(t *testing.T) {
	sender := &fakeMessageTaskSender{}
	job := NewSendMessageTaskJob(sender, nil)

	job.Exec()

	if sender.calls != 1 {
		t.Fatalf("expected 1 call, got %d", sender.calls)
	}
	if sender.limit != 10 {
		t.Fatalf("expected limit 10, got %d", sender.limit)
	}
}

type fakeMessageTaskSender struct {
	calls int
	limit int
}

func (f *fakeMessageTaskSender) SendNoSendMessageTasks(ctx context.Context, limit int) error {
	f.calls++
	f.limit = limit
	return nil
}
