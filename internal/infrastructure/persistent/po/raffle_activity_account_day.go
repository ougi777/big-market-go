package po

import "time"

type RaffleActivityAccountDay struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	UserID          string    `gorm:"column:user_id"`
	ActivityID      int64     `gorm:"column:activity_id"`
	Day             string    `gorm:"column:day"`
	DayCount        int       `gorm:"column:day_count"`
	DayCountSurplus int       `gorm:"column:day_count_surplus"`
	CreateTime      time.Time `gorm:"column:create_time"`
	UpdateTime      time.Time `gorm:"column:update_time"`
}

func (RaffleActivityAccountDay) TableName() string {
	return "raffle_activity_account_day"
}
