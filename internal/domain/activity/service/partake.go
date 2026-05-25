package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"
)

type PartakeService struct {
	repo             activity.PartakeRepository
	now              func() time.Time
	orderIDGenerator func() (string, error)
}

func NewPartakeService(repo activity.PartakeRepository) *PartakeService {
	return &PartakeService{
		repo:             repo,
		now:              time.Now,
		orderIDGenerator: func() (string, error) { return randomNumeric(12) },
	}
}

func (s *PartakeService) CreateOrder(ctx context.Context, userID string, activityID int64) (activity.UserRaffleOrderEntity, error) {
	currentTime := s.now()

	activityEntity, exists, err := s.repo.QueryActivityByActivityID(ctx, activityID)
	if err != nil {
		return activity.UserRaffleOrderEntity{}, err
	}
	if !exists || activityEntity.State != activity.ActivityStateOpen {
		return activity.UserRaffleOrderEntity{}, types.NewAppError(types.ResponseCodeActivityStateError, nil)
	}
	if activityEntity.BeginDateTime.After(currentTime) || activityEntity.EndDateTime.Before(currentTime) {
		return activity.UserRaffleOrderEntity{}, types.NewAppError(types.ResponseCodeActivityDateError, nil)
	}

	existsOrder, exists, err := s.repo.QueryNoUsedRaffleOrder(ctx, userID, activityID)
	if err != nil {
		return activity.UserRaffleOrderEntity{}, err
	}
	if exists {
		return existsOrder, nil
	}

	aggregate, err := s.buildPartakeAggregate(ctx, userID, activityID, currentTime)
	if err != nil {
		return activity.UserRaffleOrderEntity{}, err
	}

	orderID, err := s.orderIDGenerator()
	if err != nil {
		return activity.UserRaffleOrderEntity{}, err
	}
	aggregate.UserRaffleOrder = activity.UserRaffleOrderEntity{
		UserID:       userID,
		ActivityID:   activityID,
		ActivityName: activityEntity.ActivityName,
		StrategyID:   activityEntity.StrategyID,
		OrderID:      orderID,
		OrderTime:    currentTime,
		OrderState:   activity.UserRaffleOrderCreate,
		EndDateTime:  activityEntity.EndDateTime,
	}

	if err := s.repo.SaveCreatePartakeOrder(ctx, aggregate); err != nil {
		return activity.UserRaffleOrderEntity{}, err
	}
	return aggregate.UserRaffleOrder, nil
}

func (s *PartakeService) buildPartakeAggregate(ctx context.Context, userID string, activityID int64, currentTime time.Time) (activity.CreatePartakeOrderAggregate, error) {
	account, exists, err := s.repo.QueryActivityAccount(ctx, activityID, userID)
	if err != nil {
		return activity.CreatePartakeOrderAggregate{}, err
	}
	if !exists || account.TotalCountSurplus <= 0 {
		return activity.CreatePartakeOrderAggregate{}, types.NewAppError(types.ResponseCodeAccountQuotaError, nil)
	}

	month := currentTime.Format("2006-01")
	monthAccount, existMonth, err := s.repo.QueryActivityAccountMonth(ctx, activityID, userID, month)
	if err != nil {
		return activity.CreatePartakeOrderAggregate{}, err
	}
	if existMonth && monthAccount.MonthCountSurplus <= 0 {
		return activity.CreatePartakeOrderAggregate{}, types.NewAppError(types.ResponseCodeAccountMonthQuotaError, nil)
	}
	if !existMonth {
		monthAccount = activity.AccountMonthEntity{
			UserID:            userID,
			ActivityID:        activityID,
			Month:             month,
			MonthCount:        account.MonthCount,
			MonthCountSurplus: account.MonthCount,
		}
	}

	day := currentTime.Format("2006-01-02")
	dayAccount, existDay, err := s.repo.QueryActivityAccountDay(ctx, activityID, userID, day)
	if err != nil {
		return activity.CreatePartakeOrderAggregate{}, err
	}
	if existDay && dayAccount.DayCountSurplus <= 0 {
		return activity.CreatePartakeOrderAggregate{}, types.NewAppError(types.ResponseCodeAccountDayQuotaError, nil)
	}
	if !existDay {
		dayAccount = activity.AccountDayEntity{
			UserID:          userID,
			ActivityID:      activityID,
			Day:             day,
			DayCount:        account.DayCount,
			DayCountSurplus: account.DayCount,
		}
	}

	return activity.CreatePartakeOrderAggregate{
		UserID:               userID,
		ActivityID:           activityID,
		ActivityAccount:      account,
		ExistAccountMonth:    existMonth,
		ActivityAccountMonth: monthAccount,
		ExistAccountDay:      existDay,
		ActivityAccountDay:   dayAccount,
	}, nil
}

func randomNumeric(length int) (string, error) {
	value := make([]byte, length)
	for i := range value {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("generate random numeric: %w", err)
		}
		value[i] = byte('0' + n.Int64())
	}
	return string(value), nil
}
