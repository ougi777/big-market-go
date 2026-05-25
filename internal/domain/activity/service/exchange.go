package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/credit"
	"bm-go/internal/types"
)

type exchangeRepository interface {
	QuerySkuProductBySKU(ctx context.Context, sku int64) (activity.SkuProductEntity, bool, error)
	QueryActivityByActivityID(ctx context.Context, activityID int64) (activity.ActivityEntity, bool, error)
	QueryUnpaidActivityOrder(ctx context.Context, userID string, sku int64) (activity.SkuExchangeOrderEntity, bool, error)
	SaveCreditPayOrder(ctx context.Context, aggregate activity.CreateSkuExchangeOrderAggregate) error
	UpdateTaskSendMessageCompleted(ctx context.Context, userID string, messageID string) error
	UpdateTaskSendMessageFail(ctx context.Context, userID string, messageID string) error
}

type creditTradeRepository interface {
	CompleteCreditPayOrder(ctx context.Context, aggregate credit.CompleteSkuExchangeAggregate) error
}

type exchangeStockService interface {
	SubtractActivitySkuStock(ctx context.Context, sku int64, activityID int64) (bool, error)
}

type ExchangeService struct {
	repo                 exchangeRepository
	creditRepo           creditTradeRepository
	stockService         exchangeStockService
	publisher            credit.MessagePublisher
	now                  func() time.Time
	orderIDGenerator     func() (string, error)
	messageIDGenerator   func() (string, error)
	businessNoGenerator  func() (string, error)
	creditOrderGenerator func() (string, error)
}

func NewExchangeService(repo exchangeRepository, stockService exchangeStockService, publisher credit.MessagePublisher, creditRepos ...creditTradeRepository) *ExchangeService {
	var creditRepo creditTradeRepository
	if len(creditRepos) > 0 {
		creditRepo = creditRepos[0]
	}
	if creditRepo == nil {
		creditRepo = missingCreditTradeRepository{}
	}
	return &ExchangeService{
		repo:                 repo,
		creditRepo:           creditRepo,
		stockService:         stockService,
		publisher:            publisher,
		now:                  time.Now,
		orderIDGenerator:     func() (string, error) { return randomNumeric(12) },
		messageIDGenerator:   func() (string, error) { return randomNumeric(11) },
		businessNoGenerator:  func() (string, error) { return randomNumeric(12) },
		creditOrderGenerator: func() (string, error) { return randomNumeric(12) },
	}
}

type missingCreditTradeRepository struct{}

func (missingCreditTradeRepository) CompleteCreditPayOrder(ctx context.Context, aggregate credit.CompleteSkuExchangeAggregate) error {
	return errors.New("credit trade repository is not configured")
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
		return s.payCreditOrder(ctx, unpaid)
	}

	order, err := s.createCreditPayOrder(ctx, userID, sku)
	if err != nil {
		return false, err
	}
	return s.payCreditOrder(ctx, order)
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

func (s *ExchangeService) payCreditOrder(ctx context.Context, order activity.SkuExchangeOrderEntity) (bool, error) {
	creditOrderID, err := s.creditOrderGenerator()
	if err != nil {
		return false, err
	}
	messageID, err := s.messageIDGenerator()
	if err != nil {
		return false, err
	}
	message, err := json.Marshal(credit.EventMessage[credit.AdjustSuccessMessage]{
		ID:        messageID,
		Timestamp: s.now().UnixMilli(),
		Data: credit.AdjustSuccessMessage{
			UserID:        order.UserID,
			OrderID:       creditOrderID,
			Amount:        order.PayAmount,
			OutBusinessNo: order.OutBusinessNo,
		},
	})
	if err != nil {
		return false, err
	}
	err = s.creditRepo.CompleteCreditPayOrder(ctx, credit.CompleteSkuExchangeAggregate{
		UserID:        order.UserID,
		OutBusinessNo: order.OutBusinessNo,
		CreditOrder: credit.OrderEntity{
			UserID:        order.UserID,
			OrderID:       creditOrderID,
			TradeName:     "CONVERT_SKU",
			TradeType:     "reverse",
			TradeAmount:   -order.PayAmount,
			OutBusinessNo: order.OutBusinessNo,
		},
		SendTask: credit.TaskEntity{
			UserID:    order.UserID,
			Topic:     credit.TopicCreditAdjustSuccess,
			MessageID: messageID,
			Message:   string(message),
			State:     "create",
		},
	})
	if err != nil {
		return false, err
	}
	if err := s.publisher.Publish(ctx, credit.TopicCreditAdjustSuccess, string(message)); err != nil {
		_ = s.repo.UpdateTaskSendMessageFail(ctx, order.UserID, messageID)
		return false, err
	}
	_ = s.repo.UpdateTaskSendMessageCompleted(ctx, order.UserID, messageID)
	return true, nil
}
