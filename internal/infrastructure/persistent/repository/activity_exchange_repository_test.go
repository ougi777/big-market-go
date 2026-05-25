package repository

import (
	"context"
	"regexp"
	"testing"

	"bm-go/internal/domain/activity"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestActivityRepositoryQueryUnpaidActivityOrder(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	rows := sqlmock.NewRows([]string{"user_id", "sku", "order_id", "out_business_no", "pay_amount"}).
		AddRow("xiaofuge", 9011, "order-001", "biz-001", 1.68)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`sku`,`order_id`,`out_business_no`,`pay_amount` FROM `raffle_activity_order` WHERE user_id = ? and sku = ? and state = ? and order_time >= date_sub(now(), interval 1 month) ORDER BY `raffle_activity_order`.`id` LIMIT ?")).
		WithArgs("xiaofuge", int64(9011), activity.ActivityOrderWaitPay, 1).
		WillReturnRows(rows)

	order, exists, err := repo.QueryUnpaidActivityOrder(context.Background(), "xiaofuge", 9011)
	if err != nil {
		t.Fatalf("query unpaid activity order: %v", err)
	}
	if !exists {
		t.Fatal("expected unpaid order exists")
	}
	if order.OrderID != "order-001" || order.OutBusinessNo != "biz-001" || order.PayAmount != 1.68 {
		t.Fatalf("unexpected unpaid order: %+v", order)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestActivityRepositoryQueryUnpaidActivityOrderNotFound(t *testing.T) {
	db, mock := newMockGormDB(t)
	repo := NewActivityRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `user_id`,`sku`,`order_id`,`out_business_no`,`pay_amount` FROM `raffle_activity_order` WHERE user_id = ? and sku = ? and state = ? and order_time >= date_sub(now(), interval 1 month) ORDER BY `raffle_activity_order`.`id` LIMIT ?")).
		WithArgs("xiaofuge", int64(9011), activity.ActivityOrderWaitPay, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	order, exists, err := repo.QueryUnpaidActivityOrder(context.Background(), "xiaofuge", 9011)
	if err != nil {
		t.Fatalf("query unpaid activity order: %v", err)
	}
	if exists {
		t.Fatalf("expected unpaid order missing, got %+v", order)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
