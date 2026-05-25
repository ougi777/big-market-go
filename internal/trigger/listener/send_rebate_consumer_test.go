package listener

import (
	"context"
	"errors"
	"testing"

	"bm-go/internal/domain/rebate"
	"bm-go/internal/types"
)

func TestSendRebateConsumerHandle(t *testing.T) {
	processor := &fakeRebateProcessor{}
	consumer := NewSendRebateConsumer(nil, processor, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":{"userId":"xiaofuge","rebateType":"sku","rebateConfig":"9011","bizId":"xiaofuge_sku_20260525"}}`)
	if err != nil {
		t.Fatalf("handle send rebate: %v", err)
	}

	if processor.message.UserID != "xiaofuge" ||
		processor.message.RebateType != rebate.RebateTypeSKU ||
		processor.message.RebateConfig != "9011" ||
		processor.message.BizID != "xiaofuge_sku_20260525" {
		t.Fatalf("expected rebate message, got %+v", processor.message)
	}
}

func TestSendRebateConsumerHandleIgnoresDuplicate(t *testing.T) {
	processor := &fakeRebateProcessor{err: types.NewAppError(types.ResponseCodeIndexDup, nil)}
	consumer := NewSendRebateConsumer(nil, processor, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":{"userId":"xiaofuge","rebateType":"sku","rebateConfig":"9011","bizId":"xiaofuge_sku_20260525"}}`)
	if err != nil {
		t.Fatalf("handle duplicate send rebate: %v", err)
	}
}

func TestSendRebateConsumerHandleInvalidMessage(t *testing.T) {
	processor := &fakeRebateProcessor{}
	consumer := NewSendRebateConsumer(nil, processor, nil)

	err := consumer.handle(context.Background(), `{invalid`)
	if err == nil {
		t.Fatal("expected parse error")
	}
	if processor.message.UserID != "" {
		t.Fatalf("expected processor not called, got %+v", processor.message)
	}
}

func TestSendRebateConsumerHandleProcessorError(t *testing.T) {
	processor := &fakeRebateProcessor{err: errors.New("process failed")}
	consumer := NewSendRebateConsumer(nil, processor, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":{"userId":"xiaofuge","rebateType":"sku","rebateConfig":"9011","bizId":"xiaofuge_sku_20260525"}}`)
	if err == nil {
		t.Fatal("expected processor error")
	}
	if processor.message.UserID != "xiaofuge" {
		t.Fatalf("expected processor called, got %+v", processor.message)
	}
}

func TestSendRebateConsumerStart(t *testing.T) {
	messageConsumer := &fakeMessageConsumer{}
	consumer := NewSendRebateConsumer(messageConsumer, &fakeRebateProcessor{}, nil)

	if err := consumer.Start(context.Background()); err != nil {
		t.Fatalf("start consumer: %v", err)
	}
	if messageConsumer.topic != rebate.TopicSendRebate {
		t.Fatalf("expected topic %s, got %s", rebate.TopicSendRebate, messageConsumer.topic)
	}
	if messageConsumer.handler == nil {
		t.Fatal("expected handler registered")
	}
}

type fakeRebateProcessor struct {
	message rebate.SendRebateMessage
	err     error
}

func (f *fakeRebateProcessor) ProcessRebate(ctx context.Context, message rebate.SendRebateMessage) error {
	f.message = message
	if f.err != nil {
		return f.err
	}
	return nil
}
