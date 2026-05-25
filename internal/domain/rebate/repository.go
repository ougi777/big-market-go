package rebate

import "context"

type Repository interface {
	QueryDailyBehaviorRebateConfig(ctx context.Context, behaviorType string) ([]DailyBehaviorRebateEntity, error)
	SaveUserRebateRecords(ctx context.Context, aggregates []BehaviorRebateAggregate) error
	QueryOrderByOutBusinessNo(ctx context.Context, userID string, outBusinessNo string) ([]BehaviorRebateOrderEntity, error)
	UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error
	UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error
}

type MessagePublisher interface {
	Publish(ctx context.Context, topic string, message string) error
}
