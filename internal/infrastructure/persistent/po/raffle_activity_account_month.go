package po

import "time"

type RaffleActivityAccountMonth struct {
	ID                int64     `gorm:"column:id;primaryKey"`
	UserID            string    `gorm:"column:user_id"`
	ActivityID        int64     `gorm:"column:activity_id"`
	Month             string    `gorm:"column:month"`
	MonthCount        int       `gorm:"column:month_count"`
	MonthCountSurplus int       `gorm:"column:month_count_surplus"`
	CreateTime        time.Time `gorm:"column:create_time"`
	UpdateTime        time.Time `gorm:"column:update_time"`
}

func (RaffleActivityAccountMonth) TableName() string {
	return "raffle_activity_account_month"
}
