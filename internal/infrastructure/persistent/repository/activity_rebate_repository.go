package repository

import (
	"context"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/credit"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

func (r *ActivityRepository) SaveRebateSkuOrder(ctx context.Context, aggregate activity.CreateRebateSkuOrderAggregate) error {
	now := time.Now()
	order := aggregate.ActivityOrder
	return r.shardDB(ctx, aggregate.UserID).Transaction(func(tx *gorm.DB) error {
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
		return addActivityAccountQuota(tx, credit.CompleteSkuExchangeAggregate{
			UserID:        aggregate.UserID,
			ActivityID:    aggregate.ActivityID,
			TotalCount:    order.TotalCount,
			DayCount:      order.DayCount,
			MonthCount:    order.MonthCount,
			OutBusinessNo: order.OutBusinessNo,
		}, now)
	})
}
