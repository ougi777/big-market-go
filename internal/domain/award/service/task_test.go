package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/award"
)

func TestTaskServiceSendNoSendMessageTasks(t *testing.T) {
	repo := &fakeAwardTaskRepository{
		tasks: []award.TaskEntity{
			{
				UserID:    "xiaofuge",
				Topic:     award.TopicSendAward,
				MessageID: "12345678901",
				Message:   `{"id":"12345678901"}`,
				State:     award.TaskStateCreate,
			},
		},
	}
	publisher := &fakeAwardPublisher{}
	service := NewTaskService(repo, publisher)

	err := service.SendNoSendMessageTasks(context.Background(), 10)
	if err != nil {
		t.Fatalf("send no send message tasks: %v", err)
	}

	if publisher.topic != award.TopicSendAward || publisher.message != `{"id":"12345678901"}` {
		t.Fatalf("expected publisher called, got %s/%s", publisher.topic, publisher.message)
	}
	if repo.completedUserID != "xiaofuge" || repo.completedMessageID != "12345678901" {
		t.Fatalf("expected task completed, got %s/%s", repo.completedUserID, repo.completedMessageID)
	}
}
