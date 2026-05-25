package repository

import (
	"context"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/award"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

type AwardRepository struct {
	db *gorm.DB
}

var _ award.Repository = (*AwardRepository)(nil)
var _ award.TaskRepository = (*AwardRepository)(nil)

func NewAwardRepository(db *gorm.DB) *AwardRepository {
	return &AwardRepository{db: db}
}

func (r *AwardRepository) SaveUserAwardRecord(ctx context.Context, record award.UserAwardRecordEntity) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
		if err := tx.Create(&recordPO).Error; err != nil {
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

		result := tx.Model(&po.UserRaffleOrder{}).
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

func (r *AwardRepository) QueryNoSendMessageTaskList(ctx context.Context, limit int) ([]award.TaskEntity, error) {
	if limit <= 0 {
		limit = 10
	}

	var taskPOList []po.Task
	err := r.db.WithContext(ctx).
		Select("user_id", "topic", "message_id", "message", "state").
		Where("state = ? or (state = ? and update_time < date_sub(now(), interval 6 second))", award.TaskStateFail, award.TaskStateCreate).
		Limit(limit).
		Find(&taskPOList).
		Error
	if err != nil {
		return nil, err
	}

	tasks := make([]award.TaskEntity, 0, len(taskPOList))
	for _, taskPO := range taskPOList {
		tasks = append(tasks, award.TaskEntity{
			UserID:    taskPO.UserID,
			Topic:     taskPO.Topic,
			MessageID: taskPO.MessageID,
			Message:   taskPO.Message,
			State:     taskPO.State,
		})
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
	return r.db.WithContext(ctx).
		Model(&po.Task{}).
		Where("user_id = ? and message_id = ?", userID, messageID).
		Updates(map[string]any{
			"state":       state,
			"update_time": time.Now(),
		}).
		Error
}
