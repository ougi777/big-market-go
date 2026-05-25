package job

import (
	"context"
	"errors"
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

func TestSendMessageTaskJobExecSenderError(t *testing.T) {
	sender := &fakeMessageTaskSender{err: errors.New("send failed")}
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
	err   error
}

func (f *fakeMessageTaskSender) SendNoSendMessageTasks(ctx context.Context, limit int) error {
	f.calls++
	f.limit = limit
	return f.err
}
