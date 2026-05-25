package po

import "time"

type RaffleActivityAccount struct {
	ID                int64     `gorm:"column:id;primaryKey"`
	UserID            string    `gorm:"column:user_id"`
	ActivityID        int64     `gorm:"column:activity_id"`
	TotalCount        int       `gorm:"column:total_count"`
	TotalCountSurplus int       `gorm:"column:total_count_surplus"`
	DayCount          int       `gorm:"column:day_count"`
	DayCountSurplus   int       `gorm:"column:day_count_surplus"`
	MonthCount        int       `gorm:"column:month_count"`
	MonthCountSurplus int       `gorm:"column:month_count_surplus"`
	CreateTime        time.Time `gorm:"column:create_time"`
	UpdateTime        time.Time `gorm:"column:update_time"`
}

func (RaffleActivityAccount) TableName() string {
	return "raffle_activity_account"
}
