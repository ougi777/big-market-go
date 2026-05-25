package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/rebate"
	"bm-go/internal/types"
)

type rebateRepository interface {
	QuerySkuProductBySKU(ctx context.Context, sku int64) (activity.SkuProductEntity, bool, error)
	QueryActivityByActivityID(ctx context.Context, activityID int64) (activity.ActivityEntity, bool, error)
	SaveRebateSkuOrder(ctx context.Context, aggregate activity.CreateRebateSkuOrderAggregate) error
	SaveRebateIntegralOrder(ctx context.Context, rebateIntegral activity.RebateIntegralEntity) error
}

type RebateProcessor struct {
	repo             rebateRepository
	now              func() time.Time
	orderIDGenerator func() (string, error)
}

func NewRebateProcessor(repo rebateRepository) *RebateProcessor {
	return &RebateProcessor{
		repo:             repo,
		now:              time.Now,
		orderIDGenerator: func() (string, error) { return randomNumeric(12) },
	}
}

func (p *RebateProcessor) ProcessRebate(ctx context.Context, message rebate.SendRebateMessage) error {
	if strings.TrimSpace(message.UserID) == "" || strings.TrimSpace(message.BizID) == "" {
		return types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}

	switch message.RebateType {
	case rebate.RebateTypeSKU:
		return p.processSKU(ctx, message)
	case rebate.RebateTypeIntegral:
		return p.processIntegral(ctx, message)
	default:
		return types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}
}

func (p *RebateProcessor) processSKU(ctx context.Context, message rebate.SendRebateMessage) error {
	sku, err := strconv.ParseInt(message.RebateConfig, 10, 64)
	if err != nil || sku == 0 {
		return types.NewAppError(types.ResponseCodeIllegalParam, err)
	}
	product, exists, err := p.repo.QuerySkuProductBySKU(ctx, sku)
	if err != nil {
		return err
	}
	if !exists {
		return types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}
	activityEntity, exists, err := p.repo.QueryActivityByActivityID(ctx, product.ActivityID)
	if err != nil {
		return err
	}
	if !exists || activityEntity.State != activity.ActivityStateOpen {
		return types.NewAppError(types.ResponseCodeActivityStateError, nil)
	}

	now := p.now()
	orderID, err := p.orderIDGenerator()
	if err != nil {
		return err
	}
	return p.repo.SaveRebateSkuOrder(ctx, activity.CreateRebateSkuOrderAggregate{
		UserID:     message.UserID,
		ActivityID: product.ActivityID,
		ActivityOrder: activity.ActivityOrderEntity{
			UserID:        message.UserID,
			SKU:           sku,
			ActivityID:    product.ActivityID,
			ActivityName:  activityEntity.ActivityName,
			StrategyID:    activityEntity.StrategyID,
			OrderID:       orderID,
			OrderTime:     now,
			TotalCount:    product.ActivityCount.TotalCount,
			DayCount:      product.ActivityCount.DayCount,
			MonthCount:    product.ActivityCount.MonthCount,
			PayAmount:     0,
			State:         activity.ActivityOrderCompleted,
			OutBusinessNo: message.BizID,
		},
	})
}

func (p *RebateProcessor) processIntegral(ctx context.Context, message rebate.SendRebateMessage) error {
	amount, err := strconv.ParseFloat(message.RebateConfig, 64)
	if err != nil || amount <= 0 {
		return types.NewAppError(types.ResponseCodeIllegalParam, err)
	}
	orderID, err := p.orderIDGenerator()
	if err != nil {
		return err
	}
	return p.repo.SaveRebateIntegralOrder(ctx, activity.RebateIntegralEntity{
		UserID:        message.UserID,
		OrderID:       orderID,
		TradeAmount:   amount,
		OutBusinessNo: message.BizID,
	})
}
