package task

import "context"

type Repository interface {
	QueryNoSendMessageTaskList(ctx context.Context, limit int) ([]Entity, error)
	UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error
	UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error
}

type MessagePublisher interface {
	Publish(ctx context.Context, topic string, message string) error
}
