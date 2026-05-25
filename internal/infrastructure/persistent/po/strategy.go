package po

import "time"

type Strategy struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	StrategyID   int64     `gorm:"column:strategy_id"`
	StrategyDesc string    `gorm:"column:strategy_desc"`
	RuleModels   string    `gorm:"column:rule_models"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

func (Strategy) TableName() string {
	return "strategy"
}
