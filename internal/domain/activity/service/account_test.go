package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"bm-go/internal/domain/activity"
)

func TestAccountServiceQueryActivityAccount(t *testing.T) {
	repo := &fakeActivityRepository{
		account: activity.AccountEntity{
			UserID:            "xiaofuge",
			ActivityID:        100301,
			TotalCount:        100,
			TotalCountSurplus: 80,
			DayCount:          10,
			DayCountSurplus:   8,
			MonthCount:        30,
			MonthCountSurplus: 20,
		},
		accountExists: true,
		day: activity.AccountDayEntity{
			Day:             "2024-05-25",
			DayCount:        5,
			DayCountSurplus: 3,
		},
		dayExists: true,
		month: activity.AccountMonthEntity{
			Month:             "2024-05",
			MonthCount:        50,
			MonthCountSurplus: 35,
		},
		monthExists: true,
	}
	service := NewAccountService(repo)
	service.now = func() time.Time {
		return time.Date(2024, 5, 25, 10, 0, 0, 0, time.Local)
	}

	account, err := service.QueryActivityAccount(context.Background(), 100301, "xiaofuge")
	if err != nil {
		t.Fatalf("query activity account: %v", err)
	}

	if repo.day != (activity.AccountDayEntity{}) && repo.queriedDay != "2024-05-25" {
		t.Fatalf("expected queried day 2024-05-25, got %s", repo.queriedDay)
	}
	if repo.month != (activity.AccountMonthEntity{}) && repo.queriedMonth != "2024-05" {
		t.Fatalf("expected queried month 2024-05, got %s", repo.queriedMonth)
	}
	if account.TotalCount != 100 || account.TotalCountSurplus != 80 {
		t.Fatalf("expected total 100/80, got %d/%d", account.TotalCount, account.TotalCountSurplus)
	}
	if account.DayCount != 5 || account.DayCountSurplus != 3 {
		t.Fatalf("expected day 5/3, got %d/%d", account.DayCount, account.DayCountSurplus)
	}
	if account.MonthCount != 50 || account.MonthCountSurplus != 35 {
		t.Fatalf("expected month 50/35, got %d/%d", account.MonthCount, account.MonthCountSurplus)
	}
}

func TestAccountServiceQueryActivityAccountFallback(t *testing.T) {
	repo := &fakeActivityRepository{
		account: activity.AccountEntity{
			UserID:            "xiaofuge",
			ActivityID:        100301,
			TotalCount:        100,
			TotalCountSurplus: 80,
			DayCount:          10,
			DayCountSurplus:   8,
			MonthCount:        30,
			MonthCountSurplus: 20,
		},
		accountExists: true,
	}
	service := NewAccountService(repo)

	account, err := service.QueryActivityAccount(context.Background(), 100301, "xiaofuge")
	if err != nil {
		t.Fatalf("query activity account: %v", err)
	}

	if account.DayCount != 10 || account.DayCountSurplus != 10 {
		t.Fatalf("expected day fallback 10/10, got %d/%d", account.DayCount, account.DayCountSurplus)
	}
	if account.MonthCount != 30 || account.MonthCountSurplus != 30 {
		t.Fatalf("expected month fallback 30/30, got %d/%d", account.MonthCount, account.MonthCountSurplus)
	}
}

func TestAccountServiceQueryActivityAccountEmpty(t *testing.T) {
	repo := &fakeActivityRepository{}
	service := NewAccountService(repo)

	account, err := service.QueryActivityAccount(context.Background(), 100301, "xiaofuge")
	if err != nil {
		t.Fatalf("query activity account: %v", err)
	}

	if account.UserID != "xiaofuge" || account.ActivityID != 100301 {
		t.Fatalf("expected user and activity identity, got %s/%d", account.UserID, account.ActivityID)
	}
	if account.TotalCount != 0 || account.DayCount != 0 || account.MonthCount != 0 {
		t.Fatalf("expected empty account, got %+v", account)
	}
}

func TestAccountServiceQueryActivityAccountIllegalParam(t *testing.T) {
	service := NewAccountService(&fakeActivityRepository{})

	_, err := service.QueryActivityAccount(context.Background(), 0, "xiaofuge")
	if err == nil {
		t.Fatal("expected illegal param error")
	}

	_, err = service.QueryActivityAccount(context.Background(), 100301, " ")
	if err == nil {
		t.Fatal("expected illegal param error")
	}
}

func TestAccountServiceQueryActivityAccountRepositoryError(t *testing.T) {
	service := NewAccountService(&fakeActivityRepository{accountErr: errors.New("query account failed")})

	_, err := service.QueryActivityAccount(context.Background(), 100301, "xiaofuge")
	if err == nil {
		t.Fatal("expected account query error")
	}
}

func TestAccountServiceQueryActivityAccountDayError(t *testing.T) {
	repo := &fakeActivityRepository{
		account:       activity.AccountEntity{UserID: "xiaofuge", ActivityID: 100301},
		accountExists: true,
		dayErr:        errors.New("query day failed"),
	}
	service := NewAccountService(repo)

	_, err := service.QueryActivityAccount(context.Background(), 100301, "xiaofuge")
	if err == nil {
		t.Fatal("expected day query error")
	}
}

func TestAccountServiceQueryActivityAccountMonthError(t *testing.T) {
	repo := &fakeActivityRepository{
		account:       activity.AccountEntity{UserID: "xiaofuge", ActivityID: 100301},
		accountExists: true,
		monthErr:      errors.New("query month failed"),
	}
	service := NewAccountService(repo)

	_, err := service.QueryActivityAccount(context.Background(), 100301, "xiaofuge")
	if err == nil {
		t.Fatal("expected month query error")
	}
}

type fakeActivityRepository struct {
	account       activity.AccountEntity
	accountExists bool
	day           activity.AccountDayEntity
	dayExists     bool
	month         activity.AccountMonthEntity
	monthExists   bool
	queriedDay    string
	queriedMonth  string
	accountErr    error
	dayErr        error
	monthErr      error
}

func (f *fakeActivityRepository) QueryActivityAccount(ctx context.Context, activityID int64, userID string) (activity.AccountEntity, bool, error) {
	return f.account, f.accountExists, f.accountErr
}

func (f *fakeActivityRepository) QueryActivityAccountDay(ctx context.Context, activityID int64, userID string, day string) (activity.AccountDayEntity, bool, error) {
	f.queriedDay = day
	return f.day, f.dayExists, f.dayErr
}

func (f *fakeActivityRepository) QueryActivityAccountMonth(ctx context.Context, activityID int64, userID string, month string) (activity.AccountMonthEntity, bool, error) {
	f.queriedMonth = month
	return f.month, f.monthExists, f.monthErr
}
