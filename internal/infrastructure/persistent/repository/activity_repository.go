package repository

import (
	"context"
	"errors"

	"bm-go/internal/domain/activity"
	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

var _ activity.Repository = (*ActivityRepository)(nil)
var _ activity.AccountRepository = (*ActivityRepository)(nil)
var _ activity.SkuProductRepository = (*ActivityRepository)(nil)

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) QueryActivityAccount(ctx context.Context, activityID int64, userID string) (activity.AccountEntity, bool, error) {
	var accountPO po.RaffleActivityAccount
	err := r.db.WithContext(ctx).
		Select("user_id", "activity_id", "total_count", "total_count_surplus", "day_count", "day_count_surplus", "month_count", "month_count_surplus").
		Where("user_id = ? and activity_id = ?", userID, activityID).
		First(&accountPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.AccountEntity{}, false, nil
	}
	if err != nil {
		return activity.AccountEntity{}, false, err
	}

	return activity.AccountEntity{
		UserID:            accountPO.UserID,
		ActivityID:        accountPO.ActivityID,
		TotalCount:        accountPO.TotalCount,
		TotalCountSurplus: accountPO.TotalCountSurplus,
		DayCount:          accountPO.DayCount,
		DayCountSurplus:   accountPO.DayCountSurplus,
		MonthCount:        accountPO.MonthCount,
		MonthCountSurplus: accountPO.MonthCountSurplus,
	}, true, nil
}

func (r *ActivityRepository) QueryActivityAccountDay(ctx context.Context, activityID int64, userID string, day string) (activity.AccountDayEntity, bool, error) {
	var dayPO po.RaffleActivityAccountDay
	err := r.db.WithContext(ctx).
		Select("user_id", "activity_id", "day", "day_count", "day_count_surplus").
		Where("user_id = ? and activity_id = ? and day = ?", userID, activityID, day).
		First(&dayPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.AccountDayEntity{}, false, nil
	}
	if err != nil {
		return activity.AccountDayEntity{}, false, err
	}

	return activity.AccountDayEntity{
		UserID:          dayPO.UserID,
		ActivityID:      dayPO.ActivityID,
		Day:             dayPO.Day,
		DayCount:        dayPO.DayCount,
		DayCountSurplus: dayPO.DayCountSurplus,
	}, true, nil
}

func (r *ActivityRepository) QueryActivityAccountMonth(ctx context.Context, activityID int64, userID string, month string) (activity.AccountMonthEntity, bool, error) {
	var monthPO po.RaffleActivityAccountMonth
	err := r.db.WithContext(ctx).
		Select("user_id", "activity_id", "month", "month_count", "month_count_surplus").
		Where("user_id = ? and activity_id = ? and month = ?", userID, activityID, month).
		First(&monthPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.AccountMonthEntity{}, false, nil
	}
	if err != nil {
		return activity.AccountMonthEntity{}, false, err
	}

	return activity.AccountMonthEntity{
		UserID:            monthPO.UserID,
		ActivityID:        monthPO.ActivityID,
		Month:             monthPO.Month,
		MonthCount:        monthPO.MonthCount,
		MonthCountSurplus: monthPO.MonthCountSurplus,
	}, true, nil
}

func (r *ActivityRepository) QuerySkuProductListByActivityID(ctx context.Context, activityID int64) ([]activity.SkuProductEntity, error) {
	var skuPOList []po.RaffleActivitySku
	err := r.db.WithContext(ctx).
		Select("sku", "activity_id", "activity_count_id", "stock_count", "stock_count_surplus", "product_amount").
		Where("activity_id = ?", activityID).
		Find(&skuPOList).
		Error
	if err != nil {
		return nil, err
	}

	products := make([]activity.SkuProductEntity, 0, len(skuPOList))
	for _, skuPO := range skuPOList {
		activityCount, err := r.queryActivityCount(ctx, skuPO.ActivityCountID)
		if err != nil {
			return nil, err
		}

		products = append(products, activity.SkuProductEntity{
			SKU:               skuPO.SKU,
			ActivityID:        skuPO.ActivityID,
			ActivityCountID:   skuPO.ActivityCountID,
			StockCount:        skuPO.StockCount,
			StockCountSurplus: skuPO.StockCountSurplus,
			ProductAmount:     skuPO.ProductAmount,
			ActivityCount:     activityCount,
		})
	}
	return products, nil
}

func (r *ActivityRepository) queryActivityCount(ctx context.Context, activityCountID int64) (activity.ActivityCountEntity, error) {
	var countPO po.RaffleActivityCount
	err := r.db.WithContext(ctx).
		Select("activity_count_id", "total_count", "day_count", "month_count").
		Where("activity_count_id = ?", activityCountID).
		First(&countPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.ActivityCountEntity{ActivityCountID: activityCountID}, nil
	}
	if err != nil {
		return activity.ActivityCountEntity{}, err
	}

	return activity.ActivityCountEntity{
		ActivityCountID: countPO.ActivityCountID,
		TotalCount:      countPO.TotalCount,
		DayCount:        countPO.DayCount,
		MonthCount:      countPO.MonthCount,
	}, nil
}
