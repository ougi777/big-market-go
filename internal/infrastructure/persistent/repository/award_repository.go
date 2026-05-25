package repository

import (
	"context"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/award"
	taskdomain "bm-go/internal/domain/task"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/infrastructure/persistent/sharding"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

type AwardRepository struct {
	db      dbRouter
	sharder sharding.Router
}

var _ award.Repository = (*AwardRepository)(nil)
var _ award.TaskRepository = (*AwardRepository)(nil)
var _ taskdomain.Repository = (*AwardRepository)(nil)

func NewAwardRepository(db *gorm.DB, routers ...sharding.Router) *AwardRepository {
	return NewAwardRepositoryWithDBRouter(singleDBRouter{db: db}, routers...)
}

func NewAwardRepositoryWithDBRouter(db dbRouter, routers ...sharding.Router) *AwardRepository {
	router := sharding.NewRouter(1)
	if len(routers) > 0 {
		router = routers[0]
	}
	return &AwardRepository{db: db, sharder: router}
}

func (r *AwardRepository) defaultDB(ctx context.Context) *gorm.DB {
	return r.db.Default().WithContext(ctx)
}

func (r *AwardRepository) shardDB(ctx context.Context, userID string) *gorm.DB {
	return r.db.Shard(r.sharder.DBKey(userID)).WithContext(ctx)
}

func (r *AwardRepository) taskDBs(ctx context.Context) []*gorm.DB {
	connections := r.db.Connections()
	dbs := make([]*gorm.DB, 0, len(connections))
	for _, db := range connections {
		if db != nil {
			dbs = append(dbs, db.WithContext(ctx))
		}
	}
	if len(dbs) == 0 {
		return []*gorm.DB{r.defaultDB(ctx)}
	}
	return dbs
}

func (r *AwardRepository) SaveUserAwardRecord(ctx context.Context, record award.UserAwardRecordEntity) error {
	return r.shardDB(ctx, record.UserID).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		recordPO := po.UserAwardRecord{
			UserID:     record.UserID,
			ActivityID: record.ActivityID,
			StrategyID: record.StrategyID,
			OrderID:    record.OrderID,
			AwardID:    record.AwardID,
			AwardTitle: record.AwardTitle,
			AwardTime:  record.AwardTime,
			AwardState: record.AwardState,
		}
		if err := tx.Table(r.sharder.Table("user_award_record", record.UserID)).Create(&recordPO).Error; err != nil {
			return types.NewAppError(types.ResponseCodeIndexDup, err)
		}
		if record.SendTask.MessageID != "" {
			taskPO := po.Task{
				UserID:     record.SendTask.UserID,
				Topic:      record.SendTask.Topic,
				MessageID:  record.SendTask.MessageID,
				Message:    record.SendTask.Message,
				State:      record.SendTask.State,
				CreateTime: now,
				UpdateTime: now,
			}
			if err := tx.Create(&taskPO).Error; err != nil {
				return types.NewAppError(types.ResponseCodeIndexDup, err)
			}
		}

		result := tx.Table(r.sharder.Table("user_raffle_order", record.UserID)).
			Where("user_id = ? and order_id = ? and order_state = ?", record.UserID, record.OrderID, activity.UserRaffleOrderCreate).
			Update("order_state", activity.UserRaffleOrderUsed)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return types.NewAppError(types.ResponseCodeActivityOrderStateError, nil)
		}
		return nil
	})
}

func (r *AwardRepository) QueryAwardConfig(ctx context.Context, awardID int) (string, error) {
	var awardPO po.Award
	err := r.defaultDB(ctx).
		Select("award_config").
		Where("award_id = ?", awardID).
		Take(&awardPO).
		Error
	if err != nil {
		return "", err
	}
	return awardPO.AwardConfig, nil
}

