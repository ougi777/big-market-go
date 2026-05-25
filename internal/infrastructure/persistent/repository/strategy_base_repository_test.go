package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestStrategyRepositoryQueryStrategyEntityByStrategyID(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	rows := sqlmock.NewRows([]string{
		"strategy_id",
		"strategy_desc",
		"rule_models",
	}).AddRow(100001, "默认抽奖策略", "rule_weight,rule_blacklist")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `strategy_id`,`strategy_desc`,`rule_models` FROM `strategy` WHERE strategy_id = ? ORDER BY `strategy`.`id` LIMIT ?")).
		WithArgs(int64(100001), 1).
		WillReturnRows(rows)

	entity, err := repo.QueryStrategyEntityByStrategyID(context.Background(), 100001)
	if err != nil {
		t.Fatalf("query strategy: %v", err)
	}
	if entity.StrategyID != 100001 || entity.RuleModel != "rule_weight,rule_blacklist" {
		t.Fatalf("unexpected strategy: %+v", entity)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestStrategyRepositoryQueryStrategyEntityByStrategyIDNotFound(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `strategy_id`,`strategy_desc`,`rule_models` FROM `strategy` WHERE strategy_id = ? ORDER BY `strategy`.`id` LIMIT ?")).
		WithArgs(int64(100001), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	entity, err := repo.QueryStrategyEntityByStrategyID(context.Background(), 100001)
	if err != nil {
		t.Fatalf("query strategy: %v", err)
	}
	if entity.StrategyID != 0 {
		t.Fatalf("expected empty strategy, got %+v", entity)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestStrategyRepositoryQueryStrategyIDByActivityID(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	rows := sqlmock.NewRows([]string{"strategy_id"}).AddRow(100006)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `strategy_id` FROM `raffle_activity` WHERE activity_id = ? ORDER BY `raffle_activity`.`id` LIMIT ?")).
		WithArgs(int64(100301), 1).
		WillReturnRows(rows)

	strategyID, err := repo.QueryStrategyIDByActivityID(context.Background(), 100301)
	if err != nil {
		t.Fatalf("query strategy id: %v", err)
	}
	if strategyID != 100006 {
		t.Fatalf("expected strategy id 100006, got %d", strategyID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
