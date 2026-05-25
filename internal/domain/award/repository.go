package award

import "context"

type Repository interface {
	SaveUserAwardRecord(ctx context.Context, record UserAwardRecordEntity) error
}

type TaskRepository interface {
	QueryNoSendMessageTaskList(ctx context.Context, limit int) ([]TaskEntity, error)
	UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error
	UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error
}

type MessagePublisher interface {
	Publish(ctx context.Context, topic string, message string) error
}
