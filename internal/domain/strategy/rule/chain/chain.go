package chain

import (
	"context"
	"errors"
)

const (
	RuleDefault   = "rule_default"
	RuleBlacklist = "rule_blacklist"
	RuleWeight    = "rule_weight"
)

var ErrNextChainMissing = errors.New("next logic chain is missing")

type AwardResult struct {
	AwardID        int
	AwardTitle     string
	AwardIndex     int
	LogicModel     string
	AwardRuleValue string
}

type LogicChain interface {
	Logic(ctx context.Context, userID string, strategyID int64) (AwardResult, error)
	SetNext(next LogicChain)
}

type BaseChain struct {
	next LogicChain
}

func (b *BaseChain) SetNext(next LogicChain) {
	b.next = next
}

func (b *BaseChain) Next(ctx context.Context, userID string, strategyID int64) (AwardResult, error) {
	if b.next == nil {
		return AwardResult{}, ErrNextChainMissing
	}
	return b.next.Logic(ctx, userID, strategyID)
}
