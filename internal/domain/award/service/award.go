package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"bm-go/internal/domain/award"
)

type AwardService struct {
	repo             award.Repository
	taskRepo         award.TaskRepository
	publisher        award.MessagePublisher
	now              func() time.Time
	messageGenerator func() (string, error)
}

func NewAwardService(repo award.Repository, taskRepo award.TaskRepository, publisher award.MessagePublisher) *AwardService {
	return &AwardService{
		repo:             repo,
		taskRepo:         taskRepo,
		publisher:        publisher,
		now:              time.Now,
		messageGenerator: func() (string, error) { return randomNumeric(11) },
	}
}

func (s *AwardService) SaveUserAwardRecord(ctx context.Context, record award.UserAwardRecordEntity) error {
	messageID, err := s.messageGenerator()
	if err != nil {
		return err
	}
	message, err := json.Marshal(award.EventMessage[award.SendAwardMessage]{
		ID:        messageID,
		Timestamp: s.now().UnixMilli(),
		Data: award.SendAwardMessage{
			UserID:      record.UserID,
			OrderID:     record.OrderID,
			AwardID:     record.AwardID,
			AwardTitle:  record.AwardTitle,
			AwardConfig: record.AwardConfig,
		},
	})
	if err != nil {
		return err
	}

	record.SendTask = award.TaskEntity{
		UserID:    record.UserID,
		Topic:     award.TopicSendAward,
		MessageID: messageID,
		Message:   string(message),
		State:     award.TaskStateCreate,
	}

	if err := s.repo.SaveUserAwardRecord(ctx, record); err != nil {
		return err
	}
	if s.publisher == nil {
		return nil
	}
	if err := s.publisher.Publish(ctx, record.SendTask.Topic, record.SendTask.Message); err != nil {
		if s.taskRepo != nil {
			_ = s.taskRepo.UpdateTaskSendMessageFail(ctx, record.SendTask.UserID, record.SendTask.MessageID)
		}
		return err
	}
	if s.taskRepo != nil {
		return s.taskRepo.UpdateTaskSendMessageCompleted(ctx, record.SendTask.UserID, record.SendTask.MessageID)
	}
	return nil
}

func randomNumeric(length int) (string, error) {
	value := make([]byte, length)
	for i := range value {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("generate random numeric: %w", err)
		}
		value[i] = byte('0' + n.Int64())
	}
	return string(value), nil
}
