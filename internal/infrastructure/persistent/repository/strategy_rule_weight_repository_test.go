package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStrategyRepositoryQueryAwardRuleWeight(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	ruleRows := sqlmock.NewRows([]string{"rule_value"}).
		AddRow("4000:101,102")
	awardRows := sqlmock.NewRows([]string{"award_id", "award_title"}).
		AddRow(101, "积分").
		AddRow(102, "优惠券")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `rule_value` FROM `strategy_rule` WHERE strategy_id = ? and rule_model = ? ORDER BY `strategy_rule`.`id` LIMIT ?")).
		WithArgs(int64(100001), "rule_weight", 1).
		WillReturnRows(ruleRows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `award_id`,`award_title` FROM `strategy_award` WHERE strategy_id = ? and award_id in (?,?) ORDER BY sort asc")).
		WithArgs(int64(100001), 101, 102).
		WillReturnRows(awardRows)

	ruleWeights, err := repo.QueryAwardRuleWeight(context.Background(), 100001)
	if err != nil {
		t.Fatalf("query award rule weight: %v", err)
	}
	if len(ruleWeights) != 1 {
		t.Fatalf("expected one rule weight, got %d", len(ruleWeights))
	}
	if ruleWeights[0].Weight != 4000 || len(ruleWeights[0].AwardList) != 2 {
		t.Fatalf("unexpected rule weight: %+v", ruleWeights[0])
	}
	if ruleWeights[0].AwardList[0].AwardID != 101 || ruleWeights[0].AwardList[1].AwardTitle != "优惠券" {
		t.Fatalf("unexpected award list: %+v", ruleWeights[0].AwardList)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
