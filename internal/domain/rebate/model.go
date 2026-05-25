package rebate

const (
	BehaviorTypeSign   = "sign"
	RebateStateOpen    = "open"
	RebateTypeSKU      = "sku"
	RebateTypeIntegral = "integral"
	TopicSendRebate    = "send_rebate"
	TaskStateCreate    = "create"
	TaskStateComplete  = "completed"
	TaskStateFail      = "fail"
)

type DailyBehaviorRebateEntity struct {
	BehaviorType string
	RebateDesc   string
	RebateType   string
	RebateConfig string
}

type BehaviorRebateOrderEntity struct {
	UserID        string
	OrderID       string
	BehaviorType  string
	RebateDesc    string
	RebateType    string
	RebateConfig  string
	OutBusinessNo string
	BizID         string
}

type SendRebateMessage struct {
	UserID       string `json:"userId"`
	RebateDesc   string `json:"rebateDesc,omitempty"`
	RebateType   string `json:"rebateType"`
	RebateConfig string `json:"rebateConfig"`
	BizID        string `json:"bizId"`
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

type BehaviorRebateAggregate struct {
	UserID string
	Order  BehaviorRebateOrderEntity
	Task   TaskEntity
}
