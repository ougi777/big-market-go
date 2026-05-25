package service

import (
	"context"

	"bm-go/internal/domain/strategy"
	"bm-go/internal/domain/strategy/rule/chain"
	"bm-go/internal/domain/strategy/rule/tree"
)

type RaffleService struct {
	chainFactory *chain.Factory
	repo         strategy.RaffleRepository
	treeNodes    map[string]tree.Node
}

func NewRaffleService(chainFactory *chain.Factory, repo strategy.RaffleRepository, treeNodes map[string]tree.Node) *RaffleService {
	return &RaffleService{chainFactory: chainFactory, repo: repo, treeNodes: treeNodes}
}

func (s *RaffleService) PerformRaffle(ctx context.Context, userID string, strategyID int64) (chain.AwardResult, error) {
	logicChain, err := s.chainFactory.OpenLogicChain(ctx, strategyID)
	if err != nil {
		return chain.AwardResult{}, err
	}
	chainResult, err := logicChain.Logic(ctx, userID, strategyID)
	if err != nil {
		return chain.AwardResult{}, err
	}
	if chainResult.LogicModel != chain.RuleDefault {
		return s.buildAwardResult(ctx, strategyID, chainResult.AwardID, chainResult.LogicModel, chainResult.AwardRuleValue)
	}

	treeAward, err := s.raffleLogicTree(ctx, userID, strategyID, chainResult.AwardID)
	if err != nil {
		return chain.AwardResult{}, err
	}
	return s.buildAwardResult(ctx, strategyID, treeAward.AwardID, chainResult.LogicModel, treeAward.AwardRuleValue)
}

func (s *RaffleService) raffleLogicTree(ctx context.Context, userID string, strategyID int64, awardID int) (tree.StrategyAward, error) {
	ruleModels, err := s.repo.QueryStrategyAwardRuleModels(ctx, strategyID, awardID)
	if err != nil {
		return tree.StrategyAward{}, err
	}
	if ruleModels == "" {
		return tree.StrategyAward{AwardID: awardID}, nil
	}

	ruleTree, ok, err := s.repo.QueryRuleTreeByTreeID(ctx, ruleModels)
	if err != nil {
		return tree.StrategyAward{}, err
	}
	if !ok {
		return tree.StrategyAward{}, nil
	}

	engine := tree.NewEngine(s.treeNodes, ruleTree)
	treeAward, err := engine.Process(ctx, userID, strategyID, awardID)
	if err != nil {
		return tree.StrategyAward{}, err
	}
	if treeAward.AwardID == 0 {
		return tree.StrategyAward{AwardID: awardID}, nil
	}
	return treeAward, nil
}

func (s *RaffleService) buildAwardResult(ctx context.Context, strategyID int64, awardID int, logicModel string, awardRuleValue string) (chain.AwardResult, error) {
	award, ok, err := s.repo.QueryStrategyAwardEntity(ctx, strategyID, awardID)
	if err != nil {
		return chain.AwardResult{}, err
	}
	if !ok {
		return chain.AwardResult{AwardID: awardID, LogicModel: logicModel, AwardRuleValue: awardRuleValue}, nil
	}
	return chain.AwardResult{
		AwardID:        awardID,
		AwardTitle:     award.AwardTitle,
		AwardIndex:     award.Sort,
		LogicModel:     logicModel,
		AwardRuleValue: awardRuleValue,
	}, nil
}
