package repository

import (
	"context"
	"errors"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

var _ activity.Repository = (*ActivityRepository)(nil)
var _ activity.AccountRepository = (*ActivityRepository)(nil)
var _ activity.SkuProductRepository = (*ActivityRepository)(nil)
var _ activity.PartakeRepository = (*ActivityRepository)(nil)

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) QueryActivityByActivityID(ctx context.Context, activityID int64) (activity.ActivityEntity, bool, error) {
	var activityPO po.RaffleActivity
	err := r.db.WithContext(ctx).
		Select("activity_id", "activity_name", "activity_desc", "begin_date_time", "end_date_time", "strategy_id", "state").
		Where("activity_id = ?", activityID).
		First(&activityPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.ActivityEntity{}, false, nil
	}
	if err != nil {
		return activity.ActivityEntity{}, false, err
	}

	return activity.ActivityEntity{
		ActivityID:    activityPO.ActivityID,
		ActivityName:  activityPO.ActivityName,
		ActivityDesc:  activityPO.ActivityDesc,
		BeginDateTime: activityPO.BeginDateTime,
		EndDateTime:   activityPO.EndDateTime,
		StrategyID:    activityPO.StrategyID,
		State:         activityPO.State,
	}, true, nil
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

func (r *ActivityRepository) QueryNoUsedRaffleOrder(ctx context.Context, userID string, activityID int64) (activity.UserRaffleOrderEntity, bool, error) {
	var orderPO po.UserRaffleOrder
	err := r.db.WithContext(ctx).
		Select("user_id", "activity_id", "activity_name", "strategy_id", "order_id", "order_time", "order_state").
		Where("user_id = ? and activity_id = ? and order_state = ?", userID, activityID, activity.UserRaffleOrderCreate).
		First(&orderPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.UserRaffleOrderEntity{}, false, nil
	}
	if err != nil {
		return activity.UserRaffleOrderEntity{}, false, err
	}

	return activity.UserRaffleOrderEntity{
		UserID:       orderPO.UserID,
		ActivityID:   orderPO.ActivityID,
		ActivityName: orderPO.ActivityName,
		StrategyID:   orderPO.StrategyID,
		OrderID:      orderPO.OrderID,
		OrderTime:    orderPO.OrderTime,
		OrderState:   orderPO.OrderState,
	}, true, nil
}

func (r *ActivityRepository) SaveCreatePartakeOrder(ctx context.Context, aggregate activity.CreatePartakeOrderAggregate) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := subtractAccountQuota(tx, aggregate.UserID, aggregate.ActivityID); err != nil {
			return err
		}
		if err := saveOrSubtractMonthAccount(tx, aggregate); err != nil {
			return err
		}
		if err := saveOrSubtractDayAccount(tx, aggregate); err != nil {
			return err
		}

		orderPO := po.UserRaffleOrder{
			UserID:       aggregate.UserRaffleOrder.UserID,
			ActivityID:   aggregate.UserRaffleOrder.ActivityID,
			ActivityName: aggregate.UserRaffleOrder.ActivityName,
			StrategyID:   aggregate.UserRaffleOrder.StrategyID,
			OrderID:      aggregate.UserRaffleOrder.OrderID,
			OrderTime:    aggregate.UserRaffleOrder.OrderTime,
			OrderState:   aggregate.UserRaffleOrder.OrderState,
		}
		if err := tx.Create(&orderPO).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
		return nil
	})
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

