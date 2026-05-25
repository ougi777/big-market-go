package repository

import (
	"context"

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

func NewAwardRepository(db *gorm.DB) *AwardRepository {
	return &AwardRepository{db: db}
}

func (r *AwardRepository) SaveUserAwardRecord(ctx context.Context, record award.UserAwardRecordEntity) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
