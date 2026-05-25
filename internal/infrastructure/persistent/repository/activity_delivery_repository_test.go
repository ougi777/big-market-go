package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestActivityRepositoryDeliverActivityOrder(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	rows := sqlmock.NewRows([]string{"user_id", "activity_id", "total_count", "day_count", "month_count", "state"}).
		AddRow("xiaofuge", 100301, 100, 10, 30, activity.ActivityOrderWaitPay)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`activity_id`,`total_count`,`day_count`,`month_count`,`state` FROM `raffle_activity_order` WHERE user_id = ? and out_business_no = ? ORDER BY `raffle_activity_order`.`id` LIMIT ?")).
		WithArgs("xiaofuge", "biz-001", 1).
		WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_order`")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account`")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account_month`")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account_day`")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeliverActivityOrder(context.Background(), activity.DeliveryOrderEntity{
		UserID:        "xiaofuge",
		OutBusinessNo: "biz-001",
	})
	if err != nil {
		t.Fatalf("deliver activity order: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestActivityRepositoryDeliverActivityOrderCreatesAccount(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	rows := sqlmock.NewRows([]string{"user_id", "activity_id", "total_count", "day_count", "month_count", "state"}).
		AddRow("xiaofuge", 100301, 100, 10, 30, activity.ActivityOrderWaitPay)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`activity_id`,`total_count`,`day_count`,`month_count`,`state` FROM `raffle_activity_order` WHERE user_id = ? and out_business_no = ? ORDER BY `raffle_activity_order`.`id` LIMIT ?")).
		WithArgs("xiaofuge", "biz-001", 1).
		WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_order`")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account`")).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `raffle_activity_account`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account_month`")).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `raffle_activity_account_month`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `raffle_activity_account_day`")).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `raffle_activity_account_day`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeliverActivityOrder(context.Background(), activity.DeliveryOrderEntity{
		UserID:        "xiaofuge",
		OutBusinessNo: "biz-001",
	})
	if err != nil {
		t.Fatalf("deliver activity order: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestActivityRepositoryDeliverActivityOrderNotFound(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`activity_id`,`total_count`,`day_count`,`month_count`,`state` FROM `raffle_activity_order` WHERE user_id = ? and out_business_no = ? ORDER BY `raffle_activity_order`.`id` LIMIT ?")).
		WithArgs("xiaofuge", "biz-001", 1).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	err := repo.DeliverActivityOrder(context.Background(), activity.DeliveryOrderEntity{
		UserID:        "xiaofuge",
		OutBusinessNo: "biz-001",
	})

	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != types.ResponseCodeIllegalParam {
		t.Fatalf("expected illegal param code, got %s", appErr.Code.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
