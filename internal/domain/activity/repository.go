package activity

import "context"

type AccountRepository interface {
	QueryActivityAccount(ctx context.Context, activityID int64, userID string) (AccountEntity, bool, error)
	QueryActivityAccountDay(ctx context.Context, activityID int64, userID string, day string) (AccountDayEntity, bool, error)
	QueryActivityAccountMonth(ctx context.Context, activityID int64, userID string, month string) (AccountMonthEntity, bool, error)
}

type SkuProductRepository interface {
	QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]SkuProductEntity, error)
}

type Repository interface {
	AccountRepository
	SkuProductRepository
}
