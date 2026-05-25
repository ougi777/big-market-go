package repository

import (
	"context"
	"errors"
	"time"

	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

func (r *StrategyRepository) QueryActivityAccountTotalUseCount(ctx context.Context, userID string, strategyID int64) (int, error) {
	var activityPO po.RaffleActivity
	err := r.defaultDB(ctx).
		Select("activity_id").
		Where("strategy_id = ?", strategyID).
		First(&activityPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	var accountPO po.RaffleActivityAccount
	err = r.shardDB(ctx, userID).
		Select("total_count", "total_count_surplus").
		Where("user_id = ? and activity_id = ?", userID, activityPO.ActivityID).
		First(&accountPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return accountPO.TotalCount - accountPO.TotalCountSurplus, nil
}

func (r *StrategyRepository) QueryRaffleActivityAccountDayPartakeCount(ctx context.Context, activityID int64, userID string) (int, error) {
	var accountDayPO po.RaffleActivityAccountDay
	err := r.shardDB(ctx, userID).
		Select("day_count", "day_count_surplus").
		Where("user_id = ? and activity_id = ? and day = ?", userID, activityID, time.Now().Format("2006-01-02")).
		First(&accountDayPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return accountDayPO.DayCount - accountDayPO.DayCountSurplus, nil
}

func (r *StrategyRepository) QueryTodayUserRaffleCount(ctx context.Context, userID string, strategyID int64) (int, error) {
	activityID, err := r.queryActivityIDByStrategyID(ctx, strategyID)
	if err != nil {
		return 0, err
	}
	if activityID == 0 {
		return 0, nil
	}
	return r.QueryRaffleActivityAccountDayPartakeCount(ctx, activityID, userID)
}

func (r *StrategyRepository) queryActivityIDByStrategyID(ctx context.Context, strategyID int64) (int64, error) {
	var activityPO po.RaffleActivity
	err := r.defaultDB(ctx).
		Select("activity_id").
		Where("strategy_id = ?", strategyID).
		First(&activityPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return activityPO.ActivityID, nil
}

func (r *StrategyRepository) QueryRaffleActivityAccountPartakeCount(ctx context.Context, activityID int64, userID string) (int, error) {
	var accountPO po.RaffleActivityAccount
	err := r.shardDB(ctx, userID).
		Select("total_count", "total_count_surplus").
		Where("user_id = ? and activity_id = ?", userID, activityID).
		First(&accountPO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return accountPO.TotalCount - accountPO.TotalCountSurplus, nil
}
