package repository

import (
	"context"
	"errors"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/infrastructure/persistent/sharding"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

type ActivityRepository struct {
	db      *gorm.DB
	sharder sharding.Router
}

var _ activity.Repository = (*ActivityRepository)(nil)
var _ activity.AccountRepository = (*ActivityRepository)(nil)
var _ activity.CreditAccountRepository = (*ActivityRepository)(nil)
var _ activity.SkuProductRepository = (*ActivityRepository)(nil)
var _ activity.SkuStockRepository = (*ActivityRepository)(nil)
var _ activity.PartakeRepository = (*ActivityRepository)(nil)
var _ activity.RebateRepository = (*ActivityRepository)(nil)
var _ activity.DeliveryRepository = (*ActivityRepository)(nil)

func NewActivityRepository(db *gorm.DB, routers ...sharding.Router) *ActivityRepository {
	router := sharding.NewRouter(1)
	if len(routers) > 0 {
		router = routers[0]
	}
	return &ActivityRepository{db: db, sharder: router}
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

func (r *ActivityRepository) QueryUserCreditAccount(ctx context.Context, userID string) (activity.CreditAccountEntity, bool, error) {
	var accountPO po.UserCreditAccount
	err := r.db.WithContext(ctx).
		Select("user_id", "available_amount").
		Where("user_id = ?", userID).
		First(&accountPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.CreditAccountEntity{}, false, nil
	}
	if err != nil {
		return activity.CreditAccountEntity{}, false, err
	}
	return activity.CreditAccountEntity{
		UserID:          accountPO.UserID,
		AvailableAmount: accountPO.AvailableAmount,
	}, true, nil
}

func (r *ActivityRepository) QueryNoUsedRaffleOrder(ctx context.Context, userID string, activityID int64) (activity.UserRaffleOrderEntity, bool, error) {
	var orderPO po.UserRaffleOrder
	err := r.db.WithContext(ctx).
		Table(r.sharder.Table("user_raffle_order", userID)).
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
		if err := tx.Table(r.sharder.Table("user_raffle_order", aggregate.UserID)).Create(&orderPO).Error; err != nil {
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

func (r *ActivityRepository) QuerySkuProductBySKU(ctx context.Context, sku int64) (activity.SkuProductEntity, bool, error) {
	var skuPO po.RaffleActivitySku
	err := r.db.WithContext(ctx).
		Select("sku", "activity_id", "activity_count_id", "stock_count", "stock_count_surplus", "product_amount").
		Where("sku = ?", sku).
		First(&skuPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.SkuProductEntity{}, false, nil
	}
	if err != nil {
		return activity.SkuProductEntity{}, false, err
	}
	activityCount, err := r.queryActivityCount(ctx, skuPO.ActivityCountID)
	if err != nil {
		return activity.SkuProductEntity{}, false, err
	}
	return activity.SkuProductEntity{
		SKU:               skuPO.SKU,
		ActivityID:        skuPO.ActivityID,
		ActivityCountID:   skuPO.ActivityCountID,
		StockCount:        skuPO.StockCount,
		StockCountSurplus: skuPO.StockCountSurplus,
		ProductAmount:     skuPO.ProductAmount,
		ActivityCount:     activityCount,
	}, true, nil
}

func (r *ActivityRepository) QueryUnpaidActivityOrder(ctx context.Context, userID string, sku int64) (activity.SkuExchangeOrderEntity, bool, error) {
	var orderPO po.RaffleActivityOrder
	err := r.db.WithContext(ctx).
		Table(r.sharder.Table("raffle_activity_order", userID)).
		Select("user_id", "sku", "order_id", "out_business_no", "pay_amount").
		Where("user_id = ? and sku = ? and state = ? and order_time >= date_sub(now(), interval 1 month)", userID, sku, activity.ActivityOrderWaitPay).
		First(&orderPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return activity.SkuExchangeOrderEntity{}, false, nil
	}
	if err != nil {
		return activity.SkuExchangeOrderEntity{}, false, err
	}
	return activity.SkuExchangeOrderEntity{
		UserID:        orderPO.UserID,
		SKU:           orderPO.SKU,
		OrderID:       orderPO.OrderID,
		OutBusinessNo: orderPO.OutBusinessNo,
		PayAmount:     orderPO.PayAmount,
	}, true, nil
}

func (r *ActivityRepository) SaveCreditPayOrder(ctx context.Context, aggregate activity.CreateSkuExchangeOrderAggregate) error {
	order := aggregate.ActivityOrder
	orderPO := po.RaffleActivityOrder{
		UserID:        order.UserID,
		SKU:           order.SKU,
		ActivityID:    order.ActivityID,
		ActivityName:  order.ActivityName,
		StrategyID:    order.StrategyID,
		OrderID:       order.OrderID,
		OrderTime:     order.OrderTime,
		TotalCount:    order.TotalCount,
		DayCount:      order.DayCount,
		MonthCount:    order.MonthCount,
		PayAmount:     order.PayAmount,
		State:         order.State,
		OutBusinessNo: order.OutBusinessNo,
	}
	if err := r.db.WithContext(ctx).Table(r.sharder.Table("raffle_activity_order", aggregate.UserID)).Create(&orderPO).Error; err != nil {
		return types.NewAppError(types.ResponseCodeIndexDup, err)
	}
	return nil
}

func (r *ActivityRepository) CompleteCreditPayOrder(ctx context.Context, aggregate activity.CompleteSkuExchangeAggregate) error {
	now := time.Now()
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := adjustUserCreditAccount(tx, aggregate.CreditOrder); err != nil {
			return err
		}
		if err := tx.Table(r.sharder.Table("user_credit_order", aggregate.UserID)).Create(&po.UserCreditOrder{
			UserID:        aggregate.CreditOrder.UserID,
			OrderID:       aggregate.CreditOrder.OrderID,
			TradeName:     aggregate.CreditOrder.TradeName,
			TradeType:     aggregate.CreditOrder.TradeType,
			TradeAmount:   aggregate.CreditOrder.TradeAmount,
			OutBusinessNo: aggregate.CreditOrder.OutBusinessNo,
			CreateTime:    now,
			UpdateTime:    now,
		}).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
		if aggregate.SendTask.MessageID != "" {
			if err := tx.Create(&po.Task{
				UserID:     aggregate.SendTask.UserID,
				Topic:      aggregate.SendTask.Topic,
				MessageID:  aggregate.SendTask.MessageID,
				Message:    aggregate.SendTask.Message,
				State:      aggregate.SendTask.State,
				CreateTime: now,
				UpdateTime: now,
			}).Error; err != nil {
				return types.NewAppError(types.ResponseCodeIndexDup, err)
			}
		}

		return nil
	})
}

func (r *ActivityRepository) SaveRebateSkuOrder(ctx context.Context, aggregate activity.CreateRebateSkuOrderAggregate) error {
	now := time.Now()
	order := aggregate.ActivityOrder
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderPO := po.RaffleActivityOrder{
			UserID:        order.UserID,
			SKU:           order.SKU,
			ActivityID:    order.ActivityID,
			ActivityName:  order.ActivityName,
			StrategyID:    order.StrategyID,
			OrderID:       order.OrderID,
			OrderTime:     order.OrderTime,
			TotalCount:    order.TotalCount,
			DayCount:      order.DayCount,
			MonthCount:    order.MonthCount,
			PayAmount:     order.PayAmount,
			State:         order.State,
			OutBusinessNo: order.OutBusinessNo,
			CreateTime:    now,
			UpdateTime:    now,
		}
		if err := tx.Table(r.sharder.Table("raffle_activity_order", aggregate.UserID)).Create(&orderPO).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
		return addActivityAccountQuota(tx, activity.CompleteSkuExchangeAggregate{
			UserID:        aggregate.UserID,
			ActivityID:    aggregate.ActivityID,
			TotalCount:    order.TotalCount,
			DayCount:      order.DayCount,
			MonthCount:    order.MonthCount,
			OutBusinessNo: order.OutBusinessNo,
		}, now)
	})
}

