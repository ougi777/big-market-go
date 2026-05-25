package award

import "time"

import taskdomain "bm-go/internal/domain/task"

const (
	AwardStateCreate       = "create"
	AwardStateComplete     = "completed"
	AwardStateFail         = "fail"
	AwardKeyUserCreditRand = "user_credit_random"
	AccountStatusOpen      = "open"
	TaskStateCreate        = taskdomain.StateCreate
	TaskStateCompleted     = taskdomain.StateCompleted
	TaskStateFail          = taskdomain.StateFail
	TopicSendAward         = "send_award"
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

type TaskEntity = taskdomain.Entity

type DistributeAwardEntity struct {
	UserID      string
	OrderID     string
	AwardID     int
	AwardConfig string
}

type UserCreditAwardEntity struct {
	UserID       string
	CreditAmount float64
}

type GiveOutPrizesAggregate struct {
	UserID          string
	UserAwardRecord UserAwardRecordEntity
	UserCreditAward UserCreditAwardEntity
}
