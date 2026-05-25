package award

import "context"

import taskdomain "bm-go/internal/domain/task"

type Repository interface {
	SaveUserAwardRecord(ctx context.Context, record UserAwardRecordEntity) error
	QueryAwardConfig(ctx context.Context, awardID int) (string, error)
	QueryAwardKey(ctx context.Context, awardID int) (string, error)
	SaveGiveOutPrizes(ctx context.Context, aggregate GiveOutPrizesAggregate) error
}

type TaskRepository = taskdomain.Repository

type MessagePublisher = taskdomain.MessagePublisher
