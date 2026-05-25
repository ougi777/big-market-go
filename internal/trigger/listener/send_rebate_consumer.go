package listener

import (
	"context"
	"encoding/json"

	"bm-go/internal/domain/rebate"

	"go.uber.org/zap"
)

type RebateProcessor interface {
	ProcessRebate(ctx context.Context, message rebate.SendRebateMessage) error
}

type SendRebateConsumer struct {
	consumer  MessageConsumer
	processor RebateProcessor
	logger    *zap.Logger
	cancel    context.CancelFunc
}

func NewSendRebateConsumer(consumer MessageConsumer, processor RebateProcessor, logger *zap.Logger) *SendRebateConsumer {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &SendRebateConsumer{
		consumer:  consumer,
		processor: processor,
		logger:    logger,
	}
}

func (c *SendRebateConsumer) Start(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	return c.consumer.Consume(runCtx, rebate.TopicSendRebate, c.handle)
}

func (c *SendRebateConsumer) Stop(ctx context.Context) error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *SendRebateConsumer) handle(ctx context.Context, message string) error {
	var event rebate.EventMessage[rebate.SendRebateMessage]
	if err := json.Unmarshal([]byte(message), &event); err != nil {
		c.logger.Error("parse send rebate message failed", zap.Error(err), zap.String("message", message))
		return err
	}

	if err := c.processor.ProcessRebate(ctx, event.Data); err != nil {
		if isIndexDuplicateError(err) {
			c.logger.Info("process rebate ignored duplicate", zap.String("userId", event.Data.UserID), zap.String("bizId", event.Data.BizID))
			return nil
		}
		c.logger.Error("process rebate failed", zap.Error(err), zap.String("message", message))
		return err
	}
	c.logger.Info("process rebate completed", zap.String("userId", event.Data.UserID), zap.String("bizId", event.Data.BizID))
	return nil
}
