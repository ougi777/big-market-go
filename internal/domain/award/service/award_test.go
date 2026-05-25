package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"bm-go/internal/domain/award"
)

func TestAwardServiceSaveUserAwardRecord(t *testing.T) {
	repo := &fakeAwardRepository{}
	taskRepo := &fakeAwardTaskRepository{}
	publisher := &fakeAwardPublisher{}
	service := NewAwardService(repo, taskRepo, publisher)
	service.now = func() time.Time { return time.Date(2026, 5, 25, 10, 0, 0, 0, time.UTC) }
	service.messageGenerator = func() (string, error) { return "12345678901", nil }

	err := service.SaveUserAwardRecord(context.Background(), award.UserAwardRecordEntity{
		UserID:      "xiaofuge",
		ActivityID:  100301,
		StrategyID:  100006,
		OrderID:     "order-001",
		AwardID:     101,
		AwardTitle:  "随机积分",
		AwardConfig: "0.01,1",
		AwardState:  award.AwardStateCreate,
	})
	if err != nil {
		t.Fatalf("save award record: %v", err)
	}

	if repo.record.SendTask.Topic != award.TopicSendAward {
		t.Fatalf("expected topic %s, got %s", award.TopicSendAward, repo.record.SendTask.Topic)
	}
	if repo.record.SendTask.MessageID != "12345678901" {
		t.Fatalf("expected message id, got %s", repo.record.SendTask.MessageID)
	}
	if !strings.Contains(repo.record.SendTask.Message, `"awardId":101`) ||
		!strings.Contains(repo.record.SendTask.Message, `"orderId":"order-001"`) ||
		!strings.Contains(repo.record.SendTask.Message, `"timestamp":1779703200000`) {
		t.Fatalf("expected award message, got %s", repo.record.SendTask.Message)
	}
	if publisher.topic != award.TopicSendAward || publisher.message != repo.record.SendTask.Message {
		t.Fatalf("expected publisher called, got %s/%s", publisher.topic, publisher.message)
	}
	if taskRepo.completedUserID != "xiaofuge" || taskRepo.completedMessageID != "12345678901" {
		t.Fatalf("expected task completed, got %s/%s", taskRepo.completedUserID, taskRepo.completedMessageID)
	}
}

type fakeAwardRepository struct {
	record award.UserAwardRecordEntity
}

func (f *fakeAwardRepository) SaveUserAwardRecord(ctx context.Context, record award.UserAwardRecordEntity) error {
	f.record = record
	return nil
}

type fakeAwardTaskRepository struct {
	tasks              []award.TaskEntity
	completedUserID    string
	completedMessageID string
	failUserID         string
	failMessageID      string
}

func (f *fakeAwardTaskRepository) QueryNoSendMessageTaskList(ctx context.Context, limit int) ([]award.TaskEntity, error) {
	return f.tasks, nil
}

func (f *fakeAwardTaskRepository) UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error {
	f.completedUserID = userID
	f.completedMessageID = messageID
	return nil
}

func (f *fakeAwardTaskRepository) UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error {
	f.failUserID = userID
	f.failMessageID = messageID
	return nil
}

type fakeAwardPublisher struct {
	topic   string
	message string
}

func (f *fakeAwardPublisher) Publish(ctx context.Context, topic string, message string) error {
	f.topic = topic
	f.message = message
	return nil
}
