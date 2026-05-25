package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"bm-go/internal/domain/activity"
	"bm-go/internal/types"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
)

func TestActivityRepositoryQuerySkuProductBySKU(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	skuRows := sqlmock.NewRows([]string{
		"sku",
		"activity_id",
		"activity_count_id",
		"stock_count",
		"stock_count_surplus",
		"product_amount",
	}).AddRow(9011, 100301, 11101, 100000, 99890, 1.68)
	countRows := sqlmock.NewRows([]string{
		"activity_count_id",
		"total_count",
		"day_count",
		"month_count",
	}).AddRow(11101, 100, 10, 30)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `sku`,`activity_id`,`activity_count_id`,`stock_count`,`stock_count_surplus`,`product_amount` FROM `raffle_activity_sku` WHERE sku = ? ORDER BY `raffle_activity_sku`.`id` LIMIT ?")).
		WithArgs(int64(9011), 1).
		WillReturnRows(skuRows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `activity_count_id`,`total_count`,`day_count`,`month_count` FROM `raffle_activity_count` WHERE activity_count_id = ? ORDER BY `raffle_activity_count`.`id` LIMIT ?")).
		WithArgs(int64(11101), 1).
		WillReturnRows(countRows)

	product, exists, err := repo.QuerySkuProductBySKU(context.Background(), 9011)
	if err != nil {
		t.Fatalf("query sku product: %v", err)
	}
	if !exists {
		t.Fatal("expected sku product exists")
	}
	if product.SKU != 9011 || product.ActivityID != 100301 || product.ActivityCount.TotalCount != 100 {
		t.Fatalf("unexpected product: %+v", product)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestActivityRepositorySaveCreditPayOrderDuplicate(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.Local)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `raffle_activity_order`")).
		WillReturnError(&mysql.MySQLError{Number: 1062, Message: "duplicate"})
	mock.ExpectRollback()

	err := repo.SaveCreditPayOrder(context.Background(), activity.CreateSkuExchangeOrderAggregate{
		UserID:     "xiaofuge",
		ActivityID: 100301,
		ActivityOrder: activity.ActivityOrderEntity{
			UserID:        "xiaofuge",
			SKU:           9011,
			ActivityID:    100301,
			ActivityName:  "大营销抽奖",
			StrategyID:    100006,
			OrderID:       "order-001",
			OrderTime:     now,
			TotalCount:    100,
			DayCount:      10,
			MonthCount:    30,
			PayAmount:     1.68,
			State:         activity.ActivityOrderWaitPay,
			OutBusinessNo: "biz-001",
		},
	})
	var appErr types.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != types.ResponseCodeIndexDup {
		t.Fatalf("expected duplicate code, got %s", appErr.Code.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
