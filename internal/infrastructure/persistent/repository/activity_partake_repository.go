package repository

import (
	"context"
	"errors"

	"bm-go/internal/domain/activity"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

func (r *ActivityRepository) QueryNoUsedRaffleOrder(ctx context.Context, userID string, activityID int64) (activity.UserRaffleOrderEntity, bool, error) {
	var orderPO po.UserRaffleOrder
	err := r.shardDB(ctx, userID).
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
	return r.shardDB(ctx, aggregate.UserID).Transaction(func(tx *gorm.DB) error {
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
