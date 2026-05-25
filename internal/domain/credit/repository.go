package credit

import "context"

type AccountRepository interface {
	QueryUserCreditAccount(ctx context.Context, userID string) (AccountEntity, bool, error)
}

type TradeRepository interface {
	CompleteCreditPayOrder(ctx context.Context, aggregate CompleteSkuExchangeAggregate) error
	SaveRebateIntegralOrder(ctx context.Context, rebateIntegral RebateIntegralEntity) error
}

type CompleteSkuExchangeAggregate struct {
	UserID        string
	ActivityID    int64
	TotalCount    int
	DayCount      int
	MonthCount    int
	OutBusinessNo string
	CreditOrder   OrderEntity
	SendTask      TaskEntity
}

type TaskEntity struct {
	UserID    string
	Topic     string
	MessageID string
	Message   string
	State     string
}
