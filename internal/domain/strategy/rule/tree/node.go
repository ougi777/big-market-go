package tree

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type Node interface {
	Logic(ctx context.Context, userID string, strategyID int64, awardID int, ruleValue string) (TreeAction, error)
}

type Repository interface {
	QueryTodayUserRaffleCount(ctx context.Context, userID string, strategyID int64) (int, error)
	AwardStockConsumeSendQueue(ctx context.Context, strategyID int64, awardID int) error
}

type StockDispatch interface {
	SubtractionAwardStock(ctx context.Context, strategyID int64, awardID int) (bool, error)
}

type LockNode struct {
	repo Repository
}

func NewLockNode(repo Repository) *LockNode {
	return &LockNode{repo: repo}
}

func (n *LockNode) Logic(ctx context.Context, userID string, strategyID int64, awardID int, ruleValue string) (TreeAction, error) {
	raffleCount, err := strconv.Atoi(strings.TrimSpace(ruleValue))
	if err != nil {
		return TreeAction{}, fmt.Errorf("invalid rule_lock value %q: %w", ruleValue, err)
	}

	userRaffleCount, err := n.repo.QueryTodayUserRaffleCount(ctx, userID, strategyID)
	if err != nil {
		return TreeAction{}, err
	}
	if userRaffleCount >= raffleCount {
		return Allow(), nil
	}
	return TakeOver(nil), nil
}

type StockNode struct {
	repo     Repository
	dispatch StockDispatch
}

func NewStockNode(repo Repository, dispatch StockDispatch) *StockNode {
	return &StockNode{repo: repo, dispatch: dispatch}
}

func (n *StockNode) Logic(ctx context.Context, userID string, strategyID int64, awardID int, ruleValue string) (TreeAction, error) {
	ok, err := n.dispatch.SubtractionAwardStock(ctx, strategyID, awardID)
	if err != nil {
		return TreeAction{}, err
	}
	if !ok {
		return Allow(), nil
	}

	if err := n.repo.AwardStockConsumeSendQueue(ctx, strategyID, awardID); err != nil {
		return TreeAction{}, err
	}
	return TakeOver(&StrategyAward{AwardID: awardID, AwardRuleValue: ruleValue}), nil
}

type LuckAwardNode struct {
	repo Repository
}

func NewLuckAwardNode(repo Repository) *LuckAwardNode {
	return &LuckAwardNode{repo: repo}
}

func (n *LuckAwardNode) Logic(ctx context.Context, userID string, strategyID int64, awardID int, ruleValue string) (TreeAction, error) {
	parts := strings.SplitN(ruleValue, ":", 2)
	if len(parts) == 0 || strings.TrimSpace(parts[0]) == "" {
		return TreeAction{}, fmt.Errorf("invalid rule_luck_award value: %s", ruleValue)
	}

	luckAwardID, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return TreeAction{}, err
	}
	awardRuleValue := ""
	if len(parts) > 1 {
		awardRuleValue = parts[1]
	}

	if err := n.repo.AwardStockConsumeSendQueue(ctx, strategyID, luckAwardID); err != nil {
		return TreeAction{}, err
	}
	return TakeOver(&StrategyAward{AwardID: luckAwardID, AwardRuleValue: awardRuleValue}), nil
}
