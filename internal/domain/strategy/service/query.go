package service

import (
	"context"

	"bm-go/internal/domain/strategy"
)

type QueryService struct {
	repo strategy.QueryRepository
}

func NewQueryService(repo strategy.QueryRepository) *QueryService {
	return &QueryService{repo: repo}
}

func (s *QueryService) QueryRaffleAwardList(ctx context.Context, activityID int64, userID string) ([]RaffleAward, error) {
	strategyID, err := s.repo.QueryStrategyIDByActivityID(ctx, activityID)
	if err != nil {
		return nil, err
	}

	awards, err := s.repo.QueryStrategyAwardList(ctx, strategyID)
	if err != nil {
		return nil, err
	}

	treeIDs := make([]string, 0, len(awards))
	for _, award := range awards {
		if award.RuleModels != "" {
			treeIDs = append(treeIDs, award.RuleModels)
		}
	}

	lockCounts, err := s.repo.QueryAwardRuleLockCount(ctx, treeIDs)
	if err != nil {
		return nil, err
	}

	dayPartakeCount, err := s.repo.QueryRaffleActivityAccountDayPartakeCount(ctx, activityID, userID)
	if err != nil {
		return nil, err
	}

	result := make([]RaffleAward, 0, len(awards))
	for _, award := range awards {
		lockCount, hasLock := lockCounts[award.RuleModels]
		waitUnlockCount := 0
		unlocked := true
		if hasLock {
			unlocked = dayPartakeCount >= lockCount
			if lockCount > dayPartakeCount {
				waitUnlockCount = lockCount - dayPartakeCount
			}
		}

		result = append(result, RaffleAward{
			AwardID:            award.AwardID,
			AwardTitle:         award.AwardTitle,
			AwardSubtitle:      award.AwardSubtitle,
			Sort:               award.Sort,
			AwardRuleLockCount: lockCount,
			HasAwardRuleLock:   hasLock,
			IsAwardUnlock:      unlocked,
			WaitUnlockCount:    waitUnlockCount,
		})
	}
	return result, nil
}

func (s *QueryService) QueryRaffleStrategyRuleWeight(ctx context.Context, activityID int64, userID string) ([]RaffleStrategyRuleWeight, error) {
	strategyID, err := s.repo.QueryStrategyIDByActivityID(ctx, activityID)
	if err != nil {
		return nil, err
	}

	totalUseCount, err := s.repo.QueryRaffleActivityAccountPartakeCount(ctx, activityID, userID)
	if err != nil {
		return nil, err
	}

	ruleWeights, err := s.repo.QueryAwardRuleWeight(ctx, strategyID)
	if err != nil {
		return nil, err
	}

	result := make([]RaffleStrategyRuleWeight, 0, len(ruleWeights))
	for _, ruleWeight := range ruleWeights {
		awards := make([]RuleWeightAward, 0, len(ruleWeight.AwardList))
		for _, award := range ruleWeight.AwardList {
			awards = append(awards, RuleWeightAward{
				AwardID:    award.AwardID,
				AwardTitle: award.AwardTitle,
			})
		}

		result = append(result, RaffleStrategyRuleWeight{
			RuleWeightCount:                  ruleWeight.Weight,
			UserActivityAccountTotalUseCount: totalUseCount,
			StrategyAwards:                   awards,
		})
	}
	return result, nil
}

type RaffleAward struct {
	AwardID            int
	AwardTitle         string
	AwardSubtitle      string
	Sort               int
	AwardRuleLockCount int
	HasAwardRuleLock   bool
	IsAwardUnlock      bool
	WaitUnlockCount    int
}

type RaffleStrategyRuleWeight struct {
	RuleWeightCount                  int
	UserActivityAccountTotalUseCount int
	StrategyAwards                   []RuleWeightAward
}

type RuleWeightAward struct {
	AwardID    int
	AwardTitle string
}
