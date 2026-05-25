package listener

import (
	"context"
	"testing"

	"bm-go/internal/domain/rebate"
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

type fakeRebateProcessor struct {
	message rebate.SendRebateMessage
}

func (f *fakeRebateProcessor) ProcessRebate(ctx context.Context, message rebate.SendRebateMessage) error {
	f.message = message
	return nil
}
