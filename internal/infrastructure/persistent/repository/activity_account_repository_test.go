package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestActivityRepositoryQueryActivityAccount(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	rows := sqlmock.NewRows([]string{
		"user_id",
		"activity_id",
		"total_count",
		"total_count_surplus",
		"day_count",
		"day_count_surplus",
		"month_count",
		"month_count_surplus",
	}).AddRow("xiaofuge", 100301, 100, 80, 10, 8, 30, 20)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`activity_id`,`total_count`,`total_count_surplus`,`day_count`,`day_count_surplus`,`month_count`,`month_count_surplus` FROM `raffle_activity_account` WHERE user_id = ? and activity_id = ? ORDER BY `raffle_activity_account`.`id` LIMIT ?")).
		WithArgs("xiaofuge", int64(100301), 1).
		WillReturnRows(rows)

	account, exists, err := repo.QueryActivityAccount(context.Background(), 100301, "xiaofuge")
	if err != nil {
		t.Fatalf("query activity account: %v", err)
	}
	if !exists {
		t.Fatal("expected account exists")
	}
	if account.TotalCount != 100 || account.TotalCountSurplus != 80 || account.DayCount != 10 || account.MonthCount != 30 {
		t.Fatalf("unexpected account: %+v", account)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