func (r *ActivityRepository) SaveRebateIntegralOrder(ctx context.Context, rebateIntegral activity.RebateIntegralEntity) error {
	now := time.Now()
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		creditOrder := activity.CreditOrderEntity{
			UserID:        rebateIntegral.UserID,
			OrderID:       rebateIntegral.OrderID,
			TradeName:     "REBATE",
			TradeType:     "forward",
			TradeAmount:   rebateIntegral.TradeAmount,
			OutBusinessNo: rebateIntegral.OutBusinessNo,
		}
		if err := adjustOrCreateUserCreditAccount(tx, creditOrder, now); err != nil {
			return err
		}
		if err := tx.Table(r.sharder.Table("user_credit_order", rebateIntegral.UserID)).Create(&po.UserCreditOrder{
			UserID:        creditOrder.UserID,
			OrderID:       creditOrder.OrderID,
			TradeName:     creditOrder.TradeName,
			TradeType:     creditOrder.TradeType,
			TradeAmount:   creditOrder.TradeAmount,
			OutBusinessNo: creditOrder.OutBusinessNo,
			CreateTime:    now,
			UpdateTime:    now,
		}).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
		return nil
	})
}

func (r *ActivityRepository) DeliverActivityOrder(ctx context.Context, deliveryOrder activity.DeliveryOrderEntity) error {
	now := time.Now()
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var orderPO po.RaffleActivityOrder
		err := tx.
			Table(r.sharder.Table("raffle_activity_order", deliveryOrder.UserID)).
			Select("user_id", "activity_id", "total_count", "day_count", "month_count", "state").
			Where("user_id = ? and out_business_no = ?", deliveryOrder.UserID, deliveryOrder.OutBusinessNo).
			First(&orderPO).
			Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.NewAppError(types.ResponseCodeIllegalParam, nil)
		}
		if err != nil {
			return err
		}

		result := tx.Table(r.sharder.Table("raffle_activity_order", deliveryOrder.UserID)).
			Where("user_id = ? and out_business_no = ? and state = ?", deliveryOrder.UserID, deliveryOrder.OutBusinessNo, activity.ActivityOrderWaitPay).
			Updates(map[string]any{
				"state":       activity.ActivityOrderCompleted,
				"update_time": now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return types.NewAppError(types.ResponseCodeActivityOrderStateError, nil)
		}

		return addActivityAccountQuota(tx, activity.CompleteSkuExchangeAggregate{
			UserID:        deliveryOrder.UserID,
			ActivityID:    orderPO.ActivityID,
			TotalCount:    orderPO.TotalCount,
			DayCount:      orderPO.DayCount,
			MonthCount:    orderPO.MonthCount,
			OutBusinessNo: deliveryOrder.OutBusinessNo,
		}, now)
	})
}

