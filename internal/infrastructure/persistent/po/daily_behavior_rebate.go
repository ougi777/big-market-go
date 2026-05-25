package po

import "time"

type DailyBehaviorRebate struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	BehaviorType string    `gorm:"column:behavior_type"`
	RebateDesc   string    `gorm:"column:rebate_desc"`
	RebateType   string    `gorm:"column:rebate_type"`
	RebateConfig string    `gorm:"column:rebate_config"`
	State        string    `gorm:"column:state"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

func (DailyBehaviorRebate) TableName() string {
	return "daily_behavior_rebate"
}
