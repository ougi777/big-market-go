package listener

import (
	"context"
	"encoding/json"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/award"

	"go.uber.org/zap"
)

type ActivitySkuStockClearer interface {
	ClearActivitySkuStock(ctx context.Context, sku int64) error
}

type ActivitySkuStockZeroConsumer struct {
	consumer MessageConsumer
	clearer  ActivitySkuStockClearer
	logger   *zap.Logger
	cancel   context.CancelFunc
}

func NewActivitySkuStockZeroConsumer(consumer MessageConsumer, clearer ActivitySkuStockClearer, logger *zap.Logger) *ActivitySkuStockZeroConsumer {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &ActivitySkuStockZeroConsumer{
		consumer: consumer,
		clearer:  clearer,
		logger:   logger,
	}
}

func (c *ActivitySkuStockZeroConsumer) Start(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	return c.consumer.Consume(runCtx, activity.TopicActivitySkuStockZero, c.handle)
}

func (c *ActivitySkuStockZeroConsumer) Stop(ctx context.Context) error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *ActivitySkuStockZeroConsumer) handle(ctx context.Context, message string) error {
	var event award.EventMessage[int64]
	if err := json.Unmarshal([]byte(message), &event); err != nil {
		c.logger.Error("parse activity sku stock zero message failed", zap.Error(err), zap.String("message", message))
		return err
	}
	if err := c.clearer.ClearActivitySkuStock(ctx, event.Data); err != nil {
		c.logger.Error("clear activity sku stock failed", zap.Error(err), zap.Int64("sku", event.Data))
		return err
	}
	c.logger.Info("clear activity sku stock completed", zap.Int64("sku", event.Data))
	return nil
}