func subtractAccountQuota(tx *gorm.DB, userID string, activityID int64) error {
	result := tx.Model(&po.RaffleActivityAccount{}).
		Where("user_id = ? and activity_id = ? and total_count_surplus > 0", userID, activityID).
		UpdateColumn("total_count_surplus", gorm.Expr("total_count_surplus - ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return types.NewAppError(types.ResponseCodeAccountQuotaError, nil)
	}
	return nil
}

func saveOrSubtractMonthAccount(tx *gorm.DB, aggregate activity.CreatePartakeOrderAggregate) error {
	monthAccount := aggregate.ActivityAccountMonth
	if aggregate.ExistAccountMonth {
		result := tx.Model(&po.RaffleActivityAccountMonth{}).
			Where("user_id = ? and activity_id = ? and month = ? and month_count_surplus > 0", aggregate.UserID, aggregate.ActivityID, monthAccount.Month).
			UpdateColumn("month_count_surplus", gorm.Expr("month_count_surplus - ?", 1))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return types.NewAppError(types.ResponseCodeAccountMonthQuotaError, nil)
		}
		return subtractAccountMonthMirror(tx, aggregate.UserID, aggregate.ActivityID)
	}

	monthPO := po.RaffleActivityAccountMonth{
		UserID:            monthAccount.UserID,
		ActivityID:        monthAccount.ActivityID,
		Month:             monthAccount.Month,
		MonthCount:        monthAccount.MonthCount,
		MonthCountSurplus: monthAccount.MonthCountSurplus - 1,
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
	}
	if err := tx.Create(&monthPO).Error; err != nil {
		return types.NewAppError(types.ResponseCodeIndexDup, err)
	}
	return setAccountMonthMirror(tx, aggregate.UserID, aggregate.ActivityID, monthPO.MonthCountSurplus)
}

func saveOrSubtractDayAccount(tx *gorm.DB, aggregate activity.CreatePartakeOrderAggregate) error {
	dayAccount := aggregate.ActivityAccountDay
	if aggregate.ExistAccountDay {
		result := tx.Model(&po.RaffleActivityAccountDay{}).
			Where("user_id = ? and activity_id = ? and day = ? and day_count_surplus > 0", aggregate.UserID, aggregate.ActivityID, dayAccount.Day).
			UpdateColumn("day_count_surplus", gorm.Expr("day_count_surplus - ?", 1))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return types.NewAppError(types.ResponseCodeAccountDayQuotaError, nil)
		}
		return subtractAccountDayMirror(tx, aggregate.UserID, aggregate.ActivityID)
	}

	dayPO := po.RaffleActivityAccountDay{
		UserID:          dayAccount.UserID,
		ActivityID:      dayAccount.ActivityID,
		Day:             dayAccount.Day,
		DayCount:        dayAccount.DayCount,
		DayCountSurplus: dayAccount.DayCountSurplus - 1,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	if err := tx.Create(&dayPO).Error; err != nil {
		return types.NewAppError(types.ResponseCodeIndexDup, err)
	}
	return setAccountDayMirror(tx, aggregate.UserID, aggregate.ActivityID, dayPO.DayCountSurplus)
}

func subtractAccountMonthMirror(tx *gorm.DB, userID string, activityID int64) error {
	return tx.Model(&po.RaffleActivityAccount{}).
		Where("user_id = ? and activity_id = ? and month_count_surplus > 0", userID, activityID).
		UpdateColumn("month_count_surplus", gorm.Expr("month_count_surplus - ?", 1)).
		Error
}

func subtractAccountDayMirror(tx *gorm.DB, userID string, activityID int64) error {
	return tx.Model(&po.RaffleActivityAccount{}).
		Where("user_id = ? and activity_id = ? and day_count_surplus > 0", userID, activityID).
		UpdateColumn("day_count_surplus", gorm.Expr("day_count_surplus - ?", 1)).
		Error
}

func setAccountMonthMirror(tx *gorm.DB, userID string, activityID int64, surplus int) error {
	return tx.Model(&po.RaffleActivityAccount{}).
		Where("user_id = ? and activity_id = ?", userID, activityID).
		UpdateColumn("month_count_surplus", surplus).
		Error
}

func setAccountDayMirror(tx *gorm.DB, userID string, activityID int64, surplus int) error {
	return tx.Model(&po.RaffleActivityAccount{}).
		Where("user_id = ? and activity_id = ?", userID, activityID).
		UpdateColumn("day_count_surplus", surplus).
		Error
}
