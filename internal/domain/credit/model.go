package credit

const TopicCreditAdjustSuccess = "credit_adjust_success"

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
