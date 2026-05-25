package tree

import (
	"context"
	"fmt"
)

type Engine struct {
	nodes map[string]Node
	tree  RuleTree
}

func NewEngine(nodes map[string]Node, tree RuleTree) *Engine {
	return &Engine{nodes: nodes, tree: tree}
}

func (e *Engine) Process(ctx context.Context, userID string, strategyID int64, awardID int) (StrategyAward, error) {
	var result StrategyAward
	nextNode := e.tree.RootRule

	for nextNode != "" {
		ruleNode, ok := e.tree.NodeMap[nextNode]
		if !ok {
			return StrategyAward{}, fmt.Errorf("rule tree node not found: %s", nextNode)
		}

		logicNode, ok := e.nodes[ruleNode.RuleKey]
		if !ok {
			return StrategyAward{}, fmt.Errorf("logic tree node not found: %s", ruleNode.RuleKey)
		}

		action, err := logicNode.Logic(ctx, userID, strategyID, awardID, ruleNode.RuleValue)
		if err != nil {
			return StrategyAward{}, err
		}
		if action.Award != nil {
			result = *action.Award
		}

		nextNode = e.nextNode(action.CheckType, ruleNode.Lines)
	}

	return result, nil
}

func (e *Engine) nextNode(checkType CheckType, lines []RuleTreeNodeLine) string {
	for _, line := range lines {
		if line.RuleLimitType == LimitEqual && matchCheckType(checkType, line.RuleLimitValue) {
			return line.RuleNodeTo
		}
	}
	return ""
}

func matchCheckType(checkType CheckType, value string) bool {
	return checkType.Code == value || checkType.Name == value
}
