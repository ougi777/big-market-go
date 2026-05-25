package service

import (
	"context"
	"testing"

	"bm-go/internal/domain/activity"
)

func TestDeliveryServiceDeliverActivityOrder(t *testing.T) {
	repo := &fakeDeliveryRepository{}
	service := NewDeliveryService(repo)

	err := service.DeliverActivityOrder(context.Background(), "xiaofuge", "biz-001")
	if err != nil {
		t.Fatalf("deliver activity order: %v", err)
	}

	if repo.delivery.UserID != "xiaofuge" || repo.delivery.OutBusinessNo != "biz-001" {
		t.Fatalf("expected delivery order, got %+v", repo.delivery)
	}
}

type fakeDeliveryRepository struct {
	delivery activity.DeliveryOrderEntity
}

func (f *fakeDeliveryRepository) DeliverActivityOrder(ctx context.Context, deliveryOrder activity.DeliveryOrderEntity) error {
	f.delivery = deliveryOrder
	return nil
}
