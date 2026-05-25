package chain

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"bm-go/internal/domain/strategy"
)

type WeightChain struct {
	BaseChain
	repo     strategy.Repository
	dispatch strategy.Dispatch
}

func NewWeightChain(repo strategy.Repository, dispatch strategy.Dispatch) *WeightChain {
	return &WeightChain{repo: repo, dispatch: dispatch}
}

func (c *WeightChain) Logic(ctx context.Context, userID string, strategyID int64) (AwardResult, error) {
	ruleValue, err := c.repo.QueryStrategyRuleValue(ctx, strategyID, RuleWeight)
	if err != nil {
		return AwardResult{}, err
	}

	analyticalValueGroup, err := parseWeightRule(ruleValue)
	if err != nil {
		return AwardResult{}, err
	}
	if len(analyticalValueGroup) == 0 {
		return c.Next(ctx, userID, strategyID)
	}

	userScore, err := c.repo.QueryActivityAccountTotalUseCount(ctx, userID, strategyID)
	if err != nil {
		return AwardResult{}, err
	}

	thresholds := make([]int, 0, len(analyticalValueGroup))
	for threshold := range analyticalValueGroup {
		thresholds = append(thresholds, threshold)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(thresholds)))

	for _, threshold := range thresholds {
		if userScore >= threshold {
			awardID, err := c.dispatch.GetRandomAwardID(ctx, strategyID, analyticalValueGroup[threshold])
			if err != nil {
				return AwardResult{}, err
			}
			return AwardResult{AwardID: awardID, LogicModel: RuleWeight}, nil
		}
	}
	return c.Next(ctx, userID, strategyID)
}

func parseWeightRule(ruleValue string) (map[int]string, error) {
	groups := strings.Fields(ruleValue)
	result := make(map[int]string, len(groups))
	for _, group := range groups {
		parts := strings.SplitN(group, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid weight rule: %s", group)
		}

		threshold, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}
		result[threshold] = group
	}
	return result, nil
}
