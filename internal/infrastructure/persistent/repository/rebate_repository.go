package repository

import (
	"context"
	"time"

	"bm-go/internal/domain/rebate"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/infrastructure/persistent/sharding"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

type RebateRepository struct {
	db      dbRouter
	sharder sharding.Router
}

var _ rebate.Repository = (*RebateRepository)(nil)

func NewRebateRepository(db *gorm.DB, routers ...sharding.Router) *RebateRepository {
	return NewRebateRepositoryWithDBRouter(singleDBRouter{db: db}, routers...)
}

func NewRebateRepositoryWithDBRouter(db dbRouter, routers ...sharding.Router) *RebateRepository {
	router := sharding.NewRouter(1)
	if len(routers) > 0 {
		router = routers[0]
	}
	return &RebateRepository{db: db, sharder: router}
}

func (r *RebateRepository) defaultDB(ctx context.Context) *gorm.DB {
	return r.db.Default().WithContext(ctx)
}

func (r *RebateRepository) shardDB(ctx context.Context, userID string) *gorm.DB {
	return r.db.Shard(r.sharder.DBKey(userID)).WithContext(ctx)
}

func (r *RebateRepository) QueryDailyBehaviorRebateConfig(ctx context.Context, behaviorType string) ([]rebate.DailyBehaviorRebateEntity, error) {
	var configPOList []po.DailyBehaviorRebate
	err := r.defaultDB(ctx).
		Select("behavior_type", "rebate_desc", "rebate_type", "rebate_config").
		Where("behavior_type = ? and state = ?", behaviorType, rebate.RebateStateOpen).
		Find(&configPOList).
		Error
	if err != nil {
		return nil, err
	}

	configs := make([]rebate.DailyBehaviorRebateEntity, 0, len(configPOList))
	for _, configPO := range configPOList {
		configs = append(configs, rebate.DailyBehaviorRebateEntity{
			BehaviorType: configPO.BehaviorType,
			RebateDesc:   configPO.RebateDesc,
			RebateType:   configPO.RebateType,
			RebateConfig: configPO.RebateConfig,
		})
	}
	return configs, nil
}

func (r *RebateRepository) SaveUserRebateRecords(ctx context.Context, aggregates []rebate.BehaviorRebateAggregate) error {
	if len(aggregates) == 0 {
		return nil
	}
	return r.shardDB(ctx, aggregates[0].UserID).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		for _, aggregate := range aggregates {
			orderPO := po.UserBehaviorRebateOrder{
				UserID:        aggregate.Order.UserID,
				OrderID:       aggregate.Order.OrderID,
				BehaviorType:  aggregate.Order.BehaviorType,
				RebateDesc:    aggregate.Order.RebateDesc,
				RebateType:    aggregate.Order.RebateType,
				RebateConfig:  aggregate.Order.RebateConfig,
				OutBusinessNo: aggregate.Order.OutBusinessNo,
				BizID:         aggregate.Order.BizID,
				CreateTime:    now,
				UpdateTime:    now,
			}
			if err := tx.Table(r.sharder.Table("user_behavior_rebate_order", aggregate.UserID)).Create(&orderPO).Error; err != nil {
				return types.NewAppError(types.ResponseCodeIndexDup, err)
			}

			taskPO := po.Task{
				UserID:     aggregate.Task.UserID,
				Topic:      aggregate.Task.Topic,
				MessageID:  aggregate.Task.MessageID,
				Message:    aggregate.Task.Message,
				State:      aggregate.Task.State,
				CreateTime: now,
				UpdateTime: now,
			}
			if err := tx.Create(&taskPO).Error; err != nil {
				return types.NewAppError(types.ResponseCodeIndexDup, err)
			}
		}
		return nil
	})
}

func (r *RebateRepository) QueryOrderByOutBusinessNo(ctx context.Context, userID string, outBusinessNo string) ([]rebate.BehaviorRebateOrderEntity, error) {
	var orderPOList []po.UserBehaviorRebateOrder
	err := r.shardDB(ctx, userID).
		Table(r.sharder.Table("user_behavior_rebate_order", userID)).
		Select("user_id", "order_id", "behavior_type", "rebate_desc", "rebate_type", "rebate_config", "out_business_no", "biz_id").
		Where("user_id = ? and out_business_no = ?", userID, outBusinessNo).
		Find(&orderPOList).
		Error
	if err != nil {
		return nil, err
	}

	orders := make([]rebate.BehaviorRebateOrderEntity, 0, len(orderPOList))
	for _, orderPO := range orderPOList {
		orders = append(orders, rebate.BehaviorRebateOrderEntity{
			UserID:        orderPO.UserID,
			OrderID:       orderPO.OrderID,
			BehaviorType:  orderPO.BehaviorType,
			RebateDesc:    orderPO.RebateDesc,
			RebateType:    orderPO.RebateType,
			RebateConfig:  orderPO.RebateConfig,
			OutBusinessNo: orderPO.OutBusinessNo,
			BizID:         orderPO.BizID,
		})
	}
	return orders, nil
}

func (r *RebateRepository) UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error {
	return r.updateTaskState(ctx, userID, messageID, rebate.TaskStateComplete)
}

func (r *RebateRepository) UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error {
	return r.updateTaskState(ctx, userID, messageID, rebate.TaskStateFail)
}

func (r *RebateRepository) updateTaskState(ctx context.Context, userID string, messageID string, state string) error {
	return setTaskState(ctx, r.db.Shard(r.sharder.DBKey(userID)), userID, messageID, state)
}
