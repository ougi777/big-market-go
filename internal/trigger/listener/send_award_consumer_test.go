package listener

import (
	"context"
	"errors"
	"testing"

	"bm-go/internal/domain/award"
	"bm-go/internal/types"
)

func TestSendAwardConsumerHandle(t *testing.T) {
	distributor := &fakeAwardDistributor{}
	consumer := NewSendAwardConsumer(nil, distributor, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":{"userId":"xiaofuge","orderId":"order-001","awardId":101,"awardTitle":"credit","awardConfig":"0.01,1"}}`)
	if err != nil {
		t.Fatalf("handle send award: %v", err)
	}

	if distributor.distribute.UserID != "xiaofuge" ||
		distributor.distribute.OrderID != "order-001" ||
		distributor.distribute.AwardID != 101 ||
		distributor.distribute.AwardConfig != "0.01,1" {
		t.Fatalf("expected distribute award, got %+v", distributor.distribute)
	}
}

func TestSendAwardConsumerHandleIgnoresDuplicate(t *testing.T) {
	distributor := &fakeAwardDistributor{err: types.NewAppError(types.ResponseCodeIndexDup, nil)}
	consumer := NewSendAwardConsumer(nil, distributor, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":{"userId":"xiaofuge","orderId":"order-001","awardId":101,"awardTitle":"credit","awardConfig":"0.01,1"}}`)
	if err != nil {
		t.Fatalf("handle duplicate send award: %v", err)
	}
}

func TestSendAwardConsumerHandleInvalidMessage(t *testing.T) {
	distributor := &fakeAwardDistributor{}
	consumer := NewSendAwardConsumer(nil, distributor, nil)

	err := consumer.handle(context.Background(), `{invalid`)
	if err == nil {
		t.Fatal("expected parse error")
	}
	if distributor.distribute.UserID != "" {
		t.Fatalf("expected distributor not called, got %+v", distributor.distribute)
	}
}

func TestSendAwardConsumerHandleDistributorError(t *testing.T) {
	distributor := &fakeAwardDistributor{err: errors.New("distribute failed")}
	consumer := NewSendAwardConsumer(nil, distributor, nil)

	err := consumer.handle(context.Background(), `{"id":"12345678901","timestamp":1779703200000,"data":{"userId":"xiaofuge","orderId":"order-001","awardId":101,"awardTitle":"credit","awardConfig":"0.01,1"}}`)
	if err == nil {
		t.Fatal("expected distributor error")
	}
	if distributor.distribute.UserID != "xiaofuge" {
		t.Fatalf("expected distributor called, got %+v", distributor.distribute)
	}
}

func TestSendAwardConsumerStart(t *testing.T) {
	messageConsumer := &fakeMessageConsumer{}
	consumer := NewSendAwardConsumer(messageConsumer, &fakeAwardDistributor{}, nil)

	if err := consumer.Start(context.Background()); err != nil {
		t.Fatalf("start consumer: %v", err)
	}
	if messageConsumer.topic != award.TopicSendAward {
		t.Fatalf("expected topic %s, got %s", award.TopicSendAward, messageConsumer.topic)
	}
	if messageConsumer.handler == nil {
		t.Fatal("expected handler registered")
	}
}

type fakeAwardDistributor struct {
	distribute award.DistributeAwardEntity
	err        error
}

func (f *fakeAwardDistributor) DistributeAward(ctx context.Context, distribute award.DistributeAwardEntity) error {
	f.distribute = distribute
	if f.err != nil {
		return f.err
	}
	return nil
}
