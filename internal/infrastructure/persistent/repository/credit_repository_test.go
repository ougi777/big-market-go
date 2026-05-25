package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestCreditRepositoryQueryUserCreditAccount(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewCreditRepository(db)

	rows := sqlmock.NewRows([]string{"user_id", "available_amount"}).
		AddRow("xiaofuge", 12.35)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`available_amount` FROM `user_credit_account` WHERE user_id = ? ORDER BY `user_credit_account`.`id` LIMIT ?")).
		WithArgs("xiaofuge", 1).
		WillReturnRows(rows)

	account, exists, err := repo.QueryUserCreditAccount(context.Background(), "xiaofuge")
	if err != nil {
		t.Fatalf("query user credit account: %v", err)
	}
	if !exists {
		t.Fatal("expected credit account exists")
	}
	if account.UserID != "xiaofuge" || account.AvailableAmount != 12.35 {
		t.Fatalf("unexpected credit account: %+v", account)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestCreditRepositoryQueryUserCreditAccountNotFound(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewCreditRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`available_amount` FROM `user_credit_account` WHERE user_id = ? ORDER BY `user_credit_account`.`id` LIMIT ?")).
		WithArgs("xiaofuge", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, exists, err := repo.QueryUserCreditAccount(context.Background(), "xiaofuge")
	if err != nil {
		t.Fatalf("query user credit account: %v", err)
	}
	if exists {
		t.Fatal("expected credit account not exists")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
