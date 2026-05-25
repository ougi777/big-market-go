package service

import (
	"context"
	"errors"
	"testing"

	"bm-go/internal/domain/task"
)

func TestServiceSendNoSendMessageTasks(t *testing.T) {
	repo := &fakeTaskRepository{
		tasks: []task.Entity{
			{
				UserID:    "xiaofuge",
				Topic:     "send_award",
				MessageID: "12345678901",
				Message:   `{"id":"12345678901"}`,
				State:     task.StateCreate,
			},
		},
	}
	publisher := &fakeTaskPublisher{}
	service := NewService(repo, publisher)

	err := service.SendNoSendMessageTasks(context.Background(), 10)
	if err != nil {
		t.Fatalf("send no send message tasks: %v", err)
	}

	if publisher.topic != "send_award" || publisher.message != `{"id":"12345678901"}` {
		t.Fatalf("expected publisher called, got %s/%s", publisher.topic, publisher.message)
	}
	if repo.completedUserID != "xiaofuge" || repo.completedMessageID != "12345678901" {
		t.Fatalf("expected task completed, got %s/%s", repo.completedUserID, repo.completedMessageID)
	}
}

func TestServiceSendNoSendMessageTasksMarksFailOnPublishError(t *testing.T) {
	repo := &fakeTaskRepository{
		tasks: []task.Entity{
			{
				UserID:    "xiaofuge",
				Topic:     "send_award",
				MessageID: "12345678901",
				Message:   `{"id":"12345678901"}`,
				State:     task.StateCreate,
			},
		},
	}
	publisher := &fakeTaskPublisher{err: errors.New("publish failed")}
	service := NewService(repo, publisher)

	err := service.SendNoSendMessageTasks(context.Background(), 10)
	if err != nil {
		t.Fatalf("send no send message tasks: %v", err)
	}

	if repo.failUserID != "xiaofuge" || repo.failMessageID != "12345678901" {
		t.Fatalf("expected task fail, got %s/%s", repo.failUserID, repo.failMessageID)
	}
}

func TestServiceSendNoSendMessageTasksQueryError(t *testing.T) {
	repo := &fakeTaskRepository{queryErr: errors.New("query failed")}
	service := NewService(repo, &fakeTaskPublisher{})

	err := service.SendNoSendMessageTasks(context.Background(), 10)
	if err == nil {
		t.Fatal("expected query error")
	}
}

func TestServiceSendNoSendMessageTasksFailStateError(t *testing.T) {
	repo := &fakeTaskRepository{
		tasks: []task.Entity{
			{
				UserID:    "xiaofuge",
				Topic:     "send_award",
				MessageID: "12345678901",
				Message:   `{"id":"12345678901"}`,
				State:     task.StateCreate,
			},
		},
		failErr: errors.New("mark fail failed"),
	}
	publisher := &fakeTaskPublisher{err: errors.New("publish failed")}
	service := NewService(repo, publisher)

	err := service.SendNoSendMessageTasks(context.Background(), 10)
	if err == nil {
		t.Fatal("expected fail state error")
	}
}

func TestServiceSendNoSendMessageTasksCompletedStateError(t *testing.T) {
	repo := &fakeTaskRepository{
		tasks: []task.Entity{
			{
				UserID:    "xiaofuge",
				Topic:     "send_award",
				MessageID: "12345678901",
				Message:   `{"id":"12345678901"}`,
				State:     task.StateCreate,
			},
		},
		completedErr: errors.New("mark completed failed"),
	}
	service := NewService(repo, &fakeTaskPublisher{})

	err := service.SendNoSendMessageTasks(context.Background(), 10)
	if err == nil {
		t.Fatal("expected completed state error")
	}
}

type fakeTaskRepository struct {
	tasks              []task.Entity
	completedUserID    string
	completedMessageID string
	failUserID         string
	failMessageID      string
	queryErr           error
	completedErr       error
	failErr            error
}

func (f *fakeTaskRepository) QueryNoSendMessageTaskList(ctx context.Context, limit int) ([]task.Entity, error) {
	return f.tasks, f.queryErr
}

func (f *fakeTaskRepository) UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error {
	f.completedUserID = userID
	f.completedMessageID = messageID
	return f.completedErr
}

func (f *fakeTaskRepository) UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error {
	f.failUserID = userID
	f.failMessageID = messageID
	return f.failErr
}

type fakeTaskPublisher struct {
	topic   string
	message string
	err     error
}

func (f *fakeTaskPublisher) Publish(ctx context.Context, topic string, message string) error {
	f.topic = topic
	f.message = message
	if f.err != nil {
		return f.err
	}
	return nil
}
