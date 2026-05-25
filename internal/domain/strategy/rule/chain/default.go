package chain

import (
	"context"

	"bm-go/internal/domain/strategy"
)

type DefaultChain struct {
	dispatch strategy.Dispatch
}

func NewDefaultChain(dispatch strategy.Dispatch) *DefaultChain {
	return &DefaultChain{dispatch: dispatch}
}

func (c *DefaultChain) SetNext(next LogicChain) {}

func (c *DefaultChain) Logic(ctx context.Context, userID string, strategyID int64) (AwardResult, error) {
	awardID, err := c.dispatch.GetRandomAwardID(ctx, strategyID)
	if err != nil {
		return AwardResult{}, err
	}
	return AwardResult{AwardID: awardID, LogicModel: RuleDefault}, nil
}
