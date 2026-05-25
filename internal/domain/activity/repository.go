package activity

import "context"

type AccountRepository interface {
	QueryActivityAccount(ctx context.Context, activityID int64, userID string) (AccountEntity, bool, error)
	QueryActivityAccountDay(ctx context.Context, activityID int64, userID string, day string) (AccountDayEntity, bool, error)
	QueryActivityAccountMonth(ctx context.Context, activityID int64, userID string, month string) (AccountMonthEntity, bool, error)
}

type SkuProductRepository interface {
	QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]SkuProductEntity, error)
	QuerySkuProductBySKU(ctx context.Context, sku int64) (SkuProductEntity, bool, error)
}

type SkuStockStore interface {
	CacheActivitySkuStockCount(ctx context.Context, key string, stockCount int) error
	SubtractActivitySkuStock(ctx context.Context, key string) (int64, error)
}

type SkuStockQueue interface {
	SendActivitySkuStockConsumeQueue(ctx context.Context, stockKey ActivitySkuStockKey) error
	TakeActivitySkuStock(ctx context.Context) (ActivitySkuStockKey, bool, error)
	ClearActivitySkuStockQueue(ctx context.Context) error
}

type SkuStockRepository interface {
	UpdateActivitySkuStock(ctx context.Context, sku int64) error
	ClearActivitySkuStock(ctx context.Context, sku int64) error
}

type SkuExchangeRepository interface {
	QueryUnpaidActivityOrder(ctx context.Context, userID string, sku int64) (SkuExchangeOrderEntity, bool, error)
	SaveCreditPayOrder(ctx context.Context, aggregate CreateSkuExchangeOrderAggregate) error
	CompleteCreditPayOrder(ctx context.Context, aggregate CompleteSkuExchangeAggregate) error
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
	SkuExchangeRepository
	PartakeRepository
}
