package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStrategyRepositoryQueryAwardRuleLockCount(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	rows := sqlmock.NewRows([]string{"tree_id", "rule_value"}).
		AddRow("tree_lock_1", "1").
		AddRow("tree_lock_2", "2")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `tree_id`,`rule_value` FROM `rule_tree_node` WHERE rule_key = ? and tree_id in (?,?)")).
		WithArgs("rule_lock", "tree_lock_1", "tree_lock_2").
		WillReturnRows(rows)

	lockCounts, err := repo.QueryAwardRuleLockCount(context.Background(), []string{"tree_lock_1", "tree_lock_2"})
	if err != nil {
		t.Fatalf("query award rule lock count: %v", err)
	}
	if lockCounts["tree_lock_1"] != 1 || lockCounts["tree_lock_2"] != 2 {
		t.Fatalf("unexpected lock counts: %+v", lockCounts)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestStrategyRepositoryQueryAwardRuleLockCountEmpty(t *testing.T) {
	db, _ := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	lockCounts, err := repo.QueryAwardRuleLockCount(context.Background(), nil)
	if err != nil {
		t.Fatalf("query award rule lock count: %v", err)
	}
	if len(lockCounts) != 0 {
		t.Fatalf("expected empty lock counts, got %+v", lockCounts)
	}
}
