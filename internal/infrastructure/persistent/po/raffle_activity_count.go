package po

import "time"

type RaffleActivityCount struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	ActivityCountID int64     `gorm:"column:activity_count_id"`
	TotalCount      int       `gorm:"column:total_count"`
	DayCount        int       `gorm:"column:day_count"`
	MonthCount      int       `gorm:"column:month_count"`
	CreateTime      time.Time `gorm:"column:create_time"`
	UpdateTime      time.Time `gorm:"column:update_time"`
}

func (RaffleActivityCount) TableName() string {
	return "raffle_activity_count"
}
