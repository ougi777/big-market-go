package listener

import (
	"context"
	"testing"
)

func TestCreditAdjustSuccessConsumerHandle(t *testing.T) {
	deliverer := &fakeActivityOrderDeliverer{}
	consumer := NewCreditAdjustSuccessConsumer(nil, deliverer, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":{"userId":"xiaofuge","orderId":"order-001","amount":1.68,"outBusinessNo":"biz-001"}}`)
	if err != nil {
		t.Fatalf("handle credit adjust success: %v", err)
	}

	if deliverer.userID != "xiaofuge" || deliverer.outBusinessNo != "biz-001" {
		t.Fatalf("expected delivery request, got %s/%s", deliverer.userID, deliverer.outBusinessNo)
	}
}

type fakeActivityOrderDeliverer struct {
	userID        string
	outBusinessNo string
}

func (f *fakeActivityOrderDeliverer) DeliverActivityOrder(ctx context.Context, userID string, outBusinessNo string) error {
	f.userID = userID
	f.outBusinessNo = outBusinessNo
	return nil
}
