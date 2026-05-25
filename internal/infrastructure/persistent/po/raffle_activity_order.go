package po

import "time"

type RaffleActivityOrder struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	UserID        string    `gorm:"column:user_id"`
	SKU           int64     `gorm:"column:sku"`
	ActivityID    int64     `gorm:"column:activity_id"`
	ActivityName  string    `gorm:"column:activity_name"`
	StrategyID    int64     `gorm:"column:strategy_id"`
	OrderID       string    `gorm:"column:order_id"`
	OrderTime     time.Time `gorm:"column:order_time"`
	TotalCount    int       `gorm:"column:total_count"`
	DayCount      int       `gorm:"column:day_count"`
	MonthCount    int       `gorm:"column:month_count"`
	PayAmount     float64   `gorm:"column:pay_amount"`
	State         string    `gorm:"column:state"`
	OutBusinessNo string    `gorm:"column:out_business_no"`
	CreateTime    time.Time `gorm:"column:create_time"`
	UpdateTime    time.Time `gorm:"column:update_time"`
}

func (RaffleActivityOrder) TableName() string {
	return "raffle_activity_order"
}
