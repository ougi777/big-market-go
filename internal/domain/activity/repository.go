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

type SkuStockStore interface {
	CacheActivitySkuStockCount(ctx context.Context, key string, stockCount int) error
}

type PartakeRepository interface {
	QueryActivityByActivityID(ctx context.Context, activityID int64) (ActivityEntity, bool, error)
	QueryNoUsedRaffleOrder(ctx context.Context, userID string, activityID int64) (UserRaffleOrderEntity, bool, error)
	QueryActivityAccount(ctx context.Context, activityID int64, userID string) (AccountEntity, bool, error)
	QueryActivityAccountDay(ctx context.Context, activityID int64, userID string, day string) (AccountDayEntity, bool, error)
	QueryActivityAccountMonth(ctx context.Context, activityID int64, userID string, month string) (AccountMonthEntity, bool, error)
	SaveCreatePartakeOrder(ctx context.Context, aggregate CreatePartakeOrderAggregate) error
}

type Repository interface {
	AccountRepository
	SkuProductRepository
	PartakeRepository
}
