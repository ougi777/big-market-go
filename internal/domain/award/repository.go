package award

import "context"

type Repository interface {
	SaveUserAwardRecord(ctx context.Context, record UserAwardRecordEntity) error
}
