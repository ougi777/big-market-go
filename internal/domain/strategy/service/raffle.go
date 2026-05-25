package service

import (
	"context"

	"bm-go/internal/domain/strategy/rule/chain"
)

type RaffleService struct {
	chainFactory *chain.Factory
}

func NewRaffleService(chainFactory *chain.Factory) *RaffleService {
	return &RaffleService{chainFactory: chainFactory}
}

func (s *RaffleService) PerformRaffle(ctx context.Context, userID string, strategyID int64) (chain.AwardResult, error) {
	logicChain, err := s.chainFactory.OpenLogicChain(ctx, strategyID)
	if err != nil {
		return chain.AwardResult{}, err
	}
	return logicChain.Logic(ctx, userID, strategyID)
}
