package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestStrategyRepositoryQueryActivityAccountTotalUseCount(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	activityRows := sqlmock.NewRows([]string{"activity_id"}).AddRow(100301)
	accountRows := sqlmock.NewRows([]string{"total_count", "total_count_surplus"}).AddRow(100, 88)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `activity_id` FROM `raffle_activity` WHERE strategy_id = ? ORDER BY `raffle_activity`.`id` LIMIT ?")).
		WithArgs(int64(100006), 1).
		WillReturnRows(activityRows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `total_count`,`total_count_surplus` FROM `raffle_activity_account` WHERE user_id = ? and activity_id = ? ORDER BY `raffle_activity_account`.`id` LIMIT ?")).
		WithArgs("xiaofuge", int64(100301), 1).
		WillReturnRows(accountRows)

	useCount, err := repo.QueryActivityAccountTotalUseCount(context.Background(), "xiaofuge", 100006)
	if err != nil {
		t.Fatalf("query activity account total use count: %v", err)
	}
	if useCount != 12 {
		t.Fatalf("unexpected use count: %d", useCount)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestStrategyRepositoryQueryActivityAccountTotalUseCountActivityNotFound(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `activity_id` FROM `raffle_activity` WHERE strategy_id = ? ORDER BY `raffle_activity`.`id` LIMIT ?")).
		WithArgs(int64(100006), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	useCount, err := repo.QueryActivityAccountTotalUseCount(context.Background(), "xiaofuge", 100006)
	if err != nil {
		t.Fatalf("query activity account total use count: %v", err)
	}
	if useCount != 0 {
		t.Fatalf("unexpected use count: %d", useCount)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestStrategyRepositoryQueryRaffleActivityAccountDayPartakeCount(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)
	today := time.Now().Format("2006-01-02")

	rows := sqlmock.NewRows([]string{"day_count", "day_count_surplus"}).AddRow(10, 7)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `day_count`,`day_count_surplus` FROM `raffle_activity_account_day` WHERE user_id = ? and activity_id = ? and day = ? ORDER BY `raffle_activity_account_day`.`id` LIMIT ?")).
		WithArgs("xiaofuge", int64(100301), today, 1).
		WillReturnRows(rows)

	partakeCount, err := repo.QueryRaffleActivityAccountDayPartakeCount(context.Background(), 100301, "xiaofuge")
	if err != nil {
		t.Fatalf("query raffle activity account day partake count: %v", err)
	}
	if partakeCount != 3 {
		t.Fatalf("unexpected partake count: %d", partakeCount)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestStrategyRepositoryQueryRaffleActivityAccountPartakeCount(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	rows := sqlmock.NewRows([]string{"total_count", "total_count_surplus"}).AddRow(100, 90)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `total_count`,`total_count_surplus` FROM `raffle_activity_account` WHERE user_id = ? and activity_id = ? ORDER BY `raffle_activity_account`.`id` LIMIT ?")).
		WithArgs("xiaofuge", int64(100301), 1).
		WillReturnRows(rows)

	partakeCount, err := repo.QueryRaffleActivityAccountPartakeCount(context.Background(), 100301, "xiaofuge")
	if err != nil {
		t.Fatalf("query raffle activity account partake count: %v", err)
	}
	if partakeCount != 10 {
		t.Fatalf("unexpected partake count: %d", partakeCount)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
