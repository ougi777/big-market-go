package chain

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"bm-go/internal/domain/strategy"
)

type BlacklistChain struct {
	BaseChain
	repo strategy.Repository
}

func NewBlacklistChain(repo strategy.Repository) *BlacklistChain {
	return &BlacklistChain{repo: repo}
}

func (c *BlacklistChain) Logic(ctx context.Context, userID string, strategyID int64) (AwardResult, error) {
	ruleValue, err := c.repo.QueryStrategyRuleValue(ctx, strategyID, RuleBlacklist)
	if err != nil {
		return AwardResult{}, err
	}

	awardID, users, err := parseBlacklistRule(ruleValue)
	if err != nil {
		return AwardResult{}, err
	}

	for _, blackUserID := range users {
		if userID == blackUserID {
			return AwardResult{
				AwardID:        awardID,
				LogicModel:     RuleBlacklist,
				AwardRuleValue: "0.01,1",
			}, nil
		}
	}
	return c.Next(ctx, userID, strategyID)
}

func parseBlacklistRule(ruleValue string) (int, []string, error) {
	parts := strings.SplitN(ruleValue, ":", 2)
	if len(parts) != 2 {
		return 0, nil, fmt.Errorf("invalid blacklist rule: %s", ruleValue)
	}

	awardID, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, nil, err
	}

	users := strings.Split(parts[1], ",")
	for i := range users {
		users[i] = strings.TrimSpace(users[i])
	}
	return awardID, users, nil
}
