package repository

import (
	"context"
	"errors"

	"bm-go/internal/domain/activity"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

func (r *ActivityRepository) QueryUnpaidActivityOrder(ctx context.Context, userID string, sku int64) (activity.SkuExchangeOrderEntity, bool, error) {
	var orderPO po.RaffleActivityOrder
	err := r.shardDB(ctx, userID).
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
	if err := r.shardDB(ctx, aggregate.UserID).Table(r.sharder.Table("raffle_activity_order", aggregate.UserID)).Create(&orderPO).Error; err != nil {
		return types.NewAppError(types.ResponseCodeIndexDup, err)
	}
	return nil
}
