package po

import "time"

type StrategyAward struct {
	ID                int64     `gorm:"column:id;primaryKey"`
	StrategyID        int64     `gorm:"column:strategy_id"`
	AwardID           int       `gorm:"column:award_id"`
	AwardTitle        string    `gorm:"column:award_title"`
	AwardSubtitle     string    `gorm:"column:award_subtitle"`
	AwardCount        int       `gorm:"column:award_count"`
	AwardCountSurplus int       `gorm:"column:award_count_surplus"`
	AwardRate         float64   `gorm:"column:award_rate"`
	RuleModels        string    `gorm:"column:rule_models"`
	Sort              int       `gorm:"column:sort"`
	CreateTime        time.Time `gorm:"column:create_time"`
	UpdateTime        time.Time `gorm:"column:update_time"`
}

func (StrategyAward) TableName() string {
	return "strategy_award"
}
