package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStrategyRepositoryUpdateStrategyAwardStock(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewStrategyRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `strategy_award` SET `award_count_surplus`=award_count_surplus - ? WHERE strategy_id = ? and award_id = ? and award_count_surplus > 0")).
		WithArgs(1, int64(100006), 101).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := repo.UpdateStrategyAwardStock(context.Background(), 100006, 101); err != nil {
		t.Fatalf("update strategy award stock: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestActivityRepositoryUpdateActivitySkuStock(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_sku` SET `stock_count_surplus`=stock_count_surplus - ?,`update_time`=? WHERE sku = ? and stock_count_surplus > 0")).
		WithArgs(1, sqlmock.AnyArg(), int64(9011)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := repo.UpdateActivitySkuStock(context.Background(), 9011); err != nil {
		t.Fatalf("update activity sku stock: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestActivityRepositoryClearActivitySkuStock(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_sku` SET `stock_count_surplus`=?,`update_time`=? WHERE sku = ?")).
		WithArgs(0, sqlmock.AnyArg(), int64(9011)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := repo.ClearActivitySkuStock(context.Background(), 9011); err != nil {
		t.Fatalf("clear activity sku stock: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
