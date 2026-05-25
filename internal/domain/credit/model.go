package credit

import "context"

const TopicCreditAdjustSuccess = "credit_adjust_success"

type AccountEntity struct {
	UserID          string
	AvailableAmount float64
}

type OrderEntity struct {
	UserID        string
	OrderID       string
	TradeName     string
	TradeType     string
	TradeAmount   float64
	OutBusinessNo string
}

type RebateIntegralEntity struct {
	UserID        string
	OrderID       string
	TradeAmount   float64
	OutBusinessNo string
}

type MessagePublisher interface {
	Publish(ctx context.Context, topic string, message string) error
}

type EventMessage[T any] struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Data      T      `json:"data"`
}

type AdjustSuccessMessage struct {
	UserID        string  `json:"userId"`
	OrderID       string  `json:"orderId"`
	Amount        float64 `json:"amount"`
	OutBusinessNo string  `json:"outBusinessNo"`
}
