package service

import (
	"context"
	"testing"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/award"
	"bm-go/internal/domain/strategy/rule/chain"
)

func TestDrawServiceDraw(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)
	partake := &fakeDrawPartakeService{
		order: activity.UserRaffleOrderEntity{
			UserID:     "xiaofuge",
			ActivityID: 100301,
			StrategyID: 100006,
			OrderID:    "123456789012",
		},
	}
	raffle := &fakeDrawRaffleService{
		result: chain.AwardResult{
			AwardID:    101,
			AwardTitle: "积分",
			AwardIndex: 1,
		},
	}
	awardService := &fakeDrawAwardService{}
	service := NewDrawService(partake, raffle, awardService)
	service.now = func() time.Time { return now }

	result, err := service.Draw(context.Background(), "xiaofuge", 100301)
	if err != nil {
		t.Fatalf("draw: %v", err)
	}

	if result.AwardID != 101 || result.AwardTitle != "积分" || result.AwardIndex != 1 {
		t.Fatalf("expected draw result, got %+v", result)
	}
	if raffle.userID != "xiaofuge" || raffle.strategyID != 100006 {
		t.Fatalf("expected raffle factor, got %s/%d", raffle.userID, raffle.strategyID)
	}
	if awardService.record.OrderID != "123456789012" || awardService.record.AwardState != award.AwardStateCreate {
		t.Fatalf("expected award record, got %+v", awardService.record)
	}
	if !awardService.record.AwardTime.Equal(now) {
		t.Fatalf("expected award time %s, got %s", now, awardService.record.AwardTime)
	}
}

type fakeDrawPartakeService struct {
	order activity.UserRaffleOrderEntity
}

func (f *fakeDrawPartakeService) CreateOrder(ctx context.Context, userID string, activityID int64) (activity.UserRaffleOrderEntity, error) {
	return f.order, nil
}

type fakeDrawRaffleService struct {
	userID     string
	strategyID int64
	result     chain.AwardResult
}

func (f *fakeDrawRaffleService) PerformRaffle(ctx context.Context, userID string, strategyID int64) (chain.AwardResult, error) {
	f.userID = userID
	f.strategyID = strategyID
	return f.result, nil
}

type fakeDrawAwardService struct {
	record award.UserAwardRecordEntity
}

func (f *fakeDrawAwardService) SaveUserAwardRecord(ctx context.Context, record award.UserAwardRecordEntity) error {
	f.record = record
	return nil
}
