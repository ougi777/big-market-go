package job

import (
	"context"

	"go.uber.org/zap"
)

type MessageTaskSender interface {
	SendNoSendMessageTasks(ctx context.Context, limit int) error
}

type SendMessageTaskJob struct {
	sender MessageTaskSender
	logger *zap.Logger
	limit  int
}

func NewSendMessageTaskJob(sender MessageTaskSender, logger *zap.Logger) *SendMessageTaskJob {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &SendMessageTaskJob{
		sender: sender,
		logger: logger,
		limit:  10,
	}
}

func (j *SendMessageTaskJob) Exec() {
	if err := j.sender.SendNoSendMessageTasks(context.Background(), j.limit); err != nil {
		j.logger.Error("send message task job failed", zap.Error(err))
	}
}
