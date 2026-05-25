package repository

import (
	"context"
	"errors"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/credit"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

func (r *ActivityRepository) DeliverActivityOrder(ctx context.Context, deliveryOrder activity.DeliveryOrderEntity) error {
	now := time.Now()
	return r.shardDB(ctx, deliveryOrder.UserID).Transaction(func(tx *gorm.DB) error {
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

		return addActivityAccountQuota(tx, credit.CompleteSkuExchangeAggregate{
			UserID:        deliveryOrder.UserID,
			ActivityID:    orderPO.ActivityID,
			TotalCount:    orderPO.TotalCount,
			DayCount:      orderPO.DayCount,
			MonthCount:    orderPO.MonthCount,
			OutBusinessNo: deliveryOrder.OutBusinessNo,
		}, now)
	})
}
