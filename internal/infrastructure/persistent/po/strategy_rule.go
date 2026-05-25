package po

import "time"

type StrategyRule struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	StrategyID int64     `gorm:"column:strategy_id"`
	AwardID    *int      `gorm:"column:award_id"`
	RuleType   int       `gorm:"column:rule_type"`
	RuleModel  string    `gorm:"column:rule_model"`
	RuleValue  string    `gorm:"column:rule_value"`
	RuleDesc   string    `gorm:"column:rule_desc"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (StrategyRule) TableName() string {
	return "strategy_rule"
}
