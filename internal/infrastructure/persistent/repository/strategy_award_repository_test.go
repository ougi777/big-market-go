package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStrategyRepositoryQueryStrategyAwardList(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	rows := sqlmock.NewRows([]string{
		"strategy_id",
		"award_id",
		"award_title",
		"award_subtitle",
		"award_count",
		"award_count_surplus",
		"award_rate",
		"rule_models",
		"sort",
	}).AddRow(100001, 101, "积分", "抽奖1次后解锁", 100, 80, 0.1, "tree_lock_1", 1).
		AddRow(100001, 102, "优惠券", "", 50, 30, 0.2, "", 2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `strategy_id`,`award_id`,`award_title`,`award_subtitle`,`award_count`,`award_count_surplus`,`award_rate`,`rule_models`,`sort` FROM `strategy_award` WHERE strategy_id = ? ORDER BY sort asc")).
		WithArgs(int64(100001)).
		WillReturnRows(rows)

	awards, err := repo.QueryStrategyAwardList(context.Background(), 100001)
	if err != nil {
		t.Fatalf("query strategy awards: %v", err)
	}
	if len(awards) != 2 {
		t.Fatalf("expected 2 awards, got %d", len(awards))
	}
	if awards[0].AwardID != 101 || awards[0].RuleModels != "tree_lock_1" || awards[1].AwardID != 102 {
		t.Fatalf("unexpected awards: %+v", awards)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestStrategyRepositoryQueryStrategyAwardEntity(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	rows := sqlmock.NewRows([]string{
		"strategy_id",
		"award_id",
		"award_title",
		"award_subtitle",
		"award_count",
		"award_count_surplus",
		"award_rate",
		"rule_models",
		"sort",
	}).AddRow(100001, 101, "积分", "抽奖1次后解锁", 100, 80, 0.1, "tree_lock_1", 1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `strategy_id`,`award_id`,`award_title`,`award_subtitle`,`award_count`,`award_count_surplus`,`award_rate`,`rule_models`,`sort` FROM `strategy_award` WHERE strategy_id = ? and award_id = ? ORDER BY `strategy_award`.`id` LIMIT ?")).
		WithArgs(int64(100001), 101, 1).
		WillReturnRows(rows)

	award, exists, err := repo.QueryStrategyAwardEntity(context.Background(), 100001, 101)
	if err != nil {
		t.Fatalf("query strategy award: %v", err)
	}
	if !exists {
		t.Fatal("expected strategy award exists")
	}
	if award.AwardID != 101 || award.AwardTitle != "积分" || award.Sort != 1 {
		t.Fatalf("unexpected award: %+v", award)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
