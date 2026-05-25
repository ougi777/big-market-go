package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"bm-go/internal/domain/credit"
	"bm-go/internal/types"

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

func TestCreditRepositoryCompleteCreditPayOrderQuotaNotEnough(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewCreditRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_credit_account` SET `available_amount`=available_amount + ?,`total_amount`=total_amount + ?,`update_time`=? WHERE user_id = ? and available_amount + ? >= 0")).
		WithArgs(-1.68, -1.68, sqlmock.AnyArg(), "xiaofuge", -1.68).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := repo.CompleteCreditPayOrder(context.Background(), credit.CompleteSkuExchangeAggregate{
		UserID: "xiaofuge",
		CreditOrder: credit.OrderEntity{
			UserID:        "xiaofuge",
			OrderID:       "credit-001",
			TradeName:     "兑换 SKU",
			TradeType:     "reverse",
			TradeAmount:   -1.68,
			OutBusinessNo: "biz-001",
		},
	})

	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != types.ResponseCodeAccountQuotaError {
		t.Fatalf("expected account quota code, got %s", appErr.Code.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestCreditRepositorySaveRebateIntegralOrderCreatesAccount(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewCreditRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_credit_account` SET `available_amount`=available_amount + ?,`total_amount`=total_amount + ?,`update_time`=? WHERE user_id = ? and available_amount + ? >= 0")).
		WithArgs(10.0, 10.0, sqlmock.AnyArg(), "xiaofuge", 10.0).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_credit_account`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_credit_order`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.SaveRebateIntegralOrder(context.Background(), credit.RebateIntegralEntity{
		UserID:        "xiaofuge",
		OrderID:       "rebate-001",
		TradeAmount:   10,
		OutBusinessNo: "sign-20260525",
	})
	if err != nil {
		t.Fatalf("save rebate integral order: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
