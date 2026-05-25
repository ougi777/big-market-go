package po

import "time"

type UserCreditAccount struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	UserID          string    `gorm:"column:user_id"`
	TotalAmount     float64   `gorm:"column:total_amount"`
	AvailableAmount float64   `gorm:"column:available_amount"`
	AccountStatus   string    `gorm:"column:account_status"`
	CreateTime      time.Time `gorm:"column:create_time"`
	UpdateTime      time.Time `gorm:"column:update_time"`
}

func (UserCreditAccount) TableName() string {
	return "user_credit_account"
}
