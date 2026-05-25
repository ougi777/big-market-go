package po

import "time"

type RaffleActivity struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	ActivityID    int64     `gorm:"column:activity_id"`
	ActivityName  string    `gorm:"column:activity_name"`
	ActivityDesc  string    `gorm:"column:activity_desc"`
	BeginDateTime time.Time `gorm:"column:begin_date_time"`
	EndDateTime   time.Time `gorm:"column:end_date_time"`
	StrategyID    int64     `gorm:"column:strategy_id"`
	State         string    `gorm:"column:state"`
	CreateTime    time.Time `gorm:"column:create_time"`
	UpdateTime    time.Time `gorm:"column:update_time"`
}

func (RaffleActivity) TableName() string {
	return "raffle_activity"
}