func (r *AwardRepository) QueryAwardKey(ctx context.Context, awardID int) (string, error) {
	var awardPO po.Award
	err := r.defaultDB(ctx).
		Select("award_key").
		Where("award_id = ?", awardID).
		Take(&awardPO).
		Error
	if err != nil {
		return "", err
	}
	return awardPO.AwardKey, nil
}

func (r *AwardRepository) SaveGiveOutPrizes(ctx context.Context, aggregate award.GiveOutPrizesAggregate) error {
	return r.shardDB(ctx, aggregate.UserAwardRecord.UserID).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		creditPO := po.UserCreditAccount{
			UserID:          aggregate.UserCreditAward.UserID,
			TotalAmount:     aggregate.UserCreditAward.CreditAmount,
			AvailableAmount: aggregate.UserCreditAward.CreditAmount,
			AccountStatus:   award.AccountStatusOpen,
			CreateTime:      now,
			UpdateTime:      now,
		}

		updateCredit := tx.Model(&po.UserCreditAccount{}).
			Where("user_id = ?", aggregate.UserCreditAward.UserID).
			Updates(map[string]any{
				"total_amount":     gorm.Expr("total_amount + ?", aggregate.UserCreditAward.CreditAmount),
				"available_amount": gorm.Expr("available_amount + ?", aggregate.UserCreditAward.CreditAmount),
				"update_time":      now,
			})
		if updateCredit.Error != nil {
			return updateCredit.Error
		}
		if updateCredit.RowsAffected == 0 {
			if err := tx.Create(&creditPO).Error; err != nil {
				return types.NewAppError(types.ResponseCodeIndexDup, err)
			}
		}

		updateAward := tx.Table(r.sharder.Table("user_award_record", aggregate.UserAwardRecord.UserID)).
			Where("user_id = ? and order_id = ? and award_state = ?", aggregate.UserAwardRecord.UserID, aggregate.UserAwardRecord.OrderID, award.AwardStateCreate).
			Updates(map[string]any{
				"award_state": aggregate.UserAwardRecord.AwardState,
				"update_time": now,
			})
		if updateAward.Error != nil {
			return updateAward.Error
		}
		if updateAward.RowsAffected != 1 {
			return types.NewAppError(types.ResponseCodeActivityOrderStateError, nil)
		}
		return nil
	})
}

func (r *AwardRepository) QueryNoSendMessageTaskList(ctx context.Context, limit int) ([]taskdomain.Entity, error) {
	if limit <= 0 {
		limit = 10
	}

	tasks := make([]taskdomain.Entity, 0, limit)
	for _, db := range r.taskDBs(ctx) {
		remaining := limit - len(tasks)
		if remaining <= 0 {
			break
		}

		var taskPOList []po.Task
		err := db.
			Select("user_id", "topic", "message_id", "message", "state").
			Where("state = ? or (state = ? and update_time < date_sub(now(), interval 6 second))", award.TaskStateFail, award.TaskStateCreate).
			Limit(remaining).
			Find(&taskPOList).
			Error
		if err != nil {
			return nil, err
		}

		for _, taskPO := range taskPOList {
			tasks = append(tasks, taskdomain.Entity{
				UserID:    taskPO.UserID,
				Topic:     taskPO.Topic,
				MessageID: taskPO.MessageID,
				Message:   taskPO.Message,
				State:     taskPO.State,
			})
		}
	}

	return tasks, nil
}

func (r *AwardRepository) UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error {
	return r.updateTaskState(ctx, userID, messageID, award.TaskStateCompleted)
}

func (r *AwardRepository) UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error {
	return r.updateTaskState(ctx, userID, messageID, award.TaskStateFail)
}

func (r *AwardRepository) updateTaskState(ctx context.Context, userID string, messageID string, state string) error {
	return r.shardDB(ctx, userID).
		Model(&po.Task{}).
		Where("user_id = ? and message_id = ?", userID, messageID).
		Updates(map[string]any{
			"state":       state,
			"update_time": time.Now(),
		}).
		Error
}
