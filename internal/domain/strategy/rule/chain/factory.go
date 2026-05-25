package chain

import (
	"context"
	"fmt"
	"sync"

	"bm-go/internal/domain/strategy"
)

type Factory struct {
	repo     strategy.Repository
	dispatch strategy.Dispatch
	registry map[string]func() LogicChain
	cache    sync.Map
}

func NewFactory(repo strategy.Repository, dispatch strategy.Dispatch) *Factory {
	factory := &Factory{
		repo:     repo,
		dispatch: dispatch,
		registry: make(map[string]func() LogicChain),
	}
	factory.Register(RuleBlacklist, func() LogicChain {
		return NewBlacklistChain(repo)
	})
	factory.Register(RuleWeight, func() LogicChain {
		return NewWeightChain(repo, dispatch)
	})
	factory.Register(RuleDefault, func() LogicChain {
		return NewDefaultChain(dispatch)
	})
	return factory
}

func (f *Factory) Register(ruleModel string, constructor func() LogicChain) {
	f.registry[ruleModel] = constructor
}

func (f *Factory) OpenLogicChain(ctx context.Context, strategyID int64) (LogicChain, error) {
	if logicChain, ok := f.cache.Load(strategyID); ok {
		return logicChain.(LogicChain), nil
	}

	strategyEntity, err := f.repo.QueryStrategyEntityByStrategyID(ctx, strategyID)
	if err != nil {
		return nil, err
	}

	ruleModels := strategyEntity.RuleModels()
	if len(ruleModels) == 0 {
		ruleModels = []string{RuleDefault}
	} else {
		ruleModels = append(ruleModels, RuleDefault)
	}

	logicChain, err := f.build(ruleModels)
	if err != nil {
		return nil, err
	}
	f.cache.Store(strategyID, logicChain)
	return logicChain, nil
}

func (f *Factory) Evict(strategyID int64) {
	f.cache.Delete(strategyID)
}

func (f *Factory) build(ruleModels []string) (LogicChain, error) {
	var head LogicChain
	var current LogicChain

	for _, ruleModel := range ruleModels {
		constructor, ok := f.registry[ruleModel]
		if !ok {
			return nil, fmt.Errorf("unknown logic rule model: %s", ruleModel)
		}

		node := constructor()
		if head == nil {
			head = node
			current = node
			continue
		}
		current.SetNext(node)
		current = node
	}
	return head, nil
}
