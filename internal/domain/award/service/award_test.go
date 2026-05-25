package service

import (
	"context"
	"errors"
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

func TestAwardServiceSaveUserAwardRecordMarksTaskFailOnPublishError(t *testing.T) {
	repo := &fakeAwardRepository{}
	taskRepo := &fakeAwardTaskRepository{}
	publisher := &fakeAwardPublisher{err: errors.New("publish failed")}
	service := NewAwardService(repo, taskRepo, publisher)
	service.messageGenerator = func() (string, error) { return "12345678901", nil }

	err := service.SaveUserAwardRecord(context.Background(), award.UserAwardRecordEntity{
		UserID:     "xiaofuge",
		ActivityID: 100301,
		StrategyID: 100006,
		OrderID:    "order-001",
		AwardID:    101,
		AwardTitle: "随机积分",
		AwardState: award.AwardStateCreate,
	})
	if err == nil {
		t.Fatal("expected publish error")
	}
	if taskRepo.failUserID != "xiaofuge" || taskRepo.failMessageID != "12345678901" {
		t.Fatalf("expected task fail, got %s/%s", taskRepo.failUserID, taskRepo.failMessageID)
	}
}

func TestAwardServiceDistributeAward(t *testing.T) {
	repo := &fakeAwardRepository{
		awardKey: award.AwardKeyUserCreditRand,
	}
	service := NewAwardService(repo, nil, nil)
	service.creditGenerator = func(min float64, max float64) (float64, error) {
		if min != 0.01 || max != 1 {
			t.Fatalf("expected credit range 0.01/1, got %.2f/%.2f", min, max)
		}
		return 0.58, nil
	}

	err := service.DistributeAward(context.Background(), award.DistributeAwardEntity{
		UserID:      "xiaofuge",
		OrderID:     "order-001",
		AwardID:     101,
		AwardConfig: "0.01,1",
	})
	if err != nil {
		t.Fatalf("distribute award: %v", err)
	}

	if repo.aggregate.UserAwardRecord.AwardState != award.AwardStateComplete {
		t.Fatalf("expected award completed, got %s", repo.aggregate.UserAwardRecord.AwardState)
	}
	if repo.aggregate.UserCreditAward.CreditAmount != 0.58 {
		t.Fatalf("expected credit amount 0.58, got %.2f", repo.aggregate.UserCreditAward.CreditAmount)
	}
}

type fakeAwardRepository struct {
	record      award.UserAwardRecordEntity
	awardKey    string
	awardConfig string
	aggregate   award.GiveOutPrizesAggregate
}

func (f *fakeAwardRepository) SaveUserAwardRecord(ctx context.Context, record award.UserAwardRecordEntity) error {
	f.record = record
	return nil
}

func (f *fakeAwardRepository) QueryAwardConfig(ctx context.Context, awardID int) (string, error) {
	return f.awardConfig, nil
}

func (f *fakeAwardRepository) QueryAwardKey(ctx context.Context, awardID int) (string, error) {
	return f.awardKey, nil
}

func (f *fakeAwardRepository) SaveGiveOutPrizes(ctx context.Context, aggregate award.GiveOutPrizesAggregate) error {
	f.aggregate = aggregate
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
	err     error
}

func (f *fakeAwardPublisher) Publish(ctx context.Context, topic string, message string) error {
	f.topic = topic
	f.message = message
	if f.err != nil {
		return f.err
	}
	return nil
}
