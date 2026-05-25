package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"
)

func TestPartakeServiceCreateOrder(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	repo := &fakePartakeRepository{
		activity: activity.ActivityEntity{
			ActivityID:    100301,
			ActivityName:  "测试活动",
			BeginDateTime: now.Add(-time.Hour),
			EndDateTime:   now.Add(time.Hour),
			StrategyID:    100006,
			State:         activity.ActivityStateOpen,
		},
		activityExists: true,
		account: activity.AccountEntity{
			UserID:            "xiaofuge",
			ActivityID:        100301,
			TotalCount:        10,
			TotalCountSurplus: 9,
			DayCount:          3,
			MonthCount:        8,
		},
		accountExists: true,
	}
	service := NewPartakeService(repo)
	service.now = func() time.Time { return now }
	service.orderIDGenerator = func() (string, error) { return "123456789012", nil }

	order, err := service.CreateOrder(context.Background(), "xiaofuge", 100301)
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if order.OrderID != "123456789012" {
		t.Fatalf("expected order id, got %s", order.OrderID)
	}
	if order.StrategyID != 100006 || order.OrderState != activity.UserRaffleOrderCreate {
		t.Fatalf("expected strategy and state, got %+v", order)
	}
	if !repo.saved {
		t.Fatalf("expected aggregate saved")
	}
	if repo.aggregate.ExistAccountMonth || repo.aggregate.ExistAccountDay {
		t.Fatalf("expected new month and day accounts")
	}
	if repo.aggregate.ActivityAccountMonth.Month != "2026-05" {
		t.Fatalf("expected month 2026-05, got %s", repo.aggregate.ActivityAccountMonth.Month)
	}
	if repo.aggregate.ActivityAccountDay.Day != "2026-05-25" {
		t.Fatalf("expected day 2026-05-25, got %s", repo.aggregate.ActivityAccountDay.Day)
	}
}

func TestPartakeServiceCreateOrderReturnExisting(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	repo := &fakePartakeRepository{
		activity: activity.ActivityEntity{
			ActivityID:    100301,
			ActivityName:  "测试活动",
			BeginDateTime: now.Add(-time.Hour),
			EndDateTime:   now.Add(time.Hour),
			StrategyID:    100006,
			State:         activity.ActivityStateOpen,
		},
		activityExists: true,
		existingOrder: activity.UserRaffleOrderEntity{
			UserID:     "xiaofuge",
			ActivityID: 100301,
			OrderID:    "999999999999",
			StrategyID: 100006,
		},
		existingOrderExists: true,
	}
	service := NewPartakeService(repo)
	service.now = func() time.Time { return now }

	order, err := service.CreateOrder(context.Background(), "xiaofuge", 100301)
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if order.OrderID != "999999999999" {
		t.Fatalf("expected existing order, got %s", order.OrderID)
	}
	if repo.saved {
		t.Fatalf("did not expect aggregate saved")
	}
}

func TestPartakeServiceCreateOrderActivityStateError(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	repo := &fakePartakeRepository{
		activity: activity.ActivityEntity{
			BeginDateTime: now.Add(-time.Hour),
			EndDateTime:   now.Add(time.Hour),
			State:         "close",
		},
		activityExists: true,
	}
	service := NewPartakeService(repo)
	service.now = func() time.Time { return now }

	_, err := service.CreateOrder(context.Background(), "xiaofuge", 100301)
	assertAppErrorCode(t, err, types.ResponseCodeActivityStateError)
}

func TestPartakeServiceCreateOrderActivityDateError(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	repo := &fakePartakeRepository{
		activity: activity.ActivityEntity{
			BeginDateTime: now.Add(time.Hour),
			EndDateTime:   now.Add(2 * time.Hour),
			State:         activity.ActivityStateOpen,
		},
		activityExists: true,
	}
	service := NewPartakeService(repo)
	service.now = func() time.Time { return now }

	_, err := service.CreateOrder(context.Background(), "xiaofuge", 100301)
	assertAppErrorCode(t, err, types.ResponseCodeActivityDateError)
}

func TestPartakeServiceCreateOrderAccountQuotaError(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	repo := &fakePartakeRepository{
		activity: activity.ActivityEntity{
			BeginDateTime: now.Add(-time.Hour),
			EndDateTime:   now.Add(time.Hour),
			State:         activity.ActivityStateOpen,
		},
		activityExists: true,
		account: activity.AccountEntity{
			TotalCount:        10,
			TotalCountSurplus: 0,
		},
		accountExists: true,
	}
	service := NewPartakeService(repo)
	service.now = func() time.Time { return now }

	_, err := service.CreateOrder(context.Background(), "xiaofuge", 100301)
	assertAppErrorCode(t, err, types.ResponseCodeAccountQuotaError)
}

