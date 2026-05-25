package listener

import (
	"context"
	"encoding/json"

	"bm-go/internal/domain/credit"

	"go.uber.org/zap"
)

type ActivityOrderDeliverer interface {
	DeliverActivityOrder(ctx context.Context, userID string, outBusinessNo string) error
}

type CreditAdjustSuccessConsumer struct {
	consumer  MessageConsumer
	deliverer ActivityOrderDeliverer
	logger    *zap.Logger
	cancel    context.CancelFunc
}

func NewCreditAdjustSuccessConsumer(consumer MessageConsumer, deliverer ActivityOrderDeliverer, logger *zap.Logger) *CreditAdjustSuccessConsumer {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &CreditAdjustSuccessConsumer{
		consumer:  consumer,
		deliverer: deliverer,
		logger:    logger,
	}
}

func (c *CreditAdjustSuccessConsumer) Start(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	return c.consumer.Consume(runCtx, credit.TopicCreditAdjustSuccess, c.handle)
}

func (c *CreditAdjustSuccessConsumer) Stop(ctx context.Context) error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *CreditAdjustSuccessConsumer) handle(ctx context.Context, message string) error {
	var event credit.EventMessage[credit.AdjustSuccessMessage]
	if err := json.Unmarshal([]byte(message), &event); err != nil {
		c.logger.Error("parse credit adjust success message failed", zap.Error(err), zap.String("message", message))
		return err
	}

	if err := c.deliverer.DeliverActivityOrder(ctx, event.Data.UserID, event.Data.OutBusinessNo); err != nil {
		if isIndexDuplicateError(err) {
			c.logger.Info("deliver activity order ignored duplicate", zap.String("userId", event.Data.UserID), zap.String("outBusinessNo", event.Data.OutBusinessNo))
			return nil
		}
		c.logger.Error("deliver activity order failed", zap.Error(err), zap.String("message", message))
		return err
	}
	c.logger.Info("deliver activity order completed", zap.String("userId", event.Data.UserID), zap.String("outBusinessNo", event.Data.OutBusinessNo))
	return nil
}
