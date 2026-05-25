package service

import (
	"context"
	"strings"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"
)

type exchangeRepository interface {
	QuerySkuProductBySKU(ctx context.Context, sku int64) (activity.SkuProductEntity, bool, error)
	QueryActivityByActivityID(ctx context.Context, activityID int64) (activity.ActivityEntity, bool, error)
	QueryUnpaidActivityOrder(ctx context.Context, userID string, sku int64) (activity.SkuExchangeOrderEntity, bool, error)
	SaveCreditPayOrder(ctx context.Context, aggregate activity.CreateSkuExchangeOrderAggregate) error
	CompleteCreditPayOrder(ctx context.Context, aggregate activity.CompleteSkuExchangeAggregate) error
}

type exchangeStockService interface {
	SubtractActivitySkuStock(ctx context.Context, sku int64, activityID int64) (bool, error)
}

type ExchangeService struct {
	repo                 exchangeRepository
	stockService         exchangeStockService
	now                  func() time.Time
	orderIDGenerator     func() (string, error)
	businessNoGenerator  func() (string, error)
	creditOrderGenerator func() (string, error)
}

func NewExchangeService(repo exchangeRepository, stockService exchangeStockService) *ExchangeService {
	return &ExchangeService{
		repo:                 repo,
		stockService:         stockService,
		now:                  time.Now,
		orderIDGenerator:     func() (string, error) { return randomNumeric(12) },
		businessNoGenerator:  func() (string, error) { return randomNumeric(12) },
		creditOrderGenerator: func() (string, error) { return randomNumeric(12) },
	}
}

func (s *ExchangeService) CreditPayExchangeSku(ctx context.Context, userID string, sku int64) (bool, error) {
	if strings.TrimSpace(userID) == "" || sku == 0 {
		return false, types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}

	unpaid, exists, err := s.repo.QueryUnpaidActivityOrder(ctx, userID, sku)
	if err != nil {
		return false, err
	}
	if exists {
		return s.completeCreditPayOrder(ctx, unpaid)
	}

	order, err := s.createCreditPayOrder(ctx, userID, sku)
	if err != nil {
		return false, err
	}
	return s.completeCreditPayOrder(ctx, order)
}

func (s *ExchangeService) createCreditPayOrder(ctx context.Context, userID string, sku int64) (activity.SkuExchangeOrderEntity, error) {
	product, exists, err := s.repo.QuerySkuProductBySKU(ctx, sku)
	if err != nil {
		return activity.SkuExchangeOrderEntity{}, err
	}
	if !exists {
		return activity.SkuExchangeOrderEntity{}, types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}
	activityEntity, exists, err := s.repo.QueryActivityByActivityID(ctx, product.ActivityID)
	if err != nil {
		return activity.SkuExchangeOrderEntity{}, err
	}
	if !exists || activityEntity.State != activity.ActivityStateOpen {
		return activity.SkuExchangeOrderEntity{}, types.NewAppError(types.ResponseCodeActivityStateError, nil)
	}
	currentTime := s.now()
	if activityEntity.BeginDateTime.After(currentTime) || activityEntity.EndDateTime.Before(currentTime) {
		return activity.SkuExchangeOrderEntity{}, types.NewAppError(types.ResponseCodeActivityDateError, nil)
	}

	ok, err := s.stockService.SubtractActivitySkuStock(ctx, sku, product.ActivityID)
	if err != nil {
		return activity.SkuExchangeOrderEntity{}, err
	}
	if !ok {
		return activity.SkuExchangeOrderEntity{}, types.NewAppError(types.ResponseCodeActivityStateError, nil)
	}

	orderID, err := s.orderIDGenerator()
	if err != nil {
		return activity.SkuExchangeOrderEntity{}, err
	}
	outBusinessNo, err := s.businessNoGenerator()
	if err != nil {
		return activity.SkuExchangeOrderEntity{}, err
	}
	order := activity.ActivityOrderEntity{
		UserID:        userID,
		SKU:           sku,
		ActivityID:    product.ActivityID,
		ActivityName:  activityEntity.ActivityName,
		StrategyID:    activityEntity.StrategyID,
		OrderID:       orderID,
		OrderTime:     currentTime,
		TotalCount:    product.ActivityCount.TotalCount,
		DayCount:      product.ActivityCount.DayCount,
		MonthCount:    product.ActivityCount.MonthCount,
		PayAmount:     product.ProductAmount,
		State:         activity.ActivityOrderWaitPay,
		OutBusinessNo: outBusinessNo,
	}
	if err := s.repo.SaveCreditPayOrder(ctx, activity.CreateSkuExchangeOrderAggregate{
		UserID:        userID,
		ActivityID:    product.ActivityID,
		ActivityOrder: order,
	}); err != nil {
		return activity.SkuExchangeOrderEntity{}, err
	}

	return activity.SkuExchangeOrderEntity{
		UserID:        userID,
		SKU:           sku,
		OrderID:       orderID,
		OutBusinessNo: outBusinessNo,
		PayAmount:     product.ProductAmount,
	}, nil
}

func (s *ExchangeService) completeCreditPayOrder(ctx context.Context, order activity.SkuExchangeOrderEntity) (bool, error) {
	product, exists, err := s.repo.QuerySkuProductBySKU(ctx, order.SKU)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, types.NewAppError(types.ResponseCodeIllegalParam, nil)
	}
	creditOrderID, err := s.creditOrderGenerator()
	if err != nil {
		return false, err
	}
	err = s.repo.CompleteCreditPayOrder(ctx, activity.CompleteSkuExchangeAggregate{
		UserID:        order.UserID,
		ActivityID:    product.ActivityID,
		TotalCount:    product.ActivityCount.TotalCount,
		DayCount:      product.ActivityCount.DayCount,
		MonthCount:    product.ActivityCount.MonthCount,
		OutBusinessNo: order.OutBusinessNo,
		CreditOrder: activity.CreditOrderEntity{
			UserID:        order.UserID,
			OrderID:       creditOrderID,
			TradeName:     "CONVERT_SKU",
			TradeType:     "reverse",
			TradeAmount:   -order.PayAmount,
			OutBusinessNo: order.OutBusinessNo,
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
