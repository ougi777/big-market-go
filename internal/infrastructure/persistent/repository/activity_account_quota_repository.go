package repository

import (
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/credit"
	"bm-go/internal/infrastructure/persistent/po"
	"bm-go/internal/types"

	"gorm.io/gorm"
)

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

func addActivityAccountQuota(tx *gorm.DB, aggregate credit.CompleteSkuExchangeAggregate, now time.Time) error {
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
