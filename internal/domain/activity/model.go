package activity

import "time"

const (
	ActivityStateOpen         = "open"
	ActivityOrderWaitPay      = "wait_pay"
	ActivityOrderCompleted    = "completed"
	ActivityTradeRebateNoPay  = "rebate_no_pay_trade"
	TopicActivitySkuStockZero = "activity_sku_stock_zero"
	UserRaffleOrderCreate     = "create"
	UserRaffleOrderUsed       = "used"
	UserRaffleOrderCancel     = "cancel"
)

type ActivityEntity struct {
	ActivityID    int64
	ActivityName  string
	ActivityDesc  string
	BeginDateTime time.Time
	EndDateTime   time.Time
	StrategyID    int64
	State         string
}

type AccountEntity struct {
	UserID            string
	ActivityID        int64
	TotalCount        int
	TotalCountSurplus int
	DayCount          int
	DayCountSurplus   int
	MonthCount        int
	MonthCountSurplus int
}

type AccountDayEntity struct {
	UserID          string
	ActivityID      int64
	Day             string
	DayCount        int
	DayCountSurplus int
}

type AccountMonthEntity struct {
	UserID            string
	ActivityID        int64
	Month             string
	MonthCount        int
	MonthCountSurplus int
}

type SkuProductEntity struct {
	SKU               int64
	ActivityID        int64
	ActivityCountID   int64
	StockCount        int
	StockCountSurplus int
	ProductAmount     float64
	ActivityCount     ActivityCountEntity
}

type CreditAccountEntity struct {
	UserID          string
	AvailableAmount float64
}

type ActivityOrderEntity struct {
	UserID        string
	SKU           int64
	ActivityID    int64
	ActivityName  string
	StrategyID    int64
	OrderID       string
	OrderTime     time.Time
	TotalCount    int
	DayCount      int
	MonthCount    int
	PayAmount     float64
	State         string
	OutBusinessNo string
}

type CreditOrderEntity struct {
	UserID        string
	OrderID       string
	TradeName     string
	TradeType     string
	TradeAmount   float64
	OutBusinessNo string
}

type SkuExchangeOrderEntity struct {
	UserID        string
	SKU           int64
	OrderID       string
	OutBusinessNo string
	PayAmount     float64
}

type CreateSkuExchangeOrderAggregate struct {
	UserID        string
	ActivityID    int64
	ActivityOrder ActivityOrderEntity
}

type CompleteSkuExchangeAggregate struct {
	UserID        string
	ActivityID    int64
	TotalCount    int
	DayCount      int
	MonthCount    int
	OutBusinessNo string
	CreditOrder   CreditOrderEntity
	SendTask      TaskEntity
}

type RebateSkuOrderEntity struct {
	UserID        string
	SKU           int64
	OrderID       string
	OutBusinessNo string
}

type CreateRebateSkuOrderAggregate struct {
	UserID        string
	ActivityID    int64
	ActivityOrder ActivityOrderEntity
}

type RebateIntegralEntity struct {
	UserID        string
	OrderID       string
	TradeAmount   float64
	OutBusinessNo string
}

type DeliveryOrderEntity struct {
	UserID        string
	OutBusinessNo string
}

type TaskEntity struct {
	UserID    string
	Topic     string
	MessageID string
	Message   string
	State     string
}

type ActivitySkuStockKey struct {
	SKU        int64 `json:"sku"`
	ActivityID int64 `json:"activityId"`
}

type ActivityCountEntity struct {
	ActivityCountID int64
	TotalCount      int
	DayCount        int
	MonthCount      int
}

type UserRaffleOrderEntity struct {
	UserID       string
	ActivityID   int64
	ActivityName string
	StrategyID   int64
	OrderID      string
	OrderTime    time.Time
	OrderState   string
	EndDateTime  time.Time
}

type CreatePartakeOrderAggregate struct {
	UserID               string
	ActivityID           int64
	ActivityAccount      AccountEntity
	ExistAccountMonth    bool
	ActivityAccountMonth AccountMonthEntity
	ExistAccountDay      bool
	ActivityAccountDay   AccountDayEntity
	UserRaffleOrder      UserRaffleOrderEntity
}

type DrawResult struct {
	AwardID    int
	AwardTitle string
	AwardIndex int
}
