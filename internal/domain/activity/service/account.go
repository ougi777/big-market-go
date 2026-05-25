package service

import (
	"context"
	"time"

	"bm-go/internal/domain/activity"
)

type AccountService struct {
	repo activity.AccountRepository
	now  func() time.Time
}

func NewAccountService(repo activity.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
		now:  time.Now,
	}
}

func (s *AccountService) QueryActivityAccount(ctx context.Context, activityID int64, userID string) (activity.AccountEntity, error) {
	account, exists, err := s.repo.QueryActivityAccount(ctx, activityID, userID)
	if err != nil {
		return activity.AccountEntity{}, err
	}
	if !exists {
		return activity.AccountEntity{
			UserID:     userID,
			ActivityID: activityID,
		}, nil
	}

	result := activity.AccountEntity{
		UserID:            account.UserID,
		ActivityID:        account.ActivityID,
		TotalCount:        account.TotalCount,
		TotalCountSurplus: account.TotalCountSurplus,
	}

	day := s.now().Format("2006-01-02")
	dayAccount, dayExists, err := s.repo.QueryActivityAccountDay(ctx, activityID, userID, day)
	if err != nil {
		return activity.AccountEntity{}, err
	}
	if dayExists {
		result.DayCount = dayAccount.DayCount
		result.DayCountSurplus = dayAccount.DayCountSurplus
	} else {
		result.DayCount = account.DayCount
		result.DayCountSurplus = account.DayCount
	}

	month := s.now().Format("2006-01")
	monthAccount, monthExists, err := s.repo.QueryActivityAccountMonth(ctx, activityID, userID, month)
	if err != nil {
		return activity.AccountEntity{}, err
	}
	if monthExists {
		result.MonthCount = monthAccount.MonthCount
		result.MonthCountSurplus = monthAccount.MonthCountSurplus
	} else {
		result.MonthCount = account.MonthCount
		result.MonthCountSurplus = account.MonthCount
	}

	return result, nil
}
