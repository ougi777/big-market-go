package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/types"
)

type ArmoryService struct {
	repo  strategy.ArmoryRepository
	store strategy.RateTableStore
	rand  *rand.Rand
}

func NewArmoryService(repo strategy.ArmoryRepository, store strategy.RateTableStore) *ArmoryService {
	return &ArmoryService{
		repo:  repo,
		store: store,
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *ArmoryService) AssembleLotteryStrategy(ctx context.Context, strategyID int64) error {
	awards, err := s.repo.QueryStrategyAwardList(ctx, strategyID)
	if err != nil {
		return err
	}
	if len(awards) == 0 {
		return fmt.Errorf("strategy awards are empty: %d", strategyID)
	}

	for _, award := range awards {
		cacheKey := types.RedisKeyStrategyAwardCount + strconv.FormatInt(strategyID, 10) + types.Underline + strconv.Itoa(award.AwardID)
		if err := s.store.CacheStrategyAwardCount(ctx, cacheKey, award.AwardCountSurplus); err != nil {
			return err
		}
	}

	if err := s.assembleRateTable(ctx, strconv.FormatInt(strategyID, 10), awards); err != nil {
		return err
	}

	strategyEntity, err := s.repo.QueryStrategyEntityByStrategyID(ctx, strategyID)
	if err != nil {
		return err
	}
	ruleWeight := strategyEntity.RuleWeight()
	if ruleWeight == "" {
		return nil
	}

	ruleEntity, ok, err := s.repo.QueryStrategyRule(ctx, strategyID, ruleWeight)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("strategy rule weight is empty: %d", strategyID)
	}

	for weightKey, awardIDs := range ruleEntity.RuleWeightValues() {
		filtered := filterAwards(awards, awardIDs)
		if len(filtered) == 0 {
			continue
		}

		key := strconv.FormatInt(strategyID, 10) + types.Underline + weightKey
		if err := s.assembleRateTable(ctx, key, filtered); err != nil {
			return err
		}
	}
	return nil
}

func (s *ArmoryService) assembleRateTable(ctx context.Context, key string, awards []strategy.StrategyAwardEntity) error {
	minRate := minAwardRate(awards)
	rateRange := convertRateRange(minRate)

	rateTableValues := make([]int, 0, rateRange)
	for _, award := range awards {
		count := int(float64(rateRange) * award.AwardRate)
		for i := 0; i < count; i++ {
			rateTableValues = append(rateTableValues, award.AwardID)
		}
	}
	if len(rateTableValues) == 0 {
		return fmt.Errorf("strategy rate table is empty: %s", key)
	}

	s.rand.Shuffle(len(rateTableValues), func(i, j int) {
		rateTableValues[i], rateTableValues[j] = rateTableValues[j], rateTableValues[i]
	})

	rateTable := make(map[int]int, len(rateTableValues))
	for i, awardID := range rateTableValues {
		rateTable[i] = awardID
	}
	return s.store.StoreStrategyAwardSearchRateTable(ctx, key, len(rateTable), rateTable)
}

func minAwardRate(awards []strategy.StrategyAwardEntity) float64 {
	minRate := 0.0
	for _, award := range awards {
		if award.AwardRate <= 0 {
			continue
		}
		if minRate == 0 || award.AwardRate < minRate {
			minRate = award.AwardRate
		}
	}
	return minRate
}

func convertRateRange(minRate float64) int {
	if minRate <= 0 {
		return 1
	}

	current := minRate
	max := 1
	for current < 1 {
		current *= 10
		max *= 10
	}
	return max
}

func filterAwards(awards []strategy.StrategyAwardEntity, awardIDs []int) []strategy.StrategyAwardEntity {
	awardIDSet := make(map[int]struct{}, len(awardIDs))
	for _, awardID := range awardIDs {
		awardIDSet[awardID] = struct{}{}
	}

	filtered := make([]strategy.StrategyAwardEntity, 0, len(awards))
	for _, award := range awards {
		if _, ok := awardIDSet[award.AwardID]; ok {
			filtered = append(filtered, award)
		}
	}
	return filtered
}
