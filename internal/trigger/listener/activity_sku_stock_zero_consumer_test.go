package listener

import (
	"context"
	"errors"
	"testing"

	"bm-go/internal/domain/activity"
)

func TestActivitySkuStockZeroConsumerHandle(t *testing.T) {
	clearer := &fakeActivitySkuStockClearer{}
	consumer := NewActivitySkuStockZeroConsumer(nil, clearer, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":9011}`)
	if err != nil {
		t.Fatalf("handle activity sku stock zero: %v", err)
	}

	if clearer.sku != 9011 {
		t.Fatalf("expected sku 9011, got %d", clearer.sku)
	}
}

func TestActivitySkuStockZeroConsumerHandleInvalidMessage(t *testing.T) {
	clearer := &fakeActivitySkuStockClearer{}
	consumer := NewActivitySkuStockZeroConsumer(nil, clearer, nil)

	err := consumer.handle(context.Background(), `{invalid`)
	if err == nil {
		t.Fatal("expected parse error")
	}
	if clearer.sku != 0 {
		t.Fatalf("expected clearer not called, got %d", clearer.sku)
	}
}

func TestActivitySkuStockZeroConsumerHandleClearError(t *testing.T) {
	clearer := &fakeActivitySkuStockClearer{err: errors.New("clear failed")}
	consumer := NewActivitySkuStockZeroConsumer(nil, clearer, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":9011}`)
	if err == nil {
		t.Fatal("expected clear error")
	}
	if clearer.sku != 9011 {
		t.Fatalf("expected sku 9011, got %d", clearer.sku)
	}
}

func TestActivitySkuStockZeroConsumerStart(t *testing.T) {
	messageConsumer := &fakeMessageConsumer{}
	consumer := NewActivitySkuStockZeroConsumer(messageConsumer, &fakeActivitySkuStockClearer{}, nil)

	if err := consumer.Start(context.Background()); err != nil {
		t.Fatalf("start consumer: %v", err)
	}
	if messageConsumer.topic != activity.TopicActivitySkuStockZero {
		t.Fatalf("expected topic %s, got %s", activity.TopicActivitySkuStockZero, messageConsumer.topic)
	}
	if messageConsumer.handler == nil {
		t.Fatal("expected handler registered")
	}
}

type fakeActivitySkuStockClearer struct {
	sku int64
	err error
}

func (f *fakeActivitySkuStockClearer) ClearActivitySkuStock(ctx context.Context, sku int64) error {
	f.sku = sku
	return f.err
}

type fakeMessageConsumer struct {
	topic   string
	handler func(context.Context, string) error
	err     error
}

func (f *fakeMessageConsumer) Consume(ctx context.Context, topic string, handler func(context.Context, string) error) error {
	f.topic = topic
	f.handler = handler
	return f.err
}
