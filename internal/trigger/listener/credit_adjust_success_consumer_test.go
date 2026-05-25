package listener

import (
	"context"
	"errors"
	"testing"

	"bm-go/internal/domain/credit"
	"bm-go/internal/types"
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

func TestCreditAdjustSuccessConsumerHandleIgnoresDuplicate(t *testing.T) {
	deliverer := &fakeActivityOrderDeliverer{err: types.NewAppError(types.ResponseCodeIndexDup, nil)}
	consumer := NewCreditAdjustSuccessConsumer(nil, deliverer, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":{"userId":"xiaofuge","orderId":"order-001","amount":1.68,"outBusinessNo":"biz-001"}}`)
	if err != nil {
		t.Fatalf("handle duplicate credit adjust success: %v", err)
	}
}

func TestCreditAdjustSuccessConsumerHandleInvalidMessage(t *testing.T) {
	deliverer := &fakeActivityOrderDeliverer{}
	consumer := NewCreditAdjustSuccessConsumer(nil, deliverer, nil)

	err := consumer.handle(context.Background(), `{invalid`)
	if err == nil {
		t.Fatal("expected parse error")
	}
	if deliverer.userID != "" {
		t.Fatalf("expected deliverer not called, got %s/%s", deliverer.userID, deliverer.outBusinessNo)
	}
}

func TestCreditAdjustSuccessConsumerHandleDeliverError(t *testing.T) {
	deliverer := &fakeActivityOrderDeliverer{err: errors.New("deliver failed")}
	consumer := NewCreditAdjustSuccessConsumer(nil, deliverer, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":{"userId":"xiaofuge","orderId":"order-001","amount":1.68,"outBusinessNo":"biz-001"}}`)
	if err == nil {
		t.Fatal("expected deliver error")
	}
	if deliverer.userID != "xiaofuge" || deliverer.outBusinessNo != "biz-001" {
		t.Fatalf("expected deliverer called, got %s/%s", deliverer.userID, deliverer.outBusinessNo)
	}
}

func TestCreditAdjustSuccessConsumerStart(t *testing.T) {
	messageConsumer := &fakeMessageConsumer{}
	consumer := NewCreditAdjustSuccessConsumer(messageConsumer, &fakeActivityOrderDeliverer{}, nil)

	if err := consumer.Start(context.Background()); err != nil {
		t.Fatalf("start consumer: %v", err)
	}
	if messageConsumer.topic != credit.TopicCreditAdjustSuccess {
		t.Fatalf("expected topic %s, got %s", credit.TopicCreditAdjustSuccess, messageConsumer.topic)
	}
	if messageConsumer.handler == nil {
		t.Fatal("expected handler registered")
	}
}

type fakeActivityOrderDeliverer struct {
	userID        string
	outBusinessNo string
	err           error
}

func (f *fakeActivityOrderDeliverer) DeliverActivityOrder(ctx context.Context, userID string, outBusinessNo string) error {
	f.userID = userID
	f.outBusinessNo = outBusinessNo
	if f.err != nil {
		return f.err
	}
	return nil
}
