package award

import "time"

const (
	AwardStateCreate   = "create"
	AwardStateComplete = "complete"
	AwardStateFail     = "fail"
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
}
