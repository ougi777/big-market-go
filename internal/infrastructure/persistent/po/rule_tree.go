package po

import "time"

type RuleTree struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	TreeID          string    `gorm:"column:tree_id"`
	TreeName        string    `gorm:"column:tree_name"`
	TreeDesc        string    `gorm:"column:tree_desc"`
	TreeRootRuleKey string    `gorm:"column:tree_node_rule_key"`
	CreateTime      time.Time `gorm:"column:create_time"`
	UpdateTime      time.Time `gorm:"column:update_time"`
}

func (RuleTree) TableName() string {
	return "rule_tree"
}
