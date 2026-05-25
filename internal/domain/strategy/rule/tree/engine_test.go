package tree

import (
	"context"
	"strings"
	"testing"
)

func TestEngineProcessStockTakeOver(t *testing.T) {
	repo := &fakeTreeRepository{}
	dispatch := &fakeStockDispatch{ok: true}
	engine := NewEngine(map[string]Node{
		RuleStock: NewStockNode(repo, dispatch),
	}, RuleTree{
		TreeID:   "tree_stock",
		RootRule: RuleStock,
		NodeMap: map[string]RuleTreeNode{
			RuleStock: {
				RuleKey:   RuleStock,
				RuleValue: "stock",
			},
		},
	})

	award, err := engine.Process(context.Background(), "user001", 100001, 101)
	if err != nil {
		t.Fatalf("process tree: %v", err)
	}
	if award.AwardID != 101 {
		t.Fatalf("expected award 101, got %d", award.AwardID)
	}
	if repo.queuedAwardID != 101 {
		t.Fatalf("expected queued award 101, got %d", repo.queuedAwardID)
	}
}

func TestEngineProcessLuckAward(t *testing.T) {
	repo := &fakeTreeRepository{}
	engine := NewEngine(map[string]Node{
		RuleLuckAward: NewLuckAwardNode(repo),
	}, RuleTree{
		TreeID:   "tree_luck",
		RootRule: RuleLuckAward,
		NodeMap: map[string]RuleTreeNode{
			RuleLuckAward: {
				RuleKey:   RuleLuckAward,
				RuleValue: "102:0.01,1",
			},
		},
	})

	award, err := engine.Process(context.Background(), "user001", 100001, 101)
	if err != nil {
		t.Fatalf("process tree: %v", err)
	}
	if award.AwardID != 102 {
		t.Fatalf("expected luck award 102, got %d", award.AwardID)
	}
	if award.AwardRuleValue != "0.01,1" {
		t.Fatalf("expected rule value 0.01,1, got %s", award.AwardRuleValue)
	}
}

func TestEngineProcessRouteByCheckType(t *testing.T) {
	rootCalls := 0
	nextCalls := 0
	engine := NewEngine(map[string]Node{
		"root": &fakeLogicNode{
			action: Allow(),
			calls:  &rootCalls,
		},
		"next": &fakeLogicNode{
			action: TakeOver(&StrategyAward{AwardID: 202, AwardRuleValue: "next"}),
			calls:  &nextCalls,
		},
	}, RuleTree{
		TreeID:   "tree_route",
		RootRule: "root",
		NodeMap: map[string]RuleTreeNode{
			"root": {
				RuleKey: "root",
				Lines: []RuleTreeNodeLine{
					{
						RuleNodeTo:     "next",
						RuleLimitType:  LimitEqual,
						RuleLimitValue: "ALLOW",
					},
				},
			},
			"next": {
				RuleKey: "next",
			},
		},
	})

	award, err := engine.Process(context.Background(), "user001", 100001, 101)
	if err != nil {
		t.Fatalf("process tree: %v", err)
	}
	if award.AwardID != 202 {
		t.Fatalf("expected routed award 202, got %d", award.AwardID)
	}
	if rootCalls != 1 || nextCalls != 1 {
		t.Fatalf("expected root and next called once, got root=%d next=%d", rootCalls, nextCalls)
	}
}

func TestEngineProcessMissingTreeNode(t *testing.T) {
	engine := NewEngine(map[string]Node{}, RuleTree{
		TreeID:   "tree_missing",
		RootRule: "missing",
		NodeMap:  map[string]RuleTreeNode{},
	})

	_, err := engine.Process(context.Background(), "user001", 100001, 101)
	if err == nil || !strings.Contains(err.Error(), "rule tree node not found") {
		t.Fatalf("expected missing tree node error, got %v", err)
	}
}

func TestEngineProcessMissingLogicNode(t *testing.T) {
	engine := NewEngine(map[string]Node{}, RuleTree{
		TreeID:   "tree_missing_logic",
		RootRule: "root",
		NodeMap: map[string]RuleTreeNode{
			"root": {
				RuleKey: "missing_logic",
			},
		},
	})

	_, err := engine.Process(context.Background(), "user001", 100001, 101)
	if err == nil || !strings.Contains(err.Error(), "logic tree node not found") {
		t.Fatalf("expected missing logic node error, got %v", err)
	}
}

type fakeTreeRepository struct {
	todayCount    int
	queuedAwardID int
}

func (f *fakeTreeRepository) QueryTodayUserRaffleCount(ctx context.Context, userID string, strategyID int64) (int, error) {
	return f.todayCount, nil
}

func (f *fakeTreeRepository) AwardStockConsumeSendQueue(ctx context.Context, strategyID int64, awardID int) error {
	f.queuedAwardID = awardID
	return nil
}

type fakeStockDispatch struct {
	ok bool
}

func (f *fakeStockDispatch) SubtractionAwardStock(ctx context.Context, strategyID int64, awardID int) (bool, error) {
	return f.ok, nil
}

type fakeLogicNode struct {
	action TreeAction
	calls  *int
}

func (f *fakeLogicNode) Logic(ctx context.Context, userID string, strategyID int64, awardID int, ruleValue string) (TreeAction, error) {
	if f.calls != nil {
		*f.calls++
	}
	return f.action, nil
}
