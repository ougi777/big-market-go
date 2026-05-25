package po

import "time"

type UserCreditOrder struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	UserID        string    `gorm:"column:user_id"`
	OrderID       string    `gorm:"column:order_id"`
	TradeName     string    `gorm:"column:trade_name"`
	TradeType     string    `gorm:"column:trade_type"`
	TradeAmount   float64   `gorm:"column:trade_amount"`
	OutBusinessNo string    `gorm:"column:out_business_no"`
	CreateTime    time.Time `gorm:"column:create_time"`
	UpdateTime    time.Time `gorm:"column:update_time"`
}

func (UserCreditOrder) TableName() string {
	return "user_credit_order"
}
