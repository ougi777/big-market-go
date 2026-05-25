package po

import "time"

type Award struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	AwardID     int       `gorm:"column:award_id"`
	AwardKey    string    `gorm:"column:award_key"`
	AwardConfig string    `gorm:"column:award_config"`
	AwardDesc   string    `gorm:"column:award_desc"`
	CreateTime  time.Time `gorm:"column:create_time"`
	UpdateTime  time.Time `gorm:"column:update_time"`
}

func (Award) TableName() string {
	return "award"
}
