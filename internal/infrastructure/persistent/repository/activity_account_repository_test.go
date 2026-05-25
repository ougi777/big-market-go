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

func TestActivityRepositoryQueryActivityAccountDay(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	rows := sqlmock.NewRows([]string{
		"user_id",
		"activity_id",
		"day",
		"day_count",
		"day_count_surplus",
	}).AddRow("xiaofuge", 100301, "2026-05-25", 5, 3)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`activity_id`,`day`,`day_count`,`day_count_surplus` FROM `raffle_activity_account_day` WHERE user_id = ? and activity_id = ? and day = ? ORDER BY `raffle_activity_account_day`.`id` LIMIT ?")).
		WithArgs("xiaofuge", int64(100301), "2026-05-25", 1).
		WillReturnRows(rows)

	account, exists, err := repo.QueryActivityAccountDay(context.Background(), 100301, "xiaofuge", "2026-05-25")
	if err != nil {
		t.Fatalf("query activity account day: %v", err)
	}
	if !exists {
		t.Fatal("expected day account exists")
	}
	if account.Day != "2026-05-25" || account.DayCount != 5 || account.DayCountSurplus != 3 {
		t.Fatalf("unexpected day account: %+v", account)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestActivityRepositoryQueryActivityAccountMonth(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	rows := sqlmock.NewRows([]string{
		"user_id",
		"activity_id",
		"month",
		"month_count",
		"month_count_surplus",
	}).AddRow("xiaofuge", 100301, "2026-05", 30, 20)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`activity_id`,`month`,`month_count`,`month_count_surplus` FROM `raffle_activity_account_month` WHERE user_id = ? and activity_id = ? and month = ? ORDER BY `raffle_activity_account_month`.`id` LIMIT ?")).
		WithArgs("xiaofuge", int64(100301), "2026-05", 1).
		WillReturnRows(rows)

	account, exists, err := repo.QueryActivityAccountMonth(context.Background(), 100301, "xiaofuge", "2026-05")
	if err != nil {
		t.Fatalf("query activity account month: %v", err)
	}
	if !exists {
		t.Fatal("expected month account exists")
	}
	if account.Month != "2026-05" || account.MonthCount != 30 || account.MonthCountSurplus != 20 {
		t.Fatalf("unexpected month account: %+v", account)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
