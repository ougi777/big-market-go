package po

import "time"

type UserAwardRecord struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	UserID     string    `gorm:"column:user_id"`
	ActivityID int64     `gorm:"column:activity_id"`
	StrategyID int64     `gorm:"column:strategy_id"`
	OrderID    string    `gorm:"column:order_id"`
	AwardID    int       `gorm:"column:award_id"`
	AwardTitle string    `gorm:"column:award_title"`
	AwardTime  time.Time `gorm:"column:award_time"`
	AwardState string    `gorm:"column:award_state"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (UserAwardRecord) TableName() string {
	return "user_award_record"
}
