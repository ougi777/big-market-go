package service

import (
	"context"

	"bm-go/internal/domain/task"
)

type Service struct {
	repo      task.Repository
	publisher task.MessagePublisher
}

func NewService(repo task.Repository, publisher task.MessagePublisher) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) SendNoSendMessageTasks(ctx context.Context, limit int) error {
	tasks, err := s.repo.QueryNoSendMessageTaskList(ctx, limit)
	if err != nil {
		return err
	}

	for _, taskEntity := range tasks {
		if err := s.publisher.Publish(ctx, taskEntity.Topic, taskEntity.Message); err != nil {
			if updateErr := s.repo.UpdateTaskSendMessageFail(ctx, taskEntity.UserID, taskEntity.MessageID); updateErr != nil {
				return updateErr
			}
			continue
		}
		if err := s.repo.UpdateTaskSendMessageCompleted(ctx, taskEntity.UserID, taskEntity.MessageID); err != nil {
			return err
		}
	}
	return nil
}
