package po

import "time"

type UserRaffleOrder struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	UserID       string    `gorm:"column:user_id"`
	ActivityID   int64     `gorm:"column:activity_id"`
	ActivityName string    `gorm:"column:activity_name"`
	StrategyID   int64     `gorm:"column:strategy_id"`
	OrderID      string    `gorm:"column:order_id"`
	OrderTime    time.Time `gorm:"column:order_time"`
	OrderState   string    `gorm:"column:order_state"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

func (UserRaffleOrder) TableName() string {
	return "user_raffle_order"
}
