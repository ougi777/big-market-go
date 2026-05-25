package award

import "time"

const (
	AwardStateCreate   = "create"
	AwardStateComplete = "complete"
	AwardStateFail     = "fail"
	TaskStateCreate    = "create"
	TaskStateCompleted = "completed"
	TaskStateFail      = "fail"
	TopicSendAward     = "send_award"
)

type UserAwardRecordEntity struct {
	UserID      string
	ActivityID  int64
	StrategyID  int64
	OrderID     string
	AwardID     int
	AwardTitle  string
	AwardTime   time.Time
	AwardState  string
	AwardConfig string
	SendTask    TaskEntity
}

type SendAwardMessage struct {
	UserID      string `json:"userId"`
	OrderID     string `json:"orderId"`
	AwardID     int    `json:"awardId"`
	AwardTitle  string `json:"awardTitle"`
	AwardConfig string `json:"awardConfig,omitempty"`
}

type EventMessage[T any] struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Data      T      `json:"data"`
}

type TaskEntity struct {
	UserID    string
	Topic     string
	MessageID string
	Message   string
	State     string
}
