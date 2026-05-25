package listener

import (
	"context"
	"testing"

	"bm-go/internal/domain/award"
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

type fakeAwardDistributor struct {
	distribute award.DistributeAwardEntity
}

func (f *fakeAwardDistributor) DistributeAward(ctx context.Context, distribute award.DistributeAwardEntity) error {
	f.distribute = distribute
	return nil
}
