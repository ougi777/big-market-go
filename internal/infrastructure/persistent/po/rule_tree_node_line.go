package po

import "time"

type RuleTreeNodeLine struct {
	ID             int64     `gorm:"column:id;primaryKey"`
	TreeID         string    `gorm:"column:tree_id"`
	RuleNodeFrom   string    `gorm:"column:rule_node_from"`
	RuleNodeTo     string    `gorm:"column:rule_node_to"`
	RuleLimitType  string    `gorm:"column:rule_limit_type"`
	RuleLimitValue string    `gorm:"column:rule_limit_value"`
	CreateTime     time.Time `gorm:"column:create_time"`
	UpdateTime     time.Time `gorm:"column:update_time"`
}

func (RuleTreeNodeLine) TableName() string {
	return "rule_tree_node_line"
}
