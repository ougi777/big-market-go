package service

import (
	"context"
	"errors"
	"testing"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"
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

func TestDeliveryServiceDeliverActivityOrderTrimSpace(t *testing.T) {
	repo := &fakeDeliveryRepository{}
	service := NewDeliveryService(repo)

	err := service.DeliverActivityOrder(context.Background(), " xiaofuge ", " biz-001 ")
	if err != nil {
		t.Fatalf("deliver activity order: %v", err)
	}

	if repo.delivery.UserID != "xiaofuge" || repo.delivery.OutBusinessNo != "biz-001" {
		t.Fatalf("expected trimmed delivery order, got %+v", repo.delivery)
	}
}

func TestDeliveryServiceDeliverActivityOrderIllegalParam(t *testing.T) {
	service := NewDeliveryService(&fakeDeliveryRepository{})

	err := service.DeliverActivityOrder(context.Background(), "", "biz-001")
	assertAppErrorCode(t, err, types.ResponseCodeIllegalParam)
}

func TestDeliveryServiceDeliverActivityOrderRepositoryError(t *testing.T) {
	expectedErr := errors.New("deliver failed")
	service := NewDeliveryService(&fakeDeliveryRepository{err: expectedErr})

	err := service.DeliverActivityOrder(context.Background(), "xiaofuge", "biz-001")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

type fakeDeliveryRepository struct {
	delivery activity.DeliveryOrderEntity
	err      error
}

func (f *fakeDeliveryRepository) DeliverActivityOrder(ctx context.Context, deliveryOrder activity.DeliveryOrderEntity) error {
	f.delivery = deliveryOrder
	if f.err != nil {
		return f.err
	}
	return nil
}
