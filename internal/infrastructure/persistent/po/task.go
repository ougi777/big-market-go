package po

import "time"

type Task struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	UserID     string    `gorm:"column:user_id"`
	Topic      string    `gorm:"column:topic"`
	MessageID  string    `gorm:"column:message_id"`
	Message    string    `gorm:"column:message"`
	State      string    `gorm:"column:state"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (Task) TableName() string {
	return "task"
}