func TestPartakeServiceCreateOrderDayQuotaError(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	repo := &fakePartakeRepository{
		activity: activity.ActivityEntity{
			BeginDateTime: now.Add(-time.Hour),
			EndDateTime:   now.Add(time.Hour),
			State:         activity.ActivityStateOpen,
		},
		activityExists: true,
		account: activity.AccountEntity{
			TotalCountSurplus: 9,
			DayCount:          3,
			MonthCount:        8,
		},
		accountExists: true,
		day: activity.AccountDayEntity{
			DayCount:        3,
			DayCountSurplus: 0,
		},
		dayExists: true,
	}
	service := NewPartakeService(repo)
	service.now = func() time.Time { return now }

	_, err := service.CreateOrder(context.Background(), "xiaofuge", 100301)
	assertAppErrorCode(t, err, types.ResponseCodeAccountDayQuotaError)
}

func TestPartakeServiceCreateOrderMonthQuotaError(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	repo := &fakePartakeRepository{
		activity: activity.ActivityEntity{
			BeginDateTime: now.Add(-time.Hour),
			EndDateTime:   now.Add(time.Hour),
			State:         activity.ActivityStateOpen,
		},
		activityExists: true,
		account: activity.AccountEntity{
			TotalCountSurplus: 9,
			DayCount:          3,
			MonthCount:        8,
		},
		accountExists: true,
		month: activity.AccountMonthEntity{
			MonthCount:        8,
			MonthCountSurplus: 0,
		},
		monthExists: true,
	}
	service := NewPartakeService(repo)
	service.now = func() time.Time { return now }

	_, err := service.CreateOrder(context.Background(), "xiaofuge", 100301)
	assertAppErrorCode(t, err, types.ResponseCodeAccountMonthQuotaError)
}

func TestPartakeServiceCreateOrderUsesExistingMonthAndDayAccounts(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	repo := &fakePartakeRepository{
		activity: activity.ActivityEntity{
			ActivityID:    100301,
			ActivityName:  "测试活动",
			BeginDateTime: now.Add(-time.Hour),
			EndDateTime:   now.Add(time.Hour),
			StrategyID:    100006,
			State:         activity.ActivityStateOpen,
		},
		activityExists: true,
		account: activity.AccountEntity{
			TotalCountSurplus: 9,
			DayCount:          3,
			MonthCount:        8,
		},
		accountExists: true,
		month: activity.AccountMonthEntity{
			Month:             "2026-05",
			MonthCount:        8,
			MonthCountSurplus: 6,
		},
		monthExists: true,
		day: activity.AccountDayEntity{
			Day:             "2026-05-25",
			DayCount:        3,
			DayCountSurplus: 2,
		},
		dayExists: true,
	}
	service := NewPartakeService(repo)
	service.now = func() time.Time { return now }
	service.orderIDGenerator = func() (string, error) { return "123456789012", nil }

	_, err := service.CreateOrder(context.Background(), "xiaofuge", 100301)
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if !repo.aggregate.ExistAccountMonth || !repo.aggregate.ExistAccountDay {
		t.Fatalf("expected existing month and day accounts")
	}
	if repo.aggregate.ActivityAccountMonth.MonthCountSurplus != 6 {
		t.Fatalf("expected existing month surplus 6, got %d", repo.aggregate.ActivityAccountMonth.MonthCountSurplus)
	}
	if repo.aggregate.ActivityAccountDay.DayCountSurplus != 2 {
		t.Fatalf("expected existing day surplus 2, got %d", repo.aggregate.ActivityAccountDay.DayCountSurplus)
	}
}

type fakePartakeRepository struct {
	activity            activity.ActivityEntity
	activityExists      bool
	existingOrder       activity.UserRaffleOrderEntity
	existingOrderExists bool
	account             activity.AccountEntity
	accountExists       bool
	day                 activity.AccountDayEntity
	dayExists           bool
	month               activity.AccountMonthEntity
	monthExists         bool
	saved               bool
	aggregate           activity.CreatePartakeOrderAggregate
}

func (f *fakePartakeRepository) QueryActivityByActivityID(ctx context.Context, activityID int64) (activity.ActivityEntity, bool, error) {
	return f.activity, f.activityExists, nil
}

func (f *fakePartakeRepository) QueryNoUsedRaffleOrder(ctx context.Context, userID string, activityID int64) (activity.UserRaffleOrderEntity, bool, error) {
	return f.existingOrder, f.existingOrderExists, nil
}

func (f *fakePartakeRepository) QueryActivityAccount(ctx context.Context, activityID int64, userID string) (activity.AccountEntity, bool, error) {
	return f.account, f.accountExists, nil
}

func (f *fakePartakeRepository) QueryActivityAccountDay(ctx context.Context, activityID int64, userID string, day string) (activity.AccountDayEntity, bool, error) {
	return f.day, f.dayExists, nil
}

func (f *fakePartakeRepository) QueryActivityAccountMonth(ctx context.Context, activityID int64, userID string, month string) (activity.AccountMonthEntity, bool, error) {
	return f.month, f.monthExists, nil
}

func (f *fakePartakeRepository) SaveCreatePartakeOrder(ctx context.Context, aggregate activity.CreatePartakeOrderAggregate) error {
	f.saved = true
	f.aggregate = aggregate
	return nil
}

func assertAppErrorCode(t *testing.T, err error, code types.ResponseCode) {
	t.Helper()

	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != code {
		t.Fatalf("expected code %s, got %s", code.Code, appErr.Code.Code)
	}
}
