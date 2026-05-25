package po

import "time"

type UserBehaviorRebateOrder struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	UserID        string    `gorm:"column:user_id"`
	OrderID       string    `gorm:"column:order_id"`
	BehaviorType  string    `gorm:"column:behavior_type"`
	RebateDesc    string    `gorm:"column:rebate_desc"`
	RebateType    string    `gorm:"column:rebate_type"`
	RebateConfig  string    `gorm:"column:rebate_config"`
	OutBusinessNo string    `gorm:"column:out_business_no"`
	BizID         string    `gorm:"column:biz_id"`
	CreateTime    time.Time `gorm:"column:create_time"`
	UpdateTime    time.Time `gorm:"column:update_time"`
}

func (UserBehaviorRebateOrder) TableName() string {
	return "user_behavior_rebate_order"
}
