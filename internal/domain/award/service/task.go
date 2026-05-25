package service

import (
	"context"

	"bm-go/internal/domain/award"
)

type TaskService struct {
	repo      award.TaskRepository
	publisher award.MessagePublisher
}

func NewTaskService(repo award.TaskRepository, publisher award.MessagePublisher) *TaskService {
	return &TaskService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *TaskService) SendNoSendMessageTasks(ctx context.Context, limit int) error {
	tasks, err := s.repo.QueryNoSendMessageTaskList(ctx, limit)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := s.publisher.Publish(ctx, task.Topic, task.Message); err != nil {
			if updateErr := s.repo.UpdateTaskSendMessageFail(ctx, task.UserID, task.MessageID); updateErr != nil {
				return updateErr
			}
			continue
		}
		if err := s.repo.UpdateTaskSendMessageCompleted(ctx, task.UserID, task.MessageID); err != nil {
			return err
		}
	}
	return nil
}
