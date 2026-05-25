package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStrategyRepositoryQueryRuleTreeByTreeID(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	treeRows := sqlmock.NewRows([]string{"tree_id", "tree_name", "tree_desc", "tree_node_rule_key"}).
		AddRow("tree_lock_stock", "库存规则树", "锁次数后校验库存", "rule_lock")
	nodeRows := sqlmock.NewRows([]string{"tree_id", "rule_key", "rule_desc", "rule_value"}).
		AddRow("tree_lock_stock", "rule_lock", "锁次数", "1").
		AddRow("tree_lock_stock", "rule_stock", "库存", "stock")
	lineRows := sqlmock.NewRows([]string{"tree_id", "rule_node_from", "rule_node_to", "rule_limit_type", "rule_limit_value"}).
		AddRow("tree_lock_stock", "rule_lock", "rule_stock", "EQUAL", "ALLOW")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `tree_id`,`tree_name`,`tree_desc`,`tree_node_rule_key` FROM `rule_tree` WHERE tree_id = ? ORDER BY `rule_tree`.`id` LIMIT ?")).
		WithArgs("tree_lock_stock", 1).
		WillReturnRows(treeRows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `tree_id`,`rule_key`,`rule_desc`,`rule_value` FROM `rule_tree_node` WHERE tree_id = ?")).
		WithArgs("tree_lock_stock").
		WillReturnRows(nodeRows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `tree_id`,`rule_node_from`,`rule_node_to`,`rule_limit_type`,`rule_limit_value` FROM `rule_tree_node_line` WHERE tree_id = ?")).
		WithArgs("tree_lock_stock").
		WillReturnRows(lineRows)

	ruleTree, exists, err := repo.QueryRuleTreeByTreeID(context.Background(), "tree_lock_stock")
	if err != nil {
		t.Fatalf("query rule tree: %v", err)
	}
	if !exists {
		t.Fatal("expected rule tree exists")
	}
	if ruleTree.RootRule != "rule_lock" || len(ruleTree.NodeMap) != 2 {
		t.Fatalf("unexpected rule tree: %+v", ruleTree)
	}
	if len(ruleTree.NodeMap["rule_lock"].Lines) != 1 || ruleTree.NodeMap["rule_lock"].Lines[0].RuleNodeTo != "rule_stock" {
		t.Fatalf("unexpected rule tree lines: %+v", ruleTree.NodeMap["rule_lock"].Lines)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
