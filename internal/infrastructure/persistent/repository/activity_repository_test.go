package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestActivityRepositoryQueryActivityByActivityID(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)
	begin := time.Date(2026, 5, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(2026, 5, 31, 23, 59, 59, 0, time.Local)

	rows := sqlmock.NewRows([]string{
		"activity_id",
		"activity_name",
		"activity_desc",
		"begin_date_time",
		"end_date_time",
		"strategy_id",
		"state",
	}).AddRow(100301, "大营销抽奖", "测试活动", begin, end, 100006, "open")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `activity_id`,`activity_name`,`activity_desc`,`begin_date_time`,`end_date_time`,`strategy_id`,`state` FROM `raffle_activity` WHERE activity_id = ? ORDER BY `raffle_activity`.`id` LIMIT ?")).
		WithArgs(int64(100301), 1).
		WillReturnRows(rows)

	activityEntity, exists, err := repo.QueryActivityByActivityID(context.Background(), 100301)
	if err != nil {
		t.Fatalf("query activity: %v", err)
	}
	if !exists {
		t.Fatal("expected activity exists")
	}
	if activityEntity.ActivityID != 100301 || activityEntity.StrategyID != 100006 || activityEntity.State != "open" {
		t.Fatalf("unexpected activity: %+v", activityEntity)
	}
	if !activityEntity.BeginDateTime.Equal(begin) || !activityEntity.EndDateTime.Equal(end) {
		t.Fatalf("unexpected activity time: %+v", activityEntity)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
