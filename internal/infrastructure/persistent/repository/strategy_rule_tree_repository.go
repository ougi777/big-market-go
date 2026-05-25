package repository

import (
	"context"
	"errors"

	treepkg "bm-go/internal/domain/strategy/rule/tree"
	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

func (r *StrategyRepository) QueryRuleTreeByTreeID(ctx context.Context, treeID string) (treepkg.RuleTree, bool, error) {
	var ruleTreePO po.RuleTree
	err := r.defaultDB(ctx).
		Select("tree_id", "tree_name", "tree_desc", "tree_node_rule_key").
		Where("tree_id = ?", treeID).
		First(&ruleTreePO).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return treepkg.RuleTree{}, false, nil
	}
	if err != nil {
		return treepkg.RuleTree{}, false, err
	}

	var nodePOList []po.RuleTreeNode
	if err := r.defaultDB(ctx).
		Select("tree_id", "rule_key", "rule_desc", "rule_value").
		Where("tree_id = ?", treeID).
		Find(&nodePOList).
		Error; err != nil {
		return treepkg.RuleTree{}, false, err
	}

	var linePOList []po.RuleTreeNodeLine
	if err := r.defaultDB(ctx).
		Select("tree_id", "rule_node_from", "rule_node_to", "rule_limit_type", "rule_limit_value").
		Where("tree_id = ?", treeID).
		Find(&linePOList).
		Error; err != nil {
		return treepkg.RuleTree{}, false, err
	}

	lineMap := make(map[string][]treepkg.RuleTreeNodeLine)
	for _, linePO := range linePOList {
		lineMap[linePO.RuleNodeFrom] = append(lineMap[linePO.RuleNodeFrom], treepkg.RuleTreeNodeLine{
			TreeID:         linePO.TreeID,
			RuleNodeFrom:   linePO.RuleNodeFrom,
			RuleNodeTo:     linePO.RuleNodeTo,
			RuleLimitType:  linePO.RuleLimitType,
			RuleLimitValue: linePO.RuleLimitValue,
		})
	}

	nodeMap := make(map[string]treepkg.RuleTreeNode, len(nodePOList))
	for _, nodePO := range nodePOList {
		nodeMap[nodePO.RuleKey] = treepkg.RuleTreeNode{
			TreeID:    nodePO.TreeID,
			RuleKey:   nodePO.RuleKey,
			RuleDesc:  nodePO.RuleDesc,
			RuleValue: nodePO.RuleValue,
			Lines:     lineMap[nodePO.RuleKey],
		}
	}

	return treepkg.RuleTree{
		TreeID:   ruleTreePO.TreeID,
		TreeName: ruleTreePO.TreeName,
		TreeDesc: ruleTreePO.TreeDesc,
		RootRule: ruleTreePO.TreeRootRuleKey,
		NodeMap:  nodeMap,
	}, true, nil
}
