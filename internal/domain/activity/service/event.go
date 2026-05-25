package service

import (
	"encoding/json"
	"time"

	"bm-go/internal/domain/award"
)

func BuildActivitySkuStockZeroMessage(sku int64) (string, error) {
	messageID, err := randomNumeric(11)
	if err != nil {
		return "", err
	}
	message, err := json.Marshal(award.EventMessage[int64]{
		ID:        messageID,
		Timestamp: time.Now().UnixMilli(),
		Data:      sku,
	})
	if err != nil {
		return "", err
	}
	return string(message), nil
}
