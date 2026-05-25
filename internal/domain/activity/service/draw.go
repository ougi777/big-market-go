package service

import (
	"context"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/domain/award"
	"bm-go/internal/domain/strategy/rule/chain"
)

type drawPartakeService interface {
	CreateOrder(ctx context.Context, userID string, activityID int64) (activity.UserRaffleOrderEntity, error)
}

type drawRaffleService interface {
	PerformRaffle(ctx context.Context, userID string, strategyID int64) (chain.AwardResult, error)
}

type drawAwardService interface {
	SaveUserAwardRecord(ctx context.Context, record award.UserAwardRecordEntity) error
}

type DrawService struct {
	partakeService drawPartakeService
	raffleService  drawRaffleService
	awardService   drawAwardService
	now            func() time.Time
}

func NewDrawService(partakeService drawPartakeService, raffleService drawRaffleService, awardService drawAwardService) *DrawService {
	return &DrawService{
		partakeService: partakeService,
		raffleService:  raffleService,
		awardService:   awardService,
		now:            time.Now,
	}
}

func (s *DrawService) Draw(ctx context.Context, userID string, activityID int64) (activity.DrawResult, error) {
	order, err := s.partakeService.CreateOrder(ctx, userID, activityID)
	if err != nil {
		return activity.DrawResult{}, err
	}

	awardResult, err := s.raffleService.PerformRaffle(ctx, order.UserID, order.StrategyID)
	if err != nil {
		return activity.DrawResult{}, err
	}

	if err := s.awardService.SaveUserAwardRecord(ctx, award.UserAwardRecordEntity{
		UserID:      order.UserID,
		ActivityID:  order.ActivityID,
		StrategyID:  order.StrategyID,
		OrderID:     order.OrderID,
		AwardID:     awardResult.AwardID,
		AwardTitle:  awardResult.AwardTitle,
		AwardTime:   s.now(),
		AwardState:  award.AwardStateCreate,
		AwardConfig: awardResult.AwardRuleValue,
	}); err != nil {
		return activity.DrawResult{}, err
	}

	return activity.DrawResult{
		AwardID:    awardResult.AwardID,
		AwardTitle: awardResult.AwardTitle,
		AwardIndex: awardResult.AwardIndex,
	}, nil
}
