package service

import (
	"context"
	"strings"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"
)

type deliveryRepository interface {
	DeliverActivityOrder(ctx context.Context, deliveryOrder activity.DeliveryOrderEntity) error
}

type DeliveryService struct {
	repo deliveryRepository
}

func NewDeliveryService(repo deliveryRepository) *DeliveryService {
	return &DeliveryService{repo: repo}
}

func (s *DeliveryService) DeliverActivityOrder(ctx context.Context, userID string, outBusinessNo string) error {
	userID = strings.TrimSpace(userID)
	outBusinessNo = strings.TrimSpace(outBusinessNo)
	if userID == "" || outBusinessNo == "" {
		return types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}
	return s.repo.DeliverActivityOrder(ctx, activity.DeliveryOrderEntity{
		UserID:        userID,
		OutBusinessNo: outBusinessNo,
	})
}
