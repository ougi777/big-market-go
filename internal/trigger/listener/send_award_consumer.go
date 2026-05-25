package listener

import (
	"context"
	"encoding/json"

	"bm-go/internal/domain/award"

	"go.uber.org/zap"
)

type AwardDistributor interface {
	DistributeAward(ctx context.Context, distribute award.DistributeAwardEntity) error
}

type MessageConsumer interface {
	Consume(ctx context.Context, topic string, handler func(context.Context, string) error) error
}

type SendAwardConsumer struct {
	consumer    MessageConsumer
	distributor AwardDistributor
	logger      *zap.Logger
	cancel      context.CancelFunc
}

func NewSendAwardConsumer(consumer MessageConsumer, distributor AwardDistributor, logger *zap.Logger) *SendAwardConsumer {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &SendAwardConsumer{
		consumer:    consumer,
		distributor: distributor,
		logger:      logger,
	}
}

func (c *SendAwardConsumer) Start(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	return c.consumer.Consume(runCtx, award.TopicSendAward, c.handle)
}

func (c *SendAwardConsumer) Stop(ctx context.Context) error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *SendAwardConsumer) handle(ctx context.Context, message string) error {
	var event award.EventMessage[award.SendAwardMessage]
	if err := json.Unmarshal([]byte(message), &event); err != nil {
		c.logger.Error("parse send award message failed", zap.Error(err), zap.String("message", message))
		return err
	}

	err := c.distributor.DistributeAward(ctx, award.DistributeAwardEntity{
		UserID:      event.Data.UserID,
		OrderID:     event.Data.OrderID,
		AwardID:     event.Data.AwardID,
		AwardConfig: event.Data.AwardConfig,
	})
	if err != nil {
		c.logger.Error("distribute award failed", zap.Error(err), zap.String("message", message))
		return err
	}
	c.logger.Info("distribute award completed", zap.String("userId", event.Data.UserID), zap.String("orderId", event.Data.OrderID))
	return nil
}