func (r *ActivityRepository) UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error {
	return r.updateTaskState(ctx, userID, messageID, "completed")
}

func (r *ActivityRepository) UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error {
	return r.updateTaskState(ctx, userID, messageID, "fail")
}

func (r *ActivityRepository) updateTaskState(ctx context.Context, userID string, messageID string, state string) error {
	return r.db.WithContext(ctx).
		Model(&po.Task{}).
		Where("user_id = ? and message_id = ?", userID, messageID).
		Updates(map[string]any{
			"state":       state,
			"update_time": time.Now(),
		}).
		Error
}

func (r *ActivityRepository) UpdateActivitySkuStock(ctx context.Context, sku int64) error {
	return r.db.WithContext(ctx).
		Model(&po.RaffleActivitySku{}).
		Where("sku = ? and stock_count_surplus > 0", sku).
		Updates(map[string]any{
			"stock_count_surplus": gorm.Expr("stock_count_surplus - ?", 1),
			"update_time":         time.Now(),
		}).
		Error
}

func (r *ActivityRepository) ClearActivitySkuStock(ctx context.Context, sku int64) error {
	return r.db.WithContext(ctx).
		Model(&po.RaffleActivitySku{}).
		Where("sku = ?", sku).
		Updates(map[string]any{
			"stock_count_surplus": 0,
			"update_time":         time.Now(),
		}).
		Error
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

func adjustUserCreditAccount(tx *gorm.DB, creditOrder activity.CreditOrderEntity) error {
	now := time.Now()
	result := tx.Model(&po.UserCreditAccount{}).
		Where("user_id = ? and available_amount + ? >= 0", creditOrder.UserID, creditOrder.TradeAmount).
		Updates(map[string]any{
			"total_amount":     gorm.Expr("total_amount + ?", creditOrder.TradeAmount),
			"available_amount": gorm.Expr("available_amount + ?", creditOrder.TradeAmount),
			"update_time":      now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return types.NewAppError(types.ResponseCodeAccountQuotaError, nil)
	}
	return nil
}

func adjustOrCreateUserCreditAccount(tx *gorm.DB, creditOrder activity.CreditOrderEntity, now time.Time) error {
	result := tx.Model(&po.UserCreditAccount{}).
		Where("user_id = ? and available_amount + ? >= 0", creditOrder.UserID, creditOrder.TradeAmount).
		Updates(map[string]any{
			"total_amount":     gorm.Expr("total_amount + ?", creditOrder.TradeAmount),
			"available_amount": gorm.Expr("available_amount + ?", creditOrder.TradeAmount),
			"update_time":      now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		if creditOrder.TradeAmount < 0 {
			return types.NewAppError(types.ResponseCodeAccountQuotaError, nil)
		}
		if err := tx.Create(&po.UserCreditAccount{
			UserID:          creditOrder.UserID,
			TotalAmount:     creditOrder.TradeAmount,
			AvailableAmount: creditOrder.TradeAmount,
			AccountStatus:   "open",
			CreateTime:      now,
			UpdateTime:      now,
		}).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
	}
	return nil
}

func addActivityAccountQuota(tx *gorm.DB, aggregate activity.CompleteSkuExchangeAggregate, now time.Time) error {
	account := po.RaffleActivityAccount{
		UserID:            aggregate.UserID,
		ActivityID:        aggregate.ActivityID,
		TotalCount:        aggregate.TotalCount,
		TotalCountSurplus: aggregate.TotalCount,
		DayCount:          aggregate.DayCount,
		DayCountSurplus:   aggregate.DayCount,
		MonthCount:        aggregate.MonthCount,
		MonthCountSurplus: aggregate.MonthCount,
		CreateTime:        now,
		UpdateTime:        now,
	}
	updateAccount := tx.Model(&po.RaffleActivityAccount{}).
		Where("user_id = ? and activity_id = ?", aggregate.UserID, aggregate.ActivityID).
		Updates(map[string]any{
			"total_count":         gorm.Expr("total_count + ?", aggregate.TotalCount),
			"total_count_surplus": gorm.Expr("total_count_surplus + ?", aggregate.TotalCount),
			"day_count":           aggregate.DayCount,
			"day_count_surplus":   gorm.Expr("day_count_surplus + ?", aggregate.DayCount),
			"month_count":         aggregate.MonthCount,
			"month_count_surplus": gorm.Expr("month_count_surplus + ?", aggregate.MonthCount),
			"update_time":         now,
		})
	if updateAccount.Error != nil {
		return updateAccount.Error
	}
	if updateAccount.RowsAffected == 0 {
		if err := tx.Create(&account).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
	}

	month := now.Format("2006-01")
	monthAccount := po.RaffleActivityAccountMonth{
		UserID:            aggregate.UserID,
		ActivityID:        aggregate.ActivityID,
		Month:             month,
		MonthCount:        aggregate.MonthCount,
		MonthCountSurplus: aggregate.MonthCount,
		CreateTime:        now,
		UpdateTime:        now,
	}
	updateMonth := tx.Model(&po.RaffleActivityAccountMonth{}).
		Where("user_id = ? and activity_id = ? and month = ?", aggregate.UserID, aggregate.ActivityID, month).
		Updates(map[string]any{
			"month_count":         gorm.Expr("month_count + ?", aggregate.MonthCount),
			"month_count_surplus": gorm.Expr("month_count_surplus + ?", aggregate.MonthCount),
			"update_time":         now,
		})
	if updateMonth.Error != nil {
		return updateMonth.Error
	}
	if updateMonth.RowsAffected == 0 {
		if err := tx.Create(&monthAccount).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
	}

	day := now.Format("2006-01-02")
	dayAccount := po.RaffleActivityAccountDay{
		UserID:          aggregate.UserID,
		ActivityID:      aggregate.ActivityID,
		Day:             day,
		DayCount:        aggregate.DayCount,
		DayCountSurplus: aggregate.DayCount,
		CreateTime:      now,
		UpdateTime:      now,
	}
	updateDay := tx.Model(&po.RaffleActivityAccountDay{}).
		Where("user_id = ? and activity_id = ? and day = ?", aggregate.UserID, aggregate.ActivityID, day).
		Updates(map[string]any{
			"day_count":         gorm.Expr("day_count + ?", aggregate.DayCount),
			"day_count_surplus": gorm.Expr("day_count_surplus + ?", aggregate.DayCount),
			"update_time":       now,
		})
	if updateDay.Error != nil {
		return updateDay.Error
	}
	if updateDay.RowsAffected == 0 {
		if err := tx.Create(&dayAccount).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
	}
	return nil
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
