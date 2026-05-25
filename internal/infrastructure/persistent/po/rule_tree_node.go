package po

import "time"

type RuleTreeNode struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	TreeID     string    `gorm:"column:tree_id"`
	RuleKey    string    `gorm:"column:rule_key"`
	RuleDesc   string    `gorm:"column:rule_desc"`
	RuleValue  string    `gorm:"column:rule_value"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (RuleTreeNode) TableName() string {
	return "rule_tree_node"
}
